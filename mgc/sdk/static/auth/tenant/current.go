package tenant

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	"magalu.cloud/core/auth"
	mgcAuthPkg "magalu.cloud/core/auth"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/utils"
)

var getCurrent = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "current",
			Summary:     "Get the currently active Tenant",
			Description: `The current Tenant is used for all Magalu HTTP requests`,
		},
		func(ctx context.Context) (*auth.Tenant, error) {
			auth := mgcAuthPkg.FromContext(ctx)
			if auth == nil {
				return nil, fmt.Errorf("programming error: unable to get auth from context")
			}

			httpClient := mgcHttpPkg.ClientFromContext(ctx)
			if auth == nil {
				return nil, fmt.Errorf("programming error: unable to get http client from context")
			}

			return auth.CurrentTenant(ctx, &httpClient.Client)
		},
	)
})
