package keys

import (
	"context"
	"fmt"

	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"

	"magalu.cloud/core"
)

var getKeys = utils.NewLazyLoader[core.Executor](newList)

func newList() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "list",
			Description: "List valid Object Storage credentials",
		},
		list,
	)
}

func list(ctx context.Context) ([]*mgcAuthPkg.ApiKeysResult, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("unable to retrieve authentication configuration")
	}

	return auth.ListApiKeys(ctx)
}
