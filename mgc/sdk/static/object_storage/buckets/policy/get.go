package policy

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type GetBucketPolicyParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Specifies the bucket whose policy document is being requested" mgc:"positional"`
}

var getGet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get the policy document for the specified bucket",
		},
		getPolicy,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "json"
	})
	return exec
})

func getPolicy(ctx context.Context, params GetBucketPolicyParams, cfg common.Config) (result map[string]any, err error) {
	req, err := newGetPolicyRequest(ctx, cfg, params.Bucket)
	if err != nil {
		return
	}

	res, err := common.SendRequest(ctx, req, cfg)
	if err != nil {
		return
	}

	return common.UnwrapResponse[map[string]any](res, req)
}

func newGetPolicyRequest(ctx context.Context, cfg common.Config, bucketName common.BucketName) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("policy", "")
	url.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
