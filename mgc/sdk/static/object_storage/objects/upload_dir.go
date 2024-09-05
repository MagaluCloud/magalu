package objects

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	syncer "sync"

	"github.com/pterm/pterm"
	"magalu.cloud/core"
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

	files, err := walkDir(ctx, basePath.String(), params.Shallow)
	if err != nil {
		return nil, err
	}

	totalFiles := len(files)
	progressBar, _ := pterm.DefaultProgressbar.
		WithTotal(totalFiles).
		WithTitle("Uploading files").
		Start()

	err = processCurrentAndSubfolders(ctx, cfg, params.Destination, params.StorageClass, basePath.String(), files, progressBar)

	if err != nil {
		return &uploadDirResult{}, err
	}

	_, _ = progressBar.Stop()

	return &uploadDirResult{
		URI: params.Destination.String(),
		Dir: basePath.String(),
	}, nil
}

func processFile(ctx context.Context, cfg common.Config, destination mgcSchemaPkg.URI, basePath string, storageClass string, file string, progressBar *pterm.ProgressbarPrinter) error {
	dst := destination.JoinPath((strings.TrimPrefix(file, basePath)))

	_, err := upload(
		ctx,
		uploadParams{Source: mgcSchemaPkg.FilePath(file), Destination: dst, StorageClass: storageClass},
		cfg,
	)

	if err != nil {
		err = &common.ObjectError{Url: mgcSchemaPkg.URI(dst), Err: err}
		return err
	}

	progressBar.Increment()
	return nil
}

func worker(ctx context.Context, cfg common.Config, destination mgcSchemaPkg.URI, basePath string, storageClass string, files <-chan string, results chan<- error, progressBar *pterm.ProgressbarPrinter) {
	for {
		select {
		case file, ok := <-files:
			if !ok {
				return
			}
			err := processFile(ctx, cfg, destination, basePath, storageClass, file, progressBar)
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

func processCurrentAndSubfolders(ctx context.Context, cfg common.Config, destination mgcSchemaPkg.URI, storageClass string, path string, files []string, progressBar *pterm.ProgressbarPrinter) error {
	results := make(chan error, cfg.Workers)
	filesChan := make(chan string, cfg.Workers)

	var wg syncer.WaitGroup
	wg.Add(cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		go func() {
			defer wg.Done()
			worker(ctx, cfg, destination, path, storageClass, filesChan, results, progressBar)
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

	for err := range results {
		if err != nil {
			return err
		}
	}

	return nil
}

func walkDir(ctx context.Context, root string, shallow bool) ([]string, error) {
	var files []string

	var walkFn func(string) error
	walkFn = func(dir string) error {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			path := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				if !shallow {
					if err := walkFn(path); err != nil {
						return err
					}
				}
			} else {
				files = append(files, path)
			}
		}
		return nil
	}

	if err := walkFn(root); err != nil {
		return nil, err
	}
	return files, nil
}
