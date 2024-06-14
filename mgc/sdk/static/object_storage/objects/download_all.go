package objects

import (
	"context"
	"fmt"
	"math"

	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/core/pipeline"
	"magalu.cloud/core/progress_report"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

var downloadAllObjectsLogger *zap.SugaredLogger

func downloadAllLogger() *zap.SugaredLogger {
	if downloadAllObjectsLogger == nil {
		downloadAllObjectsLogger = logger().Named("download")
	}
	return downloadAllObjectsLogger
}

type downloadAllObjectsParams struct {
	Source         mgcSchemaPkg.URI      `json:"src" jsonschema:"description=Path of objects to be downloaded,example=mybucket" mgc:"positional"`
	Destination    mgcSchemaPkg.FilePath `json:"dst,omitempty" jsonschema:"description=Path to save files,example=path/to/folder" mgc:"positional"`
	common.Filters `json:",squash"`      // nolint
}

var getDownloadAll = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "download-all",
			Description: "Download all objects from a bucket",
		},
		downloadAll,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Downloaded from {{.src}} to {{.dst}}\n"
	})
})

func createObjectDownloadProcessor(
	cfg common.Config,
	params downloadAllObjectsParams,
	progressReporter *progress_report.UnitsReporter,
) pipeline.Processor[pipeline.WalkDirEntry, error] {
	return func(ctx context.Context, dirEntry pipeline.WalkDirEntry) (error, pipeline.ProcessStatus) {
		bucketName := common.NewBucketNameFromURI(params.Source)
		rootURI := bucketName.AsURI()
		var err error

		defer func() { progressReporter.Report(1, 0, err) }()

		objURI := rootURI.JoinPath(dirEntry.Path())

		if dirEntry.Err() != nil {
			err = &common.ObjectError{Url: mgcSchemaPkg.URI(objURI), Err: dirEntry.Err()}
			return err, pipeline.ProcessOutput
		}

		_, ok := dirEntry.DirEntry().(*common.BucketContent)
		if !ok {
			err = &common.ObjectError{Url: mgcSchemaPkg.URI(objURI), Err: fmt.Errorf("expected object, got directory")}
			return err, pipeline.ProcessOutput
		}

		downloadAllLogger().Infow("Downloading object", "uri", objURI)
		downloader, err := common.NewDownloader(ctx, cfg, objURI, params.Destination.Join(dirEntry.Path()), "") // since we are downloading N objects, can't set a version
		if err != nil {
			return err, pipeline.ProcessAbort
		}

		if err = downloader.Download(ctx); err != nil {
			return err, pipeline.ProcessAbort
		}

		if err != nil {
			err = &common.ObjectError{Url: mgcSchemaPkg.URI(objURI), Err: err}
			return err, pipeline.ProcessOutput
		}

		return nil, pipeline.ProcessOutput
	}
}

func downloadMultipleFiles(ctx context.Context, cfg common.Config, params downloadAllObjectsParams) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	listParams := common.ListObjectsParams{
		Destination: params.Source,
		Recursive:   true,
		PaginationParams: common.PaginationParams{
			MaxItems: math.MaxInt64,
		},
	}

	progressReportMsg := "Downloading objects from: " + params.Source.String()
	progressReporter := progress_report.NewUnitsReporter(ctx, progressReportMsg, 0)
	progressReporter.Start()
	defer progressReporter.End()

	onNewPage := func(objCount uint64) {
		progressReporter.Report(0, objCount, nil)
	}

	objs := common.ListGenerator(ctx, listParams, cfg, onNewPage)
	objs = common.ApplyFilters(ctx, objs, params.FilterParams, cancel)

	downloadObjectsErrorChan := pipeline.ParallelProcess(ctx, cfg.Workers, objs, createObjectDownloadProcessor(cfg, params, progressReporter), nil)
	downloadObjectsErrorChan = pipeline.Filter(ctx, downloadObjectsErrorChan, pipeline.FilterNonNil[error]{})

	objErr, err := pipeline.SliceItemConsumer[utils.MultiError](ctx, downloadObjectsErrorChan)
	if err != nil {
		return err
	}
	if len(objErr) > 0 {
		return objErr
	}

	return nil
}

func downloadAll(ctx context.Context, p downloadAllObjectsParams, cfg common.Config) (result common.DownloadObjectParams, err error) {
	p.Destination, err = common.GetDownloadFileDst(p.Destination, p.Source)
	if err != nil {
		return result, fmt.Errorf("no destination specified and could not use local dir: %w", err)
	}
	err = downloadMultipleFiles(ctx, cfg, p)

	if err != nil {
		return result, err
	}

	return common.DownloadObjectParams{Source: p.Source, Destination: p.Destination}, nil
}
