package buckets

import (
	"context"
	"net/http"

	"magalu.cloud/core"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type createParams struct {
	Name     common.BucketName `json:"name" jsonschema:"description=Name of the bucket to be created" mgc:"positional"`
	ACL      string            `json:"acl,omitempty" jsonschema:"description=ACL Rules for the bucket"`
	Location string            `json:"location,omitempty" jsonschema:"description=Location constraint for the bucket,default=br-ne-1"`
}

var getCreate = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewReflectedSimpleExecutor[createParams, common.Config, core.Value](
		core.ExecutorSpec{
			DescriptorSpec: core.DescriptorSpec{
				Name:        "create",
				Summary:     "Create a new Bucket",
				Description: `Buckets are "containers" that are able to store various Objects inside`,
			},
			Links: utils.NewLazyLoaderWithArg(func(e core.Executor) core.Links {
				return core.Links{
					"delete": core.NewSimpleLink(e, getDelete()),
					"list":   core.NewSimpleLink(e, getList()),
				}
			}),
		},
		create,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Created bucket {{.name}}\n"
	})
})

func newCreateRequest(ctx context.Context, cfg common.Config, bucket common.BucketName) (*http.Request, error) {
	url, err := common.BuildBucketHost(cfg, bucket)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}
	return http.NewRequestWithContext(ctx, http.MethodPut, string(url), nil)
}

func create(ctx context.Context, params createParams, cfg common.Config) (core.Value, error) {
	req, err := newCreateRequest(ctx, cfg, params.Name)
	if err != nil {
		return nil, err
	}

	res, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return mgcHttpPkg.UnwrapResponse[any](res)
}
