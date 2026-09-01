package cmd

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
)

// resolveProjectScope devolve o projeto selecionado nesta invocação, ou "" para
// não escopar nada. É o valor que vai no header 'x-project-id' dos produtos e,
// no IAM, o que preenche o `parent_id` (ver applyIamParentScope).
//
// A precedência flag > env > arquivo já vem de Config.Get: --project-id grava
// no temp config, consultado antes do viper.
//
// Houve aqui uma flag --scope, com os valores 'tenant' e 'default', de quando o
// plano era escopar o IAM por header. O IAM não lê header nenhum: ele escopa por
// `parent_type`/`parent_id`, e quem diz "a organização inteira" é
// `--parent-type organization`. A flag sobrevivia dizendo uma coisa e fazendo
// outra — `--scope tenant` apagava o projeto e mandava `parent_type=project`
// apontando para o id do tenant.
func resolveProjectScope(cfg *config.Config, scope core.ProjectScope) string {
	if scope == "" || cfg == nil {
		return "" // produto não declarou escopo: não participa
	}

	var id string
	if err := cfg.Get(config.ProjectKey, &id); err != nil {
		return ""
	}
	return id
}

// withProjectScope resolve o escopo do comando e o injeta no contexto, de onde
// o transport o lê. Comando de produto não escopável passa reto — o contexto
// fica sem valor e nenhum header é enviado.
func withProjectScope(ctx context.Context, sdk *mgcSdk.Sdk, exec core.Executor) context.Context {
	scope := exec.DescriptorSpec().ProjectScope
	if scope == "" {
		return ctx
	}
	return config.NewProjectScopeContext(ctx, resolveProjectScope(sdk.Config(), scope))
}
