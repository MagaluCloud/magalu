package buckets

import (
	"context"
	"net/http"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type getParams struct {
	BucketName common.BucketName `json:"bucket" jsonschema:"description=Name of the bucket to be created" mgc:"positional"`
}

var getBucket = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewReflectedSimpleExecutor[getParams, common.Config, *createParams](
		core.ExecutorSpec{
			DescriptorSpec: core.DescriptorSpec{
				Name:        "get",
				Description: "Get bucket",
				IsInternal:  utils.BoolPtr(true),
				// Scopes:      core.Scopes{"object-storage.read"},
			},
		},
		getValidBucket,
	)

})

func getValidBucket(ctx context.Context, params getParams, cfg common.Config) (*createParams, error) {
	req, err := newGetRequest(ctx, params.BucketName, cfg)
	if err != nil {
		return nil, err
	}

	res, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	err = common.ExtractErr(res, req)
	if err != nil {
		return nil, err
	}

	return &createParams{BucketName: params.BucketName}, nil
}

func newGetRequest(ctx context.Context, bucketName common.BucketName, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Set("versioning", "")

	url.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
