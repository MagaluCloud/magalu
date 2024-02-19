package objects

import (
	"context"
	"fmt"
	"io/fs"

	"magalu.cloud/core"
	"magalu.cloud/core/pipeline"
	"magalu.cloud/core/progress_report"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type uploadDirParams struct {
	Source         mgcSchemaPkg.DirPath `json:"src" jsonschema:"description=Source directory path for upload,example=path/to/folder" mgc:"positional"`
	Destination    mgcSchemaPkg.URI     `json:"dst" jsonschema:"description=Full destination path in the bucket,example=my-bucket/dir/" mgc:"positional"`
	common.Filters `json:",squash"`     // nolint
}

type uploadDirResult struct {
	Dir string `json:"dir"`
	URI string `json:"uri"`
}

var getUploadDir = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "upload-dir",
			Description: "Upload a directory to a bucket",
		},
		uploadDir,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Uploaded directory {{.dir}} to {{.uri}}\n"
	})
})

func createObjectUploadProcessor(cfg common.Config, destination mgcSchemaPkg.URI) pipeline.Processor[pipeline.WalkDirEntry, error] {
	return func(ctx context.Context, dirEntry pipeline.WalkDirEntry) (error, pipeline.ProcessStatus) {
		if err := dirEntry.Err(); err != nil {
			return &common.ObjectError{Err: err}, pipeline.ProcessAbort
		}

		if dirEntry.DirEntry().IsDir() {
			return nil, pipeline.ProcessOutput
		}

		filePath := dirEntry.Path()
		objURI := destination.JoinPath(filePath)

		_, err := upload(
			ctx,
			uploadParams{Source: mgcSchemaPkg.FilePath(filePath), Destination: mgcSchemaPkg.URI(objURI)},
			cfg,
		)

		if err != nil {
			return &common.ObjectError{Url: mgcSchemaPkg.URI(objURI), Err: err}, pipeline.ProcessOutput
		}

		return nil, pipeline.ProcessOutput
	}
}

func uploadDir(ctx context.Context, params uploadDirParams, cfg common.Config) (*uploadDirResult, error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	if params.Source.String() == "" {
		return nil, core.UsageError{Err: fmt.Errorf("source cannot be empty")}
	}

	progressReportMsg := "Uploading directory: " + params.Source.String()
	progressReporter := progress_report.NewUnitsReporter(ctx, progressReportMsg, 1)
	progressReporter.Start()
	defer progressReporter.End()

	entries := pipeline.WalkDirEntries(ctx, params.Source.String(), func(path string, d fs.DirEntry, err error) error {
		return err
	})

	entries = common.ApplyFilters(ctx, entries, params.FilterParams, cancel)
	uploadObjectsErrorChan := pipeline.ParallelProcess(ctx, cfg.Workers, entries, createObjectUploadProcessor(cfg, params.Destination), nil)
	uploadObjectsErrorChan = pipeline.Filter(ctx, uploadObjectsErrorChan, pipeline.FilterNonNil[error]{})

	objErr, err := pipeline.SliceItemConsumer[utils.MultiError](ctx, uploadObjectsErrorChan)
	if err != nil {
		progressReporter.Report(0, 0, err)
		return nil, err
	}
	if len(objErr) > 0 {
		progressReporter.Report(0, 0, objErr)
		return nil, objErr
	}

	progressReporter.Report(1, 1, nil)

	return &uploadDirResult{
		URI: params.Destination.String(),
		Dir: params.Source.String(),
	}, nil
}
