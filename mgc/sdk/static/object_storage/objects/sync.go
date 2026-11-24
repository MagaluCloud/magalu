package objects

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
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
	BatchSize int              `json:"batch_size,omitempty" jsonschema:"description=Limit of items per batch to delete,default=1000,minimum=1,maximum=1000,required" example:"1000"`
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
		return "template={{if and (eq .deleted 0) (eq .uploaded 0)}}No files uploaded or deleted{{- else}}" +
			"Synced files from {{.src}} to {{.dst}}\n- {{.uploaded}} files uploaded\n- {{.deleted}} files deleted{{if .hasDeleted}}\n\nDeleted files:\n-{{.deletedFiles}}{{end}}{{- end}}\n"
	})
})

type syncState struct {
	mu          sy.Mutex
	bucketFiles map[string]bool
	uploadFiles UploadCounter
}

func sync(ctx context.Context, params syncParams, cfg common.Config) (result core.Value, err error) {
	if params.BatchSize < common.MinBatchSize || params.BatchSize > common.MaxBatchSize {
		return nil, core.UsageError{Err: fmt.Errorf("invalid item limit per request BatchSize, must not be lower than %d and must not be higher than %d: %d", common.MinBatchSize, common.MaxBatchSize, params.BatchSize)}
	}
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

	info, err := os.Stat(basePath.String())
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
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
		defer func() { _, _ = progressBar.Stop() }()
	}

	bucketFiles, err := fillBucketFiles(ctx, params, cfg)
	if err != nil {
		return nil, err
	}
	state := &syncState{bucketFiles: bucketFiles}

	err = processSyncFiles(ctx, cfg, params.Bucket, basePath.String(), files, progressBar, state)

	if err != nil {
		return nil, err
	}

	_, _ = progressBar.Stop()

	deletedFiles := make([]string, 0, len(state.bucketFiles))

	if params.Delete && len(state.bucketFiles) > 0 {
		for file := range state.bucketFiles {
			deletedFiles = append(deletedFiles, file)
		}
		delOb := common.DeleteObjectsParams{
			Destination: params.Bucket,
			ToDelete:    bucketObjectsToWalkDirEntry(ctx, deletedFiles),
			BatchSize:   params.BatchSize,
		}
		err = common.DeleteObjects(ctx, delOb, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to delete objects (some objects may already have been deleted): %w", err)
		}
	}

	return syncResult{
		Source:        params.Local,
		Destination:   params.Bucket,
		FilesDeleted:  len(deletedFiles),
		FilesUploaded: int(state.uploadFiles.Value()),
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
			select {
			case out <- entry:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func fillBucketFiles(ctx context.Context, params syncParams, cfg common.Config) (map[string]bool, error) {
	bucketFiles := make(map[string]bool)
	logger().Debug("Getting bucket files")

	dirBucketFiles := common.ListGenerator(ctx, common.ListObjectsParams{
		Destination: params.Bucket,
		Recursive:   true,
		PaginationParams: common.PaginationParams{
			MaxItems: math.MaxInt,
		},
	}, cfg, nil)

	for file := range dirBucketFiles {
		if err := file.Err(); err != nil {
			return nil, err
		}
		bucketFiles["/"+file.Path()] = true
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return bucketFiles, nil
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

func processSyncFiles(ctx context.Context, cfg common.Config, destination mgcSchemaPkg.URI, basePath string, files []string, progressBar *pterm.ProgressbarPrinter, state *syncState) error {
	results := make(chan error, cfg.Workers)
	filesChan := make(chan string, cfg.Workers)

	var wg sy.WaitGroup
	wg.Add(cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		go func() {
			defer wg.Done()
			syncWorker(ctx, cfg, destination, basePath, filesChan, results, progressBar, state)
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

	return ctx.Err()
}

func syncWorker(ctx context.Context, cfg common.Config, destination mgcSchemaPkg.URI, basePath string, files <-chan string, results chan<- error, progressBar *pterm.ProgressbarPrinter, state *syncState) {
	for {
		select {
		case file, ok := <-files:
			if !ok {
				return
			}
			err := processSyncFile(ctx, cfg, destination, basePath, file, progressBar, state)
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

func processSyncFile(ctx context.Context, cfg common.Config, destination mgcSchemaPkg.URI, basePath, file string, progressBar *pterm.ProgressbarPrinter, state *syncState) error {
	normalizedSource, err := common.GetAbsSystemURI(mgcSchemaPkg.URI(file))
	if err != nil {
		logger().Debugw("error with path", "error", err)
		return err
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

	state.mu.Lock()
	delete(state.bucketFiles, "/"+normalizedDestination.Path())
	state.mu.Unlock()

	fileStats, err := getFileStats(ctx, normalizedDestination, cfg)
	if err != nil {
		logger().Debugw("error getting file stats", "error", err)
	}

	isSameSize := info.Size() == fileStats.SourceLength
	isLocalOlderThenBucket := info.ModTime().Unix() < fileStats.SourceModTime
	if err == nil && isSameSize && isLocalOlderThenBucket {
		logger().Debugf("Skipping file [%s] - no change", normalizedSource)
		state.mu.Lock()
		progressBar.Increment()
		state.mu.Unlock()
		return nil
	}

	err = uploadFile(ctx, normalizedSource, normalizedDestination, cfg)
	if err != nil {
		return &common.ObjectError{Url: mgcSchemaPkg.URI(normalizedSource.Path()), Err: err}
	}

	state.uploadFiles.Increment()
	state.mu.Lock()
	progressBar.Increment()
	state.mu.Unlock()
	return nil
}
