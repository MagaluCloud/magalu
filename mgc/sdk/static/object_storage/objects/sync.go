package objects

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	sy "sync"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/pipeline"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/openapi"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
	"github.com/pterm/pterm"
)

type UploadCounter struct {
	mu sy.Mutex
	v  uint64
}

type fileSyncStats struct {
	SourceLength  int64
	SourceModTime int64
	Etag          string
}

type syncParams struct {
	Local     mgcSchemaPkg.URI `json:"local" jsonschema:"description=Local path,example=./" mgc:"positional"`
	Bucket    mgcSchemaPkg.URI `json:"bucket" jsonschema:"description=Bucket path,example=my-bucket/dir/" mgc:"positional"`
	Delete    bool             `json:"delete,omitempty" jsonschema:"description=Deletes any item at the bucket not present on the local,default=false"`
	BatchSize int              `json:"batch_size,omitempty" jsonschema:"description=Limit of items per batch to delete,default=1000,minimum=1,maximum=1000" example:"1000"`
}

type syncResult struct {
	Source        mgcSchemaPkg.URI `json:"src" jsonschema:"description=Source path to sync the remote with,example=./" mgc:"positional"`
	Destination   mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Full destination path to sync with the source path,example=s3://my-bucket/dir/" mgc:"positional"`
	FilesDeleted  int              `json:"deleted"`
	FilesUploaded int              `json:"uploaded"`
	Deleted       bool             `json:"hasDeleted"`
	DeletedFiles  string           `json:"deletedFiles"`
}

var getSync = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "sync",
			Summary:     "Synchronizes a local path with a bucket",
			Description: "This command uploads any file from the local path to the bucket if it is not already present or has modified time changed.",
		},
		sync,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template={{if and (eq .deleted 0) (eq .uploaded 0)}}Already Synced{{- else}}" +
			"Synced files from {{.src}} to {{.dst}}\n- {{.uploaded}} files uploaded\n- {{if .hasDeleted}}{{.deleted}} files deleted\n\nDeleted files:\n-{{.deletedFiles}}{{- else}}{{.deleted}} files to be deleted with the --delete parameter{{- end}}{{- end}}\n"
	})
})

// bucketFiles holds the objects found in the bucket that were not seen
// locally yet, keyed by their path relative to the sync destination. The
// value is the full object key, used when deleting.
type bucketFiles struct {
	mu   sy.Mutex
	keys map[string]string
}

func (b *bucketFiles) markSeen(relPath string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.keys, relPath)
}

func sync(ctx context.Context, params syncParams, cfg common.Config) (result core.Value, err error) {
	if !strings.HasPrefix(string(params.Bucket), common.URIPrefix) {
		logger().Debugw("Bucket path missing prefix, adding prefix")
		params.Bucket = common.URIPrefix + params.Bucket
	}

	if strings.HasPrefix(string(params.Local), common.URIPrefix) {
		return nil, fmt.Errorf("local cannot be an bucket! To copy or move between buckets, use \"mgc object-storage objects copy/move\"")
	}

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	basePath, err := common.GetAbsSystemURI(params.Local)
	if err != nil {
		return nil, err
	}

	f, err := os.Stat(basePath.String())
	if err != nil {
		return nil, err
	}
	if !f.IsDir() {
		return nil, fmt.Errorf("local path must be a folder")
	}

	files, err := walkDir(ctx, basePath.String(), false)
	if err != nil {
		return nil, err
	}

	totalFiles := len(files)
	progressBar := pterm.DefaultProgressbar.
		WithTotal(totalFiles).
		WithTitle("Syncing files").
		WithRemoveWhenDone(true)

	if !openapi.GetRawOutputFlag(ctx) {
		progressBar, _ = progressBar.Start()
	}

	remoteFiles, err := fillBucketFiles(ctx, params, cfg)
	if err != nil {
		return nil, err
	}

	uploadFiles := &UploadCounter{}

	err = processSyncFiles(ctx, cfg, params.Local, params.Bucket, basePath.String(), files, remoteFiles, uploadFiles, progressBar)

	if err != nil {
		return nil, err
	}

	_, _ = progressBar.Stop()

	deletedFiles := make([]string, 0, len(remoteFiles.keys))

	if params.Delete {
		for _, key := range remoteFiles.keys {
			deletedFiles = append(deletedFiles, key)
		}
		delOb := common.DeleteObjectsParams{
			Destination: params.Bucket,
			ToDelete:    bucketObjectsToWalkDirEntry(ctx, deletedFiles),
			BatchSize:   params.BatchSize,
		}
		err = common.DeleteObjects(ctx, delOb, cfg)
		if err != nil {
			logger().Debugw("error deleting objects", "error", err)
		}
	}

	return syncResult{
		Source:        params.Local,
		Destination:   params.Bucket,
		FilesDeleted:  len(remoteFiles.keys),
		FilesUploaded: int(uploadFiles.Value()),
		Deleted:       len(deletedFiles) > 0,
		DeletedFiles:  strings.Join(deletedFiles, ", "),
	}, nil
}

func bucketObjectsToWalkDirEntry(ctx context.Context, bucketObjects []string) <-chan pipeline.WalkDirEntry {
	out := make(chan pipeline.WalkDirEntry)
	go func() {
		defer close(out)
		var err error
		for _, obj := range bucketObjects {
			if ctx.Err() != nil {
				return
			}
			entry := pipeline.NewSimpleWalkDirEntry(obj, &common.BucketContent{
				Key: strings.TrimPrefix(obj, "/"),
			}, err)
			out <- entry
		}
	}()
	return out
}

