package versioning

import (
	"context"
	"net/http"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type getBucketVersioningParams struct {
	Bucket common.BucketName `json:"bucket" jsonschema:"description=Bucket name to get versioning info from" mgc:"positional"`
}

var getGet = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get versioning info for a Bucket",
		},
		getBucketVersioning,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "table"
	})
	return exec
})

func getBucketVersioning(ctx context.Context, params getBucketVersioningParams, cfg common.Config) (result versioningConfiguration, err error) {
	req, err := newGetBucketVersioningRequest(ctx, params.Bucket, cfg)
	if err != nil {
		return
	}

	res, err := common.SendRequest(ctx, req, cfg)
	if err != nil {
		return
	}

	return common.UnwrapResponse[versioningConfiguration](res, req)
}

func newGetBucketVersioningRequest(ctx context.Context, bucketName common.BucketName, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Set("versioning", "")

	url.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
