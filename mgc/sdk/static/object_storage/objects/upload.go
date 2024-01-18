package objects

import (
	"context"
	"fmt"
	"strings"

	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type uploadParams struct {
	Source      mgcSchemaPkg.FilePath `json:"src" jsonschema:"description=Source file path to be uploaded,example=./file.txt" mgc:"positional"`
	Destination mgcSchemaPkg.URI      `json:"dst" jsonschema:"description=Full destination path in the bucket with desired filename,example=s3://my-bucket/dir/file.txt" mgc:"positional"`
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

	fileName := params.Source.String()
	if strings.HasSuffix(fullDstPath.String(), "/") {
		// If it isn't a file path, don't rename, just append source with bucket URI
		fullDstPath = fullDstPath.JoinPath(fileName)
	}
	if params.Destination.String() == "s3://"+params.Destination.Hostname() {
		fullDstPath = fullDstPath.JoinPath(fileName)
	}

	uploader, err := common.NewUploader(cfg, params.Source, fullDstPath)
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
