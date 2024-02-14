package tenant

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"
)

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "list",
			Description: "List all available tenants for current login",
		},
		listTenants,
	)
})

func listTenants(ctx context.Context) ([]*mgcAuthPkg.Tenant, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("unable to get auth from context")
	}
	return auth.ListTenants(ctx)
}
