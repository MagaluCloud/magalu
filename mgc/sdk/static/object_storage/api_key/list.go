package api_key

import (
	"context"

	"magalu.cloud/core/utils"

	"magalu.cloud/core"
)

var getKeys = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "list",
			Description: "List valid Object Storage credentials",
		},
		list,
	)
})

func list(ctx context.Context) ([]*ApiKeysResult, error) {
	return ListApiKeys(ctx)
}
