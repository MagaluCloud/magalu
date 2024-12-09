package object_lock

import (
	"context"
	"net/http"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type GetBucketObjectLockParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Specifies the bucket whose ACL is being requested" mgc:"positional"`
}

type objectLockingResponse struct {
	ObjectLockEnabled string
	Rule              common.ObjectLockRule
}

var getGet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get object locking configuration for the specified bucket",
		},
		getObjectLocking,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "json"
	})
	return exec
})

func getObjectLocking(ctx context.Context, params GetBucketObjectLockParams, cfg common.Config) (result objectLockingResponse, err error) {
	req, err := newGetObjectLockingRequest(ctx, cfg, params.Bucket)
	if err != nil {
		return
	}

	res, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	result, err = common.UnwrapResponse[objectLockingResponse](res, req)
	return
}

func newGetObjectLockingRequest(ctx context.Context, cfg common.Config, bucketName common.BucketName) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("object-lock", "")
	url.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
