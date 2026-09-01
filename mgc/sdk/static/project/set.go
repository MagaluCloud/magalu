package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type setParams struct {
	IDOrName string `json:"id-or-name" jsonschema:"description=UUID or name of the project to use. Run 'mgc project list' to see the available ones,required,example=my-project" mgc:"positional"`
}

type setResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var getSet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:          "set",
			ProjectScoped: true,
			Summary:       "Set the project used by the CLI",
			Description:   "All subsequent requests are scoped to this project, IAM included. Undo it with 'mgc project default'",
			Observations:  "Changing the tenant clears the selected project.",
		},
		set,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Success! Current project changed to {{.name}} ({{.id}})\n"
	})
})

// set resolve SEMPRE contra a lista: o positional aceita id ou nome, e nenhum
// valor é reservado ali.
//
// Houve um atalho aqui, `set default` gravando o literal sem consultar nada.
// Ele tornava inalcançável pelo nome um projeto que a pessoa tivesse criado
// chamado `default` — e em silêncio, selecionando o default do tenant no lugar.
// Nome de projeto não tem restrição nenhuma no schema da API, então essa
// colisão é permitida. Quem diz "o default do tenant" é 'mgc project default',
// que não compartilha sintaxe com nome de projeto e por isso nunca empata.
func set(ctx context.Context, params setParams, cfg projectConfig) (*setResult, error) {
	return setScope(ctx, mgcConfigPkg.ProjectKey, params.IDOrName, cfg)
}

// setScope recebe a chave de config por parâmetro por herança do tempo em que
// havia duas ('project' e 'iamProject'). Hoje só uma chega aqui.
func setScope(ctx context.Context, key, query string, cfg projectConfig) (*setResult, error) {
	// Valida antes de gravar: um id inválido em configuração faria TODO comando
	// seguinte apontar para um escopo inexistente, e o erro apareceria longe daqui.
	available, err := list(ctx, listParams{}, cfg)
	if err != nil {
		return nil, err
	}
	return applyScope(ctx, key, query, available)
}

// applyScope resolve a consulta contra a lista e grava. Separado do setScope
// para que a resolução e a escrita possam ser testadas sem rede.
func applyScope(ctx context.Context, key, query string, available []projectResult) (*setResult, error) {
	found, err := resolveProject(available, query)
	if err != nil {
		return nil, err
	}

	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve config from context")
	}
	if err := config.Set(key, found.ID); err != nil {
		return nil, err
	}

	return &setResult{ID: found.ID, Name: found.Name}, nil
}

// resolveProject aceita id ou nome. O id vence: é o identificador não-ambíguo,
// então uma consulta que casa um id nunca é desviada por um nome igual em outro
// projeto. O nome ignora caixa, mas nome repetido não elege um vencedor
// arbitrário — devolve erro pedindo o id.
func resolveProject(available []projectResult, query string) (*projectResult, error) {
	if len(available) == 0 {
		return nil, fmt.Errorf("no projects available for this tenant. Create one with 'mgc project create'")
	}

	for i := range available {
		if available[i].ID == query {
			return &available[i], nil
		}
	}

	var matches []projectResult
	for _, p := range available {
		if strings.EqualFold(p.Name, query) {
			matches = append(matches, p)
		}
	}

	switch len(matches) {
	case 1:
		return &matches[0], nil
	case 0:
		// Quem digita `set default` quase sempre quer o default do tenant, não um
		// projeto com esse nome. Sem esta linha o erro seria tecnicamente correto
		// e inútil: lista projetos e não diz onde está o que a pessoa procura.
		hint := ""
		if strings.EqualFold(query, mgcConfigPkg.ProjectDefault) {
			hint = "Did you mean 'mgc project default' (the tenant's default project)?\n"
		}
		return nil, fmt.Errorf("project %q not found.\n%sAvailable projects:%s", query, hint, projectLines(available))
	default:
		return nil, fmt.Errorf("more than one project named %q. Use the id instead:%s", query, projectLines(matches))
	}
}

func projectLines(projects []projectResult) string {
	var b strings.Builder
	for _, p := range projects {
		fmt.Fprintf(&b, "\n  %s  %s", p.ID, p.Name)
	}
	return b.String()
}
