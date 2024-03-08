package buckets

import (
	"context"
	"net/http"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type BucketResponse struct {
	CreationDate string `xml:"CreationDate"`
	Name         string `xml:"Name"`
}

type ListResponse struct {
	Buckets []*BucketResponse `xml:"Buckets>Bucket"`
	Owner   *common.Owner     `xml:"Owner"`
}

func newListRequest(ctx context.Context, cfg common.Config) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, string(common.BuildHost(cfg)), nil)
}

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "list",
			Description: "List all existing Buckets",
			Scopes:      core.Scopes{"object-storage.read"},
		},
		list,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "table"
	})
	return exec
})

func list(ctx context.Context, _ struct{}, cfg common.Config) (result ListResponse, err error) {
	req, err := newListRequest(ctx, cfg)
	if err != nil {
		return
	}

	resp, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return common.UnwrapResponse[ListResponse](resp, req)
}
