package api_key

import (
	"context"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Scopes:      core.Scopes{"pa:api-keys:read"},
			Name:        "list",
			Description: "List valid Object Storage credentials",
		},
		list,
	)
})

func list(ctx context.Context) ([]*apiKeysResult, error) {
	return ListApiKeys(ctx)
}
