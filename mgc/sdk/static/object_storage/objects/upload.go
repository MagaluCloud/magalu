package objects

import (
	"context"
	"fmt"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type uploadParams struct {
	Source       mgcSchemaPkg.FilePath `json:"src" jsonschema:"description=Source file path to be uploaded,example=./file.txt" mgc:"positional"`
	Destination  mgcSchemaPkg.URI      `json:"dst" jsonschema:"description=Full destination path in the bucket with desired filename,example=my-bucket/dir/file.txt" mgc:"positional"`
	StorageClass string                `json:"storage_class,omitempty" jsonschema:"description=Type of Storage in which to store object,example=cold,enum=,enum=standard,enum=cold,enum=glacier_ir,enum=cold_instant,default="`
}

type uploadTemplateResult struct {
	File string `json:"file"`
	URI  string `json:"uri"`
}

var getUpload = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "upload",
			Description: "Upload a file to a bucket",
		},
		upload,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Uploaded file {{.file}} to {{.uri}}\n"
	})
})

func upload(ctx context.Context, params uploadParams, cfg common.Config) (*uploadTemplateResult, error) {
	fullDstPath := params.Destination
	if fullDstPath == "" {
		return nil, core.UsageError{Err: fmt.Errorf("destination cannot be empty")}
	}

	fileName := params.Source.AsURI().Filename()
	if params.Destination.IsRoot() || strings.HasSuffix(fullDstPath.String(), "/") {
		fullDstPath = fullDstPath.JoinPath(fileName)
	}

	uploader, err := common.NewUploader(cfg, params.Source, fullDstPath, params.StorageClass)
	if err != nil {
		return nil, err
	}

	if err = uploader.Upload(ctx); err != nil {
		return nil, err
	}

	return &uploadTemplateResult{
		URI:  fullDstPath.String(),
		File: fileName,
	}, nil
}
