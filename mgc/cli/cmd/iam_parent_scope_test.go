package cmd

import (
	"context"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
)

// execWithParentType simula um executor do IAM: o gatilho da regra é o schema
// declarar `parent_type`, não o nome do comando.
func execWithParentType(t *testing.T, withParent bool) core.Executor {
	t.Helper()
	props := map[string]*mgcSchemaPkg.Schema{"name": mgcSchemaPkg.NewStringSchema()}
	if withParent {
		props[parentTypeParam] = mgcSchemaPkg.NewStringSchema()
		props[parentIDParam] = mgcSchemaPkg.NewStringSchema()
	}
	return core.NewSimpleExecutor(core.ExecutorSpec{
		DescriptorSpec:   core.DescriptorSpec{Name: "x", Description: "x"},
		ParametersSchema: mgcSchemaPkg.NewObjectSchema(props, nil),
		ConfigsSchema:    mgcSchemaPkg.NewObjectSchema(nil, nil),
		ResultSchema:     mgcSchemaPkg.NewAnySchema(),
		Execute: func(core.Executor, context.Context, core.Parameters, core.Configs) (core.Result, error) {
			return nil, nil
		},
	})
}

// organization não tem pai: mandar um id junto seria contraditório, então ele
// sai — mesmo que a pessoa o tenha digitado.
func TestParentOrganizationDropsID(t *testing.T) {
	params := core.Parameters{parentTypeParam: parentOrganization, parentIDParam: "id-digitado"}

	applyIamParentScope(context.Background(), nil, execWithParentType(t, true), params)

	if _, ok := params[parentIDParam]; ok {
		t.Errorf("parent_id deveria ter sido removido, veio %v", params[parentIDParam])
	}
	if params[parentTypeParam] != parentOrganization {
		t.Errorf("parent_type não podia mudar: %v", params[parentTypeParam])
	}
}

// Id explícito vence: a regra só preenche o que está faltando.
func TestParentProjectKeepsExplicitID(t *testing.T) {
	params := core.Parameters{parentTypeParam: parentProject, parentIDParam: "id-explicito"}

	applyIamParentScope(context.Background(), nil, execWithParentType(t, true), params)

	if params[parentIDParam] != "id-explicito" {
		t.Errorf("parent_id = %v, quer id-explicito", params[parentIDParam])
	}
}

// Sem id digitado, entra o escopo já resolvido para esta request — que é onde
// --project-id e a configuração chegam.
func TestParentProjectFillsFromResolvedScope(t *testing.T) {
	ctx := config.NewProjectScopeContext(context.Background(), "id-do-escopo")
	params := core.Parameters{parentTypeParam: parentProject}

	applyIamParentScope(ctx, nil, execWithParentType(t, true), params)

	if params[parentIDParam] != "id-do-escopo" {
		t.Errorf("parent_id = %v, quer id-do-escopo", params[parentIDParam])
	}
}

// Comando sem `parent_type` no schema não é tocado: o gatilho é o dado, não o
// nome do produto.
func TestParentScopeIgnoresUnrelatedCommands(t *testing.T) {
	ctx := config.NewProjectScopeContext(context.Background(), "id-do-escopo")
	params := core.Parameters{"name": "vm-1"}

	applyIamParentScope(ctx, nil, execWithParentType(t, false), params)

	if len(params) != 1 {
		t.Errorf("nada podia ter sido acrescentado: %v", params)
	}
}

// parent_type ausente ou com valor fora do enum não inventa nada.
func TestParentScopeLeavesOtherValuesAlone(t *testing.T) {
	for _, v := range []any{nil, "", "outro"} {
		params := core.Parameters{parentIDParam: "id"}
		if v != nil {
			params[parentTypeParam] = v
		}
		applyIamParentScope(context.Background(), nil, execWithParentType(t, true), params)
		if params[parentIDParam] != "id" {
			t.Errorf("parent_type=%v não podia mexer no id, veio %v", v, params[parentIDParam])
		}
	}
}
