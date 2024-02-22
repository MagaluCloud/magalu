package tenant

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"
	mgcAuthScope "magalu.cloud/sdk/static/auth/scopes"
)

type tenantSetParams struct {
	UUID string `json:"uuid" jsonschema_description:"The UUID of the desired Tenant. To list all possible IDs, run auth tenant list" mgc:"positional"`
}

var getSet = utils.NewLazyLoader[core.Executor](newSet)

func newSet() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "Set the active Tenant to be used for all subsequential requests",
		},
		setTenant,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Success! Current tenant changed to {{.uuid}}\n"
	})
}

func setTenant(ctx context.Context, params tenantSetParams, _ struct{}) (*mgcAuthPkg.TokenExchangeResult, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("unable to get auth from context")
	}

	allScopes, err := mgcAuthScope.ListAllAvailable(ctx)
	if err != nil {
		return nil, err
	}

	return auth.SelectTenant(ctx, params.UUID, allScopes.AsScopesString())
}
