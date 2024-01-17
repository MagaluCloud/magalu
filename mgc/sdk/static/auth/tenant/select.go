package tenant

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"
)

type tenantSetParams struct {
	ID string `json:"id" jsonschema_description:"The UUID of the desired Tenant. To list all possible IDs, run auth tenant list" mgc:"positional"`
}

var getSelect = utils.NewLazyLoader[core.Executor](newSelect)

func newSelect() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "select",
			Description: "Set the active Tenant to be used for all subsequential requests",
		},
		selectTenant,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Success! Current tenant changed to {{.id}}\n"
	})
}

func selectTenant(ctx context.Context, params tenantSetParams, _ struct{}) (*mgcAuthPkg.TokenExchangeResult, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("Unable to get auth from context")
	}
	return auth.SelectTenant(ctx, params.ID)
}
