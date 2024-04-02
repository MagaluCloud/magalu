package objects

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"magalu.cloud/core"
	"magalu.cloud/core/pipeline"
	"magalu.cloud/core/progress_report"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type syncUploadPair struct {
	Source      mgcSchemaPkg.URI
	Destination mgcSchemaPkg.URI
	Stats       fileSyncStats
}

type fileSyncStats struct {
	SourceLength  int64
	SourceModTime int64
}

type syncParams struct {
	Source      mgcSchemaPkg.URI `json:"src" jsonschema:"description=Source path to sync the remote with,example=./" mgc:"positional"`
	Destination mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Full destination path to sync with the source path,example=s3://my-bucket/dir/" mgc:"positional"`
	Delete      bool             `json:"delete,omitempty" jsonschema:"description=Deletes any item at the destination not present on the source,default=false"`
	BatchSize   int              `json:"batch_size,omitempty" jsonschema:"description=Limit of items per batch to delete,default=1000,minimum=1,maximum=1000" example:"1000"`
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
			Summary:     "Synchronizes a local path to a bucket",
			Description: "This command uploads any file from the source to the destination if it's not present or has a different size. Additionally any file in the destination not present on the source is deleted.",
			// Scopes:      core.Scopes{"object-storage.read", "object-storage.write"},
		},
		sync,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template={{if and (eq .deleted 0) (eq .uploaded 0)}}Already Synced{{- else}}" +
			"Synced files from {{.src}} to {{.dst}}\n- {{.uploaded}} uploaded\n- {{if .hasDeleted}}{{.deleted}} deleted\n\nDeleted files:\n-{{.deletedFiles}}{{- else}}{{.deleted}} to be deleted with the parameter --delete{{- end}}{{- end}}\n"
	})
})

func sync(ctx context.Context, params syncParams, cfg common.Config) (result core.Value, err error) {
	if !strings.HasPrefix(string(params.Source), common.URIPrefix) && !strings.HasPrefix(string(params.Destination), common.URIPrefix) {
		return nil, fmt.Errorf("enter at least one parameter using the standard prefix \"s3://\"")
	}

	// TODO GA - Improve and remove this lock
	if strings.HasPrefix(string(params.Source), common.URIPrefix) && strings.HasPrefix(string(params.Destination), common.URIPrefix) {
		return nil, fmt.Errorf("to copy or move between buckets, use \"mgc object-storage objects copy/move\"")
	}

	srcIsRemote := isRemote(params.Source)

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	var srcObjects <-chan pipeline.WalkDirEntry

	if srcIsRemote {
		srcObjects = common.ListGenerator(ctx, common.ListObjectsParams{
			Destination: params.Source,
			Recursive:   true,
			PaginationParams: common.PaginationParams{
				MaxItems: common.MaxBatchSize,
			},
		}, cfg, nil)
	} else {
		srcObjects = pipeline.WalkDirEntries(ctx, params.Source.String(), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			return nil
		})
	}

	progressReporter := progress_report.NewUnitsReporter(ctx, "Sync Download", 0)
	progressReporter.Start()

	uploadChannel := pipeline.Process(ctx, srcObjects, createObjectSyncFilePairProcessor(cfg, params.Source, params.Destination, progressReporter), nil)
	uploadObjectsErrorChan := pipeline.ParallelProcess(ctx, cfg.Workers, uploadChannel, createSyncObjectProcessor(cfg, progressReporter), nil)
	objErr, err := pipeline.SliceItemConsumer[utils.MultiError](ctx, uploadObjectsErrorChan)

	for _, er := range objErr {
		if er != nil {
			progressReporter.Report(0, 0, objErr)
			return nil, objErr
		}
	}

	progressReporter.End()
	return syncResult{
		Source:        params.Source,
		Destination:   params.Destination,
		FilesDeleted:  0,
		FilesUploaded: 0,
		Deleted:       params.Delete,
		DeletedFiles:  "",
	}, nil
}

