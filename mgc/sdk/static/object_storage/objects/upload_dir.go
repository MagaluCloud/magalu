package objects

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	syncer "sync"

	"magalu.cloud/core"
	"magalu.cloud/core/progress_report"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type uploadDirParams struct {
	Source         mgcSchemaPkg.DirPath `json:"src" jsonschema:"description=Source directory path for upload,example=path/to/folder" mgc:"positional"`
	Destination    mgcSchemaPkg.URI     `json:"dst" jsonschema:"description=Full destination path in the bucket,example=my-bucket/dir/" mgc:"positional"`
	Shallow        bool                 `json:"shallow,omitempty" jsonschema:"description=Don't upload subdirectories,default=false"`
	StorageClass   string               `json:"storage_class,omitempty" jsonschema:"description=Type of Storage in which to store object,example=cold,enum=,enum=standard,enum=cold,enum=glacier_ir,enum=cold_instant,default="`
	common.Filters `json:",squash"`     // nolint
}

type uploadDirResult struct {
	Dir string `json:"dir"`
	URI string `json:"uri"`
}

type FileInfo struct {
	Path string
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

func uploadDir(ctx context.Context, params uploadDirParams, cfg common.Config) (*uploadDirResult, error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	if params.Source.String() == "" {
		return nil, core.UsageError{Err: fmt.Errorf("source cannot be empty")}
	}

	basePath, err := common.GetAbsSystemURI(mgcSchemaPkg.URI(params.Source.String()))
	if err != nil {
		return nil, err
	}

	progressReportMsg := "Uploading directory: " + basePath.String()
	progressReporter := progress_report.NewUnitsReporter(ctx, progressReportMsg, 0)
	progressReporter.Start()
	defer progressReporter.End()

	err = processCurrentAndSubfolders(ctx, cfg, params.Destination, params.StorageClass, basePath.String(), params.Shallow, progressReporter)

	if err != nil {
		return &uploadDirResult{}, err
	}

	return &uploadDirResult{
		URI: params.Destination.String(),
		Dir: basePath.String(),
	}, nil
}

func processFile(ctx context.Context, cfg common.Config, destination mgcSchemaPkg.URI, basePath string, storageClass string, file FileInfo, progressReporter *progress_report.UnitsReporter) error {
	filePath := file.Path
	dst := destination.JoinPath((strings.TrimPrefix(filePath, basePath)))

	_, err := upload(
		ctx,
		uploadParams{Source: mgcSchemaPkg.FilePath(filePath), Destination: dst, StorageClass: storageClass},
		cfg,
	)

	if err != nil {
		err = &common.ObjectError{Url: mgcSchemaPkg.URI(dst), Err: err}
		progressReporter.Report(1, 0, err)
		return err
	}
	return nil
}

func worker(ctx context.Context, cfg common.Config, destination mgcSchemaPkg.URI, basePath string, storageClass string, files <-chan FileInfo, results chan<- error, progressReporter *progress_report.UnitsReporter) {
	for {
		select {
		case file, ok := <-files:
			if !ok {
				return
			}
			err := processFile(ctx, cfg, destination, basePath, storageClass, file, progressReporter)
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

func processCurrentAndSubfolders(ctx context.Context, cfg common.Config, destination mgcSchemaPkg.URI, storageClass string, path string, shallow bool, progressReporter *progress_report.UnitsReporter) error {
	files := make(chan FileInfo, cfg.Workers)
	results := make(chan error, cfg.Workers)

	var wg syncer.WaitGroup
	wg.Add(cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		go func() {
			defer wg.Done()
			worker(ctx, cfg, destination, path, storageClass, files, results, progressReporter)
		}()
	}

	go func() {
		defer close(files)
		err := walkDir(ctx, path, files, shallow, progressReporter)
		if err != nil {
			select {
			case results <- err:
			case <-ctx.Done():
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for err := range results {
		if err != nil {
			return err
		}
	}

	return nil
}

func walkDir(ctx context.Context, root string, files chan<- FileInfo, shallow bool, progressReporter *progress_report.UnitsReporter) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	fileCount := uint64(0)
	for _, entry := range entries {
		if entry.IsDir() {
			if !shallow {
				subdir := filepath.Join(root, entry.Name())
				if err := walkDir(ctx, subdir, files, shallow, progressReporter); err != nil {
					return err
				}
			}
		} else {
			fileCount++
		}
	}

	progressReporter.Report(0, fileCount, nil)

	for _, entry := range entries {
		if !entry.IsDir() {
			select {
			case files <- FileInfo{Path: filepath.Join(root, entry.Name())}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}
