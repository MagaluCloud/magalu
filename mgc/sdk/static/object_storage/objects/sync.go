package objects

import (
	"context"
	"io/fs"
	"path/filepath"

	"magalu.cloud/core"
	"magalu.cloud/core/pipeline"
	"magalu.cloud/core/progress_report"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type syncParams struct {
	Source      mgcSchemaPkg.FilePath `json:"src" jsonschema:"description=Source path to sync the remote with,example=./" mgc:"positional"`
	Destination mgcSchemaPkg.URI      `json:"dst" jsonschema:"description=Full destination path to sync with the source path,example=s3://my-bucket/dir/" mgc:"positional"`
	Delete      bool                  `json:"delete,omitempty" jsonschema:"description=Deletes any item at the destination not present on the source,default=false"`
	BatchSize   int                   `json:"batch_size,omitempty" jsonschema:"description=Limit of items per batch to delete,default=1000,minimum=1,maximum=1000" example:"1000"`
}

type syncResult struct {
	Source        mgcSchemaPkg.FilePath `json:"src" jsonschema:"description=Source path to sync the remote with,example=./" mgc:"positional"`
	Destination   mgcSchemaPkg.URI      `json:"dst" jsonschema:"description=Full destination path to sync with the source path,example=s3://my-bucket/dir/" mgc:"positional"`
	FilesDeleted  int                   `json:"deleted"`
	FilesUploaded int                   `json:"uploaded"`
	Deleted       bool                  `json:"hasDeleted"`
}

var getSync = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "sync",
			Summary:     "Synchronizes a local path to a bucket",
			Description: "This command uploads any file from the source to the destination if it's not present or has a different size. Additionally any file in the destination not present on the source is deleted.",
		},
		sync,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template={{if and (eq .deleted 0) (eq .uploaded 0)}}Already Synced{{- else}}" +
			"Synced files from {{.src}} to {{.dst}} {{.uploaded}} uploaded and {{if .hasDeleted}}{{.deleted}} deleted{{- else}}{{.deleted}} to be deleted with the parameter --delete{{- end}} {{- end}}\n"
	})
})

func createObjectSyncProcessor(cfg common.Config, destination mgcSchemaPkg.URI, srcAbs string) pipeline.Processor[pipeline.WalkDirEntry, error] {
	return func(ctx context.Context, dirEntry pipeline.WalkDirEntry) (error, pipeline.ProcessStatus) {
		if err := dirEntry.Err(); err != nil {
			return &common.ObjectError{Err: err}, pipeline.ProcessAbort
		}

		if dirEntry.DirEntry().IsDir() {
			return nil, pipeline.ProcessOutput
		}

		absEntry, err := filepath.Abs(dirEntry.Path())
		if err != nil {
			return &common.ObjectError{Err: err}, pipeline.ProcessAbort
		}
		relative, err := filepath.Rel(srcAbs, absEntry)
		if err != nil {
			return &common.ObjectError{Err: err}, pipeline.ProcessAbort
		}
		objURI := destination.JoinPath(relative)

		_, err = upload(
			ctx,
			uploadParams{Source: mgcSchemaPkg.FilePath(dirEntry.Path()), Destination: mgcSchemaPkg.URI(objURI)},
			cfg,
		)

		if err != nil {
			return &common.ObjectError{Url: mgcSchemaPkg.URI(objURI), Err: err}, pipeline.ProcessOutput
		}

		return nil, pipeline.ProcessOutput
	}
}

func sync(ctx context.Context, params syncParams, cfg common.Config) (result core.Value, err error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	srcObjects := pipeline.WalkDirEntries(ctx, params.Source.String(), func(path string, d fs.DirEntry, err error) error {
		return err
	})
	srcAbs, err := filepath.Abs(params.Source.String())
	if err != nil {
		return nil, err
	}

	dstObjects := common.ListGenerator(ctx, common.ListObjectsParams{
		Destination: params.Destination,
		Recursive:   true,
		PaginationParams: common.PaginationParams{
			MaxItems: common.MaxBatchSize,
		},
	}, cfg, nil)

	srcMap := make(map[string]pipeline.WalkDirEntry)
	var toUpload []pipeline.WalkDirEntry
	var toDelete []pipeline.WalkDirEntry
	var objErr utils.MultiError

	for entry := range srcObjects {
		if entry.Err() != nil {
			objErr = append(objErr, err)
			continue
		}
		if entry.DirEntry().IsDir() {
			continue
		}
		absPath, err := filepath.Abs(entry.Path())
		if err != nil {
			objErr = append(objErr, err)
			continue
		}
		relPath, err := filepath.Rel(srcAbs, absPath)
		if err != nil {
			objErr = append(objErr, err)
			continue
		}
		srcMap[relPath] = entry
	}

	for entry := range dstObjects {
		if entry.Err() != nil {
			objErr = append(objErr, err)
			continue
		}
		if entry.DirEntry().IsDir() {
			continue
		}
		relPath, err := filepath.Rel(params.Destination.Path(), entry.Path())
		if err != nil {
			objErr = append(objErr, err)
			continue
		}
		value, err := entry.DirEntry().Info()
		if err != nil {
			objErr = append(objErr, err)
			continue
		}
		srcValue, ok := srcMap[relPath]

		if !ok {
			toDelete = append(toDelete, entry)
			continue
		}
		srcInfo, _ := srcValue.DirEntry().Info()
		if value == nil || srcInfo.Size() != value.Size() {
			toUpload = append(toUpload, srcValue)
			continue
		}
		delete(srcMap, relPath)
	}
	if len(objErr) > 0 {
		return nil, objErr
	}

	for _, v := range srcMap {
		if v != nil {
			toUpload = append(toUpload, v)
		}
	}

	if len(toUpload) > 0 {
		reportProgress := progress_report.FromContext(ctx)
		total := uint64(len(toUpload))

		reportProgress("Sync Upload", 0, total, progress_report.UnitsNone, nil)
		uploadChannel := pipeline.SliceItemGenerator[pipeline.WalkDirEntry](ctx, toUpload)
		uploadObjectsErrorChan := pipeline.ParallelProcess(ctx, cfg.Workers, uploadChannel, createObjectSyncProcessor(cfg, params.Destination, srcAbs), nil)
		uploadObjectsErrorChan = pipeline.Filter(ctx, uploadObjectsErrorChan, pipeline.FilterNonNil[error]{})

		objErr, _ = pipeline.SliceItemConsumer[utils.MultiError](ctx, uploadObjectsErrorChan)
		reportProgress("Sync Upload", total, total, progress_report.UnitsNone, nil)
		if len(objErr) > 0 {
			return nil, objErr
		}
	}

	if params.Delete && len(toDelete) > 0 {
		deleteChannel := pipeline.SliceItemGenerator[pipeline.WalkDirEntry](ctx, toDelete)
		err := common.DeleteObjects(ctx, common.DeleteObjectsParams{
			Destination: params.Destination,
			ToDelete:    deleteChannel,
			BatchSize:   params.BatchSize,
		}, cfg)
		if err != nil {
			return nil, err
		}
	}

	return syncResult{
		FilesDeleted:  len(toDelete),
		FilesUploaded: len(toUpload),
		Source:        params.Source,
		Destination:   params.Destination,
		Deleted:       params.Delete,
	}, nil
}
