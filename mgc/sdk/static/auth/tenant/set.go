package tenant

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
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
			Observations: "Changing the tenant unsets the current API Key and the selected projects, since both belong to the tenant.",
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

	// Antes da troca, não depois: se o token exchange falhar, o pior que
	// acontece é o usuário perder a seleção de projeto. Limpar depois arriscaria
	// o inverso — tenant novo com o projeto do antigo ainda valendo.
	hadProject, err := unsetProjectsForTenantChange(ctx)
	if err != nil {
		return nil, err
	}
	if hadProject {
		fmt.Print("📁 This operation unset the selected project. \n\n")
	}

	return auth.SelectTenant(ctx, params.UUID, allScopes.AsScopesString())
}

// unsetProjectsForTenantChange limpa os dois escopos de projeto e diz se havia
// algum. O projeto pertence ao tenant: mantê-lo após a troca faria as
// requisições apontarem para um escopo de outra conta.
func unsetProjectsForTenantChange(ctx context.Context) (had bool, err error) {
	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return false, fmt.Errorf("programming error: unable to retrieve config from context")
	}

	had = config.Project() != ""
	return had, config.UnsetProject()
}
