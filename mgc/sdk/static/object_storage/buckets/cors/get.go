package cors

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type GetBucketCorsParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Specifies the bucket whose CORS document is being requested" mgc:"positional"`
}

var getGet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get the CORS document for the specified bucket",
		},
		getCors,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "json"
	})
	return exec
})

func getCors(ctx context.Context, params GetBucketCorsParams, cfg common.Config) (result map[string]any, err error) {
	req, err := newGetCorsRequest(ctx, cfg, params.Bucket)
	if err != nil {
		return
	}

	res, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return common.UnwrapResponse[map[string]any](res, req)
}

func newGetCorsRequest(ctx context.Context, cfg common.Config, bucketName common.BucketName) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("cors", "")
	url.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
