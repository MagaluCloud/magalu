package acl

import (
	"context"
	"net/http"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type GetBucketACLParams struct {
	Bucket common.BucketName `json:"bucket" jsonschema:"description=Specifies the bucket whose ACL is being requested" mgc:"positional"`
}

var getGet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get the ACL for the specified bucket",
		},
		getACL,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "json"
	})
	return exec
})

func getACL(ctx context.Context, params GetBucketACLParams, cfg common.Config) (result common.AccessControlPolicy, err error) {
	req, err := newGetACLRequest(ctx, cfg, params.Bucket)
	if err != nil {
		return
	}

	res, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return common.UnwrapResponse[common.AccessControlPolicy](res, req)
}

func newGetACLRequest(ctx context.Context, cfg common.Config, bucketName common.BucketName) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("acl", "")
	url.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
