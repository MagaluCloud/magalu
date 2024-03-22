package objects

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"

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
	Shallow        bool                 `json:"shallow,omitempty" jsonschema:"description=Don't upload subdirectories,default=false"`
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
			// Scopes:      core.Scopes{"object-storage.write"},
		},
		uploadDir,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Uploaded directory {{.dir}} to {{.uri}}\n"
	})
})

func createObjectUploadProcessor(cfg common.Config, destination mgcSchemaPkg.URI, progressReporter *progress_report.UnitsReporter) pipeline.Processor[pipeline.WalkDirEntry, error] {
	return func(ctx context.Context, dirEntry pipeline.WalkDirEntry) (error, pipeline.ProcessStatus) {
		var err error
		defer func() { progressReporter.Report(1, 0, err) }()

		if err = dirEntry.Err(); err != nil {
			err = &common.ObjectError{Err: err}
			return err, pipeline.ProcessAbort
		}

		if dirEntry.DirEntry().IsDir() {
			return nil, pipeline.ProcessOutput
		}

		filePath := dirEntry.Path()
		objURI := destination.JoinPath(filePath)

		_, err = upload(
			ctx,
			uploadParams{Source: mgcSchemaPkg.FilePath(filePath), Destination: mgcSchemaPkg.URI(objURI)},
			cfg,
		)

		if err != nil {
			err = &common.ObjectError{Url: mgcSchemaPkg.URI(objURI), Err: err}
			return err, pipeline.ProcessOutput
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
	progressReporter := progress_report.NewUnitsReporter(ctx, progressReportMsg, 0)
	progressReporter.Start()
	defer progressReporter.End()

	entries := pipeline.WalkDirEntries(ctx, params.Source.String(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path != params.Source.String() && path+"/" != params.Source.String() {
			if d.IsDir() && params.Shallow {
				return filepath.SkipDir
			}
		}

		if d.IsDir() {
			fileCount, err := getFileCount(path)
			if err != nil {
				return err
			}

			progressReporter.Report(0, fileCount, err)
		}

		return nil
	})

	entries = common.ApplyFilters(ctx, entries, params.FilterParams, cancel)
	uploadObjectsErrorChan := pipeline.ParallelProcess(ctx, cfg.Workers, entries, createObjectUploadProcessor(cfg, params.Destination, progressReporter), nil)
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

func getFileCount(dirPath string) (count uint64, err error) {
	i := 0
	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		defer func() { i += 1 }()
		if err != nil {
			return err
		}

		// First loop will always be the dir represented by 'dirPath' itself, so skip it
		// bud don't return 'fs.SkipDir'
		if i == 0 {
			return nil
		}

		if d.IsDir() {
			return fs.SkipDir
		}

		count += 1
		return nil
	})

	return
}