func createObjectSyncFilePairProcessor(
	cfg common.Config,
	source mgcSchemaPkg.URI,
	destination mgcSchemaPkg.URI,
	progressReporter *progress_report.UnitsReporter,
) pipeline.Processor[pipeline.WalkDirEntry, syncUploadPair] {
	return func(ctx context.Context, entry pipeline.WalkDirEntry) (syncUploadPair, pipeline.ProcessStatus) {
		if err := entry.Err(); err != nil {
			return syncUploadPair{}, pipeline.ProcessSkip
		}
		if entry.DirEntry().IsDir() {
			return syncUploadPair{}, pipeline.ProcessSkip
		}

		normalizedSource, err := normalizeURI(source, entry.Path())
		if err != nil {
			return syncUploadPair{}, pipeline.ProcessSkip
		}

		normalizedDestination, err := normalizeURI(destination, entry.Path())
		if err != nil {
			return syncUploadPair{}, pipeline.ProcessSkip
		}

		info, err := entry.DirEntry().Info()
		if err != nil {
			return syncUploadPair{}, pipeline.ProcessAbort
		}

		progressReporter.Report(0, 1, nil)

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

func normalizeURI(uri mgcSchemaPkg.URI, path string) (mgcSchemaPkg.URI, error) {
	if uri.Scheme() == "s3" {
		return uri.JoinPath(filepath.Base(path)), nil
	}

	if uri == "." {
		value, err := filepath.Abs(uri.Path())
		if err != nil {
			return uri, err
		}

		return mgcSchemaPkg.FilePath(value).AsURI().JoinPath(filepath.Base(path)), nil
	}
	return uri.JoinPath(filepath.Base(path)), nil
}

func createSyncObjectProcessor(
	cfg common.Config,
	progressReporter *progress_report.UnitsReporter,
) pipeline.Processor[syncUploadPair, error] {
	return func(ctx context.Context, entry syncUploadPair) (error, pipeline.ProcessStatus) {
		var err error
		defer func(cause error) { progressReporter.Report(1, 0, err) }(err)

		logger().Debug("%s %s\n", entry.Source, entry.Destination)

		fileStats, err := getFileStats(ctx, entry.Destination, cfg)

		if err == nil && entry.Stats.SourceLength == fileStats.SourceLength && entry.Stats.SourceModTime == fileStats.SourceModTime {
			return nil, pipeline.ProcessSkip // TODO
		}

		err = sourceDestinationProcessor(ctx, entry.Source, entry.Destination, cfg)

		if err != nil {
			return &common.ObjectError{Url: mgcSchemaPkg.URI(entry.Source.Path()), Err: err}, pipeline.ProcessOutput
		}

		return nil, pipeline.ProcessOutput
	}
}

func getFileStats(ctx context.Context, destination mgcSchemaPkg.URI, cfg common.Config) (fileSyncStats, error) {
	if isRemote(destination) {
		dstHead, err := headObject(ctx, headObjectParams{
			Destination: destination,
		}, cfg)
		if err != nil {
			return fileSyncStats{}, err
		}
		dstModTime, err := time.Parse(time.RFC3339, dstHead.LastModified)
		if err != nil {
			logger().Debug("%s %s\n", dstModTime, err)
			return fileSyncStats{}, err
		}
		return fileSyncStats{
			SourceLength:  dstHead.ContentLength,
			SourceModTime: dstModTime.Unix(),
		}, nil
	}

	stat, err := os.Stat(string(destination))
	if err != nil {
		logger().Debug("%s %s\n", stat, err)
		return fileSyncStats{}, err
	}
	return fileSyncStats{
		SourceLength:  stat.Size(),
		SourceModTime: stat.ModTime().Unix(),
	}, nil
}

func sourceDestinationProcessor(ctx context.Context, source mgcSchemaPkg.URI, destination mgcSchemaPkg.URI, cfg common.Config) error {
	if isRemote(destination) {
		sourcePath := mgcSchemaPkg.FilePath("/" + source.Path())
		_, err := upload(
			ctx,
			uploadParams{Source: sourcePath, Destination: destination},
			cfg,
		)
		return err
	} else {
		destinationPath := mgcSchemaPkg.FilePath("/" + destination.Path())
		_, err := download(
			ctx,
			common.DownloadObjectParams{Source: source, Destination: destinationPath},
			cfg,
		)
		return err
	}
}
