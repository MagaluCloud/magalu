package objects

import (
	"context"
	"fmt"
	"io/fs"
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

	err = processCurrentAndSubfolders(ctx, cfg, params.Destination, params.StorageClass, basePath.String(), progressReporter)

	if err != nil {
		return &uploadDirResult{}, err
	}

	return &uploadDirResult{
		URI: params.Destination.String(),
		Dir: basePath.String(),
	}, nil
}

func walkFiles(ctx context.Context, root string, progressReporter *progress_report.UnitsReporter) (<-chan FileInfo, <-chan error) {
	files := make(chan FileInfo)
	errc := make(chan error, 1)

	go func() {
		defer close(files)
		errc <- filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				fileCount, err := getFileCount(path)
				if err != nil {
					return err
				}

				progressReporter.Report(0, fileCount, err)
				return nil
			}
			select {
			case files <- FileInfo{path}:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	}()
	return files, errc
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

func processCurrentAndSubfolders(ctx context.Context, cfg common.Config, destination mgcSchemaPkg.URI, storageClass string, path string, progressReporter *progress_report.UnitsReporter) error {
	files, errc := walkFiles(ctx, path, progressReporter)

	results := make(chan error)

	var wg syncer.WaitGroup
	wg.Add(cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		go func() {
			defer wg.Done()
			worker(ctx, cfg, destination, path, storageClass, files, results, progressReporter)
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for {
		select {
		case err, ok := <-results:
			if !ok {
				return nil
			}
			if err != nil {
				return err
			}
		case err := <-errc:
			if err != nil {
				return err
			}
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func getFileCount(dirPath string) (count uint64, err error) {
	i := 0
	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		defer func() { i += 1 }()
		if err != nil {
			return err
		}

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
