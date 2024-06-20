package objects

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"
	sy "sync"
	"time"

	"magalu.cloud/core"
	"magalu.cloud/core/pipeline"
	"magalu.cloud/core/progress_report"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type UploadCounter struct {
	mu sy.Mutex
	v  uint64
}

type syncUploadPair struct {
	Source      mgcSchemaPkg.URI
	Destination mgcSchemaPkg.URI
	Stats       fileSyncStats
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

var (
	allBucketFiles = make(map[string]bool)
	uploadFiles    = &UploadCounter{}
)

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

	srcObjects := pipeline.WalkDirEntries(ctx, basePath.String(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return nil
	})

	progressReporter := progress_report.NewUnitsReporter(ctx, "Syncing Folder", 0)
	progressReporter.Start()

	defer progressReporter.End()

	fillBucketFiles(ctx, params, cfg)

	if f, _ := os.Stat(basePath.String()); !f.IsDir() {
		return nil, fmt.Errorf("local path must be a folder")
	}

	uploadChannel := pipeline.Process(ctx, srcObjects, createObjectSyncFilePairProcessor(cfg, params.Local, params.Bucket, progressReporter, basePath), nil)
	uploadObjectsErrorChan := pipeline.ParallelProcess(ctx, cfg.Workers, uploadChannel, createSyncObjectProcessor(cfg, progressReporter), nil)
	objErr, err := pipeline.SliceItemConsumer[utils.MultiError](ctx, uploadObjectsErrorChan)

	for _, er := range objErr {
		if er != nil {
			progressReporter.Report(0, 0, er)
			return nil, objErr
		}
	}

	deletedFiles := make([]string, 0, len(allBucketFiles))

	if params.Delete {
		for file := range allBucketFiles {
			if err != nil {
				logger().Debugw("error deleting file", "error", err)
			}
			deletedFiles = append(deletedFiles, file)
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
		FilesDeleted:  len(allBucketFiles),
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

func fillBucketFiles(ctx context.Context, params syncParams, cfg common.Config) {
	logger().Debug("Getting bucket files")

	dirBucketFiles := common.ListGenerator(ctx, common.ListObjectsParams{
		Destination: params.Bucket,
		Recursive:   true,
		PaginationParams: common.PaginationParams{
			MaxItems: common.MaxBatchSize,
		},
	}, cfg, nil)

	for file := range dirBucketFiles {
		allBucketFiles["/"+file.Path()] = true
	}
}

func createObjectSyncFilePairProcessor(
	cfg common.Config,
	source mgcSchemaPkg.URI,
	destination mgcSchemaPkg.URI,
	progressReporter *progress_report.UnitsReporter,
	basePath mgcSchemaPkg.URI,
) pipeline.Processor[pipeline.WalkDirEntry, syncUploadPair] {
	return func(ctx context.Context, entry pipeline.WalkDirEntry) (syncUploadPair, pipeline.ProcessStatus) {
		if err := entry.Err(); err != nil {
			return syncUploadPair{}, pipeline.ProcessSkip
		}
		if entry.DirEntry().IsDir() {
			return syncUploadPair{}, pipeline.ProcessSkip
		}
		normalizedSource, err := common.GetAbsSystemURI(mgcSchemaPkg.URI(entry.Path()))
		if err != nil {
			logger().Debugw("error with path", "error", err)
			return syncUploadPair{}, pipeline.ProcessSkip
		}

		pathWithFolder := strings.TrimPrefix(entry.Path(), basePath.String())
		normalizedDestination := destination.JoinPath(pathWithFolder)

		info, err := entry.DirEntry().Info()
		if err != nil {
			return syncUploadPair{}, pipeline.ProcessAbort
		}

		progressReporter.Report(0, 1, nil)

		if allBucketFiles[pathWithFolder] {
			delete(allBucketFiles, pathWithFolder)
		}

		return syncUploadPair{
			Source:      normalizedSource,
			Destination: normalizedDestination,
			Stats: fileSyncStats{
				SourceLength:  info.Size(),
				SourceModTime: info.ModTime().Unix(),
			},
		}, pipeline.ProcessOutput
	}
}

func createSyncObjectProcessor(
	cfg common.Config,
	progressReporter *progress_report.UnitsReporter,
) pipeline.Processor[syncUploadPair, error] {
	return func(ctx context.Context, entry syncUploadPair) (error, pipeline.ProcessStatus) {
		var err error
		defer func(cause error) {
			progressReporter.Report(1, 0, err)
		}(err)

		logger().Debug("%s %s\n", entry.Source, entry.Destination)

		fileStats, err := getFileStats(ctx, entry.Destination, cfg)
		if err != nil {
			logger().Debugw("error getting file stats", "error", err)
		}

		isSameSize := entry.Stats.SourceLength == fileStats.SourceLength
		isLocalOlderThenBucket := entry.Stats.SourceModTime < fileStats.SourceModTime
		if err == nil && isSameSize && isLocalOlderThenBucket {
			logger().Debug("Skipping file [%s] - no change", entry.Source)
			return nil, pipeline.ProcessSkip
		}

		err = uploadFile(ctx, entry.Source, entry.Destination, cfg)
		if err != nil {
			return &common.ObjectError{Url: mgcSchemaPkg.URI(entry.Source.Path()), Err: err}, pipeline.ProcessOutput
		}

		uploadFiles.Increment()
		return nil, pipeline.ProcessOutput
	}
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
	return c.v
}
