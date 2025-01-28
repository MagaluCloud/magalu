package object_lock

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type GetBucketObjectLockParams struct {
	Object mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Specifies the object whose lock is being requested" mgc:"positional"`
}

type objectLockRetentionResponse struct {
	Mode            common.ObjectLockMode
	RetainUntilDate string `xml:",omitempty"`
}

var getGet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get object locking configuration for the specified object",
		},
		getObjectLocking,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "json"
	})
	return exec
})

func getObjectLocking(ctx context.Context, params GetBucketObjectLockParams, cfg common.Config) (result objectLockRetentionResponse, err error) {
	objectURI := mgcSchemaPkg.URI(params.Object)

	req, err := newGetObjectLockingRequest(ctx, cfg, objectURI)
	if err != nil {
		return
	}

	res, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return common.UnwrapResponse[objectLockRetentionResponse](res, req)
}

func newGetObjectLockingRequest(ctx context.Context, cfg common.Config, src mgcSchemaPkg.URI) (*http.Request, error) {
	url, err := common.BuildBucketHostWithPath(cfg, common.NewBucketNameFromURI(src), src.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, string(url), nil)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := req.URL.Query()
	query.Add("retention", "")
	req.URL.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, req.URL.String(), nil)
}
