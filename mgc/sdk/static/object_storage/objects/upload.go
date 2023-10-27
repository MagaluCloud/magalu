package objects

import (
	"context"
	"path"
	"strings"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type uploadParams struct {
	Destination string `json:"dst" jsonschema:"description=Full destination path in the bucket with s3 prefix, i.e: s3://<bucket-name>"`
	Source      string `json:"src" jsonschema:"description=Source file path to be uploaded"`
}

type uploadTemplateResult struct {
	File string `json:"file"`
	URI  string `json:"uri"`
}

var getUpload = utils.NewLazyLoader[core.Executor](newUpload)

func newUpload() core.Executor {
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
}

func formatURI(uri string) string {
	if !strings.Contains(uri, common.URIPrefix) {
		return common.URIPrefix + uri
	}
	return uri
}

func upload(ctx context.Context, params uploadParams, cfg common.Config) (*uploadTemplateResult, error) {
	dst, _ := strings.CutPrefix(params.Destination, common.URIPrefix)
	_, fileName := path.Split(params.Source)
	if isDirPath(dst) {
		// If it isn't a file path, don't rename, just append source with bucket URI
		dst = path.Join(dst, fileName)
	}

	uploader, err := common.NewUploader(cfg, params.Source, dst)
	if err != nil {
		return nil, err
	}

	if err = uploader.Upload(ctx); err != nil {
		return nil, err
	}

	return &uploadTemplateResult{
		// path.Join will remove URI double slash prefix on s3://
		URI:  formatURI(dst),
		File: fileName,
	}, nil
}
