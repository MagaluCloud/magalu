package tenant

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	mgcAuthScope "github.com/MagaluCloud/magalu/mgc/sdk/static/auth/scopes"
)

type tenantSetParams struct {
	UUID string `json:"uuid" jsonschema_description:"The UUID of the desired Tenant. To list all possible IDs, run auth tenant list" mgc:"positional"`
}

var getSet = utils.NewLazyLoader[core.Executor](newSet)

func newSet() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:         "set",
			Description:  "Set the active Tenant to be used for all subsequent requests",
			Observations: "Changing the tenant unsets the current API Key and selects the new tenant's default project.",
		},
		setTenant,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Success! Current tenant changed to {{.uuid}}\n"
	})
}

func setTenant(ctx context.Context, params tenantSetParams, _ struct{}) (
	*mgcAuthPkg.TokenExchangeResult, error,
) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("unable to get auth from context")
	}

	allScopes, err := mgcAuthScope.ListAllAvailable(ctx)
	if err != nil {
		return nil, err
	}

	id, key := auth.AccessKeyPair()
	if id != "" && key != "" {
		fmt.Print("🔐 This operation unset the current api key. \n\n")
		err = auth.UnsetAccessKey()
		if err != nil {
			return nil, err
		}
	}

	result, err := auth.SelectTenant(ctx, params.UUID, allScopes.AsScopesString())
	if err != nil {
		return nil, err
	}

	// Depois da troca, nunca antes: gravar o projeto de um tenant que o token
	// seguinte pode nem ter assumido deixaria a configuração apontando para o
	// lugar errado justamente quando a troca falha.
	if err = SelectDefaultProject(ctx, auth); err != nil {
		return nil, err
	}
	fmt.Print("📁 The selected project changed to this tenant's default. \n\n")

	return result, nil
}
