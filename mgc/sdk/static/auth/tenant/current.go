package tenant

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
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

			return auth.CurrentTenant(ctx)
		},
	)
})
