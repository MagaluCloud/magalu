package cmd

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
)

// Nomes dos parâmetros que o IAM usa para dizer onde a operação age. Só o IAM
// os tem; o gatilho é a presença deles no schema, não uma lista de rotas.
const (
	parentTypeParam = "parent_type"
	parentIDParam   = "parent_id"

	parentOrganization = "organization"
	parentProject      = "project"
)

// applyIamParentScope faz o `parent_type` mandar no `parent_id`, que passa a ser
// derivado em vez de digitado:
//
//	organization -> remove o id. É regra de negócio da API: papel de organização
//	                não tem pai, e mandar um id junto seria contraditório.
//	project      -> preenche o id quando ausente, com o escopo já resolvido para
//	                esta request (--project-id, config) ou, na falta dele, o id
//	                do tenant — que é como o IAM codifica o projeto default.
//
// Roda antes da validação de schema, então o payload validado é o final.
func applyIamParentScope(ctx context.Context, sdk *mgcSdk.Sdk, exec core.Executor, parameters core.Parameters) {
	schema := exec.ParametersSchema()
	if schema == nil || schema.Properties[parentTypeParam] == nil {
		return
	}

	switch parameters[parentTypeParam] {
	case parentOrganization:
		delete(parameters, parentIDParam)
	case parentProject:
		if id, _ := parameters[parentIDParam].(string); id == "" {
			// Gravar string vazia seria pior que omitir: viraria `parent_id=`
			// na query, que a API recebe como valor, não como ausência.
			if resolved := parentProjectID(ctx, sdk); resolved != "" {
				parameters[parentIDParam] = resolved
			} else {
				delete(parameters, parentIDParam)
			}
		}
	}
}

// parentProjectID devolve o projeto a usar. O contexto já traz o escopo
// resolvido pela tabela de precedência (flag > env > arquivo); quando não há
// nenhum, o default do IAM é o próprio tenant.
//
// Buscar o tenant aqui é seguro: as operações do IAM que têm `parent_type` são
// todas `OAuth2`, nunca api-key — então sempre há access token com as claims.
func parentProjectID(ctx context.Context, sdk *mgcSdk.Sdk) string {
	if id := config.ProjectScopeFromContext(ctx); id != "" {
		return id
	}
	id, _ := sdk.Auth().CurrentTenantID()
	return id
}