func fillBucketFiles(ctx context.Context, params syncParams, cfg common.Config) (*bucketFiles, error) {
	logger().Debug("Getting bucket files")

	dirBucketFiles := common.ListGenerator(ctx, common.ListObjectsParams{
		Destination: params.Bucket,
		Recursive:   true,
		PaginationParams: common.PaginationParams{
			MaxItems: math.MaxInt64,
		},
	}, cfg, nil)

	// Listed keys are relative to the bucket root, local paths are relative
	// to the destination, so strip the destination prefix before comparing.
	prefix := params.Bucket.Path()
	if prefix != "" {
		prefix += "/"
	}

	files := &bucketFiles{keys: make(map[string]string)}
	for file := range dirBucketFiles {
		if err := file.Err(); err != nil {
			return nil, err
		}
		key := file.Path()
		files.keys[strings.TrimPrefix(key, prefix)] = key
	}
	return files, nil
}

func getFileStats(ctx context.Context, destination mgcSchemaPkg.URI, cfg common.Config) (fileSyncStats, error) {
	dstHead, err := headObject(ctx, headObjectParams{
		Destination: destination,
	}, cfg)
	if err != nil {
		return fileSyncStats{}, err
	}
	dstModTime, err := time.Parse(time.RFC1123, dstHead.LastModified)
	if err != nil {
		logger().Debug("%s %s\n", dstModTime, err)
		return fileSyncStats{}, err
	}
	return fileSyncStats{
		SourceLength:  dstHead.ContentLength,
		SourceModTime: dstModTime.Unix(),
		Etag:          cleanEtag(dstHead.ETag),
	}, nil
}

func cleanEtag(etag string) string {
	return strings.Trim(etag, "\"")
}

func uploadFile(ctx context.Context, local mgcSchemaPkg.URI, bucket mgcSchemaPkg.URI, cfg common.Config) error {
	_, err := upload(
		ctx,
		uploadParams{Source: mgcSchemaPkg.FilePath(local), Destination: bucket},
		cfg,
	)
	return err
}

func (c *UploadCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v++
}

func (c *UploadCounter) Value() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v
}

func processSyncFiles(ctx context.Context, cfg common.Config, source, destination mgcSchemaPkg.URI, basePath string, files []string, remoteFiles *bucketFiles, uploadFiles *UploadCounter, progressBar *pterm.ProgressbarPrinter) error {
	results := make(chan error, cfg.Workers)
	filesChan := make(chan string, cfg.Workers)

	var wg sy.WaitGroup
	wg.Add(cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		go func() {
			defer wg.Done()
			syncWorker(ctx, cfg, source, destination, basePath, filesChan, results, remoteFiles, uploadFiles, progressBar)
		}()
	}

	go func() {
		defer close(filesChan)
		for _, file := range files {
			select {
			case filesChan <- file:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var errs []error
	for err := range results {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d of %d file(s) failed to sync:\n%w", len(errs), len(files), errors.Join(errs...))
	}

	return nil
}

func syncWorker(ctx context.Context, cfg common.Config, source, destination mgcSchemaPkg.URI, basePath string, files <-chan string, results chan<- error, remoteFiles *bucketFiles, uploadFiles *UploadCounter, progressBar *pterm.ProgressbarPrinter) {
	for {
		select {
		case file, ok := <-files:
			if !ok {
				return
			}
			err := processSyncFile(ctx, cfg, source, destination, basePath, file, remoteFiles, uploadFiles, progressBar)
			if err != nil {
				select {
				case results <- err:
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func processSyncFile(ctx context.Context, cfg common.Config, source, destination mgcSchemaPkg.URI, basePath, file string, remoteFiles *bucketFiles, uploadFiles *UploadCounter, progressBar *pterm.ProgressbarPrinter) error {
	normalizedSource, err := common.GetAbsSystemURI(mgcSchemaPkg.URI(file))
	if err != nil {
		logger().Debugw("error with path", "error", err)
		return nil
	}

	pathWithFolder := strings.TrimPrefix(file, basePath)

	fixedPathWithFolder, converted, err := common.FixFilenameEncoding(pathWithFolder)
	if err != nil {
		return &common.ObjectError{Url: mgcSchemaPkg.URI(pathWithFolder), Err: fmt.Errorf("skipping file: %w", err)}
	}
	if converted {
		logger().Warnw("converted filename encoding for sync", "original", pathWithFolder, "converted", fixedPathWithFolder)
	}
	pathWithFolder = fixedPathWithFolder

	normalizedDestination := destination.JoinPath(pathWithFolder)

	info, err := os.Stat(file)
	if err != nil {
		return err
	}

	remoteFiles.markSeen(strings.TrimPrefix(filepath.ToSlash(pathWithFolder), "/"))

	fileStats, err := getFileStats(ctx, normalizedDestination, cfg)
	if err != nil {
		logger().Debugw("error getting file stats", "error", err)
	}

	isSameSize := info.Size() == fileStats.SourceLength
	isLocalOlderThenBucket := info.ModTime().Unix() < fileStats.SourceModTime
	if err == nil && isSameSize && isLocalOlderThenBucket {
		logger().Debug("Skipping file [%s] - no change", normalizedSource)
		progressBar.Increment()
		return nil
	}

	err = uploadFile(ctx, normalizedSource, normalizedDestination, cfg)
	if err != nil {
		return &common.ObjectError{Url: mgcSchemaPkg.URI(normalizedSource.Path()), Err: err}
	}

	uploadFiles.Increment()
	progressBar.Increment()
	return nil
}
