package project

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type currentResult struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	// Warning explica por que o nome está ausente, quando está. Vazio quando o
	// projeto foi resolvido.
	Warning string `json:"warning,omitempty"`
}

const (
	warnNotInTenant  = "project not found in this tenant: it may have been deleted, or belong to another tenant"
	warnUnresolvable = "could not reach the API to resolve the project name"
)

var getCurrent = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:         "current",
			Summary:      "Show the project the CLI is using",
			Description:  "Resolves the selected project against the API to show its name. The id comes from the local configuration, so it is still reported when the API cannot be reached",
			Observations: "An empty id means no project is selected: requests go to the tenant's default project.",
		},
		current,
	)
})

func current(ctx context.Context, _ struct{}, cfg projectConfig) (*currentResult, error) {
	return currentScope(ctx, mgcConfigPkg.ProjectKey, cfg)
}

func currentScope(ctx context.Context, key string, cfg projectConfig) (*currentResult, error) {
	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve config from context")
	}

	var id string
	if err := config.Get(key, &id); err != nil || id == "" {
		return &currentResult{}, nil // sem projeto: nada a resolver, e nenhuma request
	}

	// Não há GET /{id} na API de projects: resolver o nome é listar e procurar.
	available, err := list(ctx, listParams{}, cfg)
	return resolveCurrent(id, available, err), nil
}

// resolveCurrent decide o que reportar em cada desfecho. Nunca falha: o id vem
// da configuração local e está sempre disponível; o nome é enriquecimento, e
// perdê-lo não pode derrubar o comando que existe para dizer onde você está.
func resolveCurrent(id string, available []projectResult, lookupErr error) *currentResult {
	if lookupErr != nil {
		return &currentResult{ID: id, Warning: warnUnresolvable}
	}
	for _, p := range available {
		if p.ID == id {
			return &currentResult{ID: id, Name: p.Name, Type: p.Type}
		}
	}
	// Configuração apontando para projeto inexistente é informação útil: hoje
	// isso só apareceria como 403 em outro comando, longe da causa.
	return &currentResult{ID: id, Warning: warnNotInTenant}
}
