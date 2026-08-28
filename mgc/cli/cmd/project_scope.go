package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
	"github.com/spf13/cobra"
)

// Valores aceitos por --scope. Só existem no IAM: os demais produtos não têm
// nível de tenant, e para eles "default" já é a ausência de configuração.
const (
	scopeDefault = "default"
	scopeTenant  = "tenant"
)

// resolveProjectScope devolve o id que deve ir no 'x-project-id', ou "" para
// omitir o header.
//
// A tabela que ele implementa:
//
//	produto não escopável       -> ""            (nunca carimba)
//	--scope tenant   (só IAM)   -> ""            (omitir = o tenant inteiro)
//	--scope default  (só IAM)   -> id do tenant  (codificação do IAM p/ o default)
//	sem flag                    -> a chave de config do produto
//
// A precedência flag > env > arquivo já vem de Config.Get: a flag global grava
// no temp config, consultado antes do viper.
func resolveProjectScope(cfg *config.Config, scope core.ProjectScope, flagScope string) string {
	return resolveProjectScopeWithTenant(cfg, scope, flagScope, "")
}

// resolveProjectScopeWithTenant é a forma testável: recebe o tenant já
// resolvido, porque obtê-lo depende do token e não cabe numa função pura.
func resolveProjectScopeWithTenant(cfg *config.Config, scope core.ProjectScope, flagScope, tenantID string) string {
	if scope == "" || cfg == nil {
		return "" // produto não declarou escopo: não participa
	}

	switch flagScope {
	case scopeTenant:
		// Omitir É como se diz "tenant inteiro" para o IAM.
		return ""
	case scopeDefault:
		// O IAM codifica o projeto default como o id do tenant. Quirk da API,
		// escondido do usuário — ele nunca digita esse id.
		return tenantID
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
func withProjectScope(ctx context.Context, sdk *mgcSdk.Sdk, cmd *cobra.Command, exec core.Executor) context.Context {
	spec := exec.DescriptorSpec()
	scope := spec.ProjectScope
	if scope == "" {
		return ctx
	}

	tenantID := ""
	// Só o `--scope default` precisa do tenant, e obtê-lo pode falhar (token
	// ausente ou expirado). Falhar aqui não deve derrubar o comando: sem
	// tenant, o escopo fica vazio e a request segue sem header — o mesmo que
	// não ter configurado nada.
	if getScopeFlag(cmd) == scopeDefault {
		if id, err := sdk.Auth().CurrentTenantID(); err == nil {
			tenantID = id
		}
	}

	// MODO AVISO. A recusa está desenhada e testada, mas quebraria todo script
	// que hoje roda `mgc iam ...` sem escopo. Avisar por uma release dá janela de
	// migração e, de quebra, mede quais operações são chamadas sem escopo — que é
	// como descobrir as isenções com dado em vez de palpite.
	//
	// Para virar recusa: trocar este Print por um erro devolvido ao chamador.
	if msg := checkScopeRequired(sdk.Config(), spec.ScopeRequired, scope, getScopeFlag(cmd)); msg != "" {
		fmt.Fprintf(os.Stderr, "\n⚠  %s\n\n", msg)
	}

	return config.NewProjectScopeContext(
		ctx,
		resolveProjectScopeWithTenant(sdk.Config(), scope, getScopeFlag(cmd), tenantID),
	)
}

// scopeRequiredMessage ensina as três saídas. Erro que diz "faltou escopo" sem
// dizer como fornecê-lo deixa a pessoa presa.
const scopeRequiredMessage = `this IAM write applies to the ENTIRE tenant unless you scope it. Choose one:
  --scope default        the tenant's default project
  --scope tenant         the entire tenant, explicitly
  --project-id <id>      a specific project
Or select a project once with 'mgc project set <id-or-name>'.`

// checkScopeRequired devolve a mensagem a exibir, ou "" quando está tudo certo.
//
// Devolver mensagem em vez de erro é o que permite os dois modos: hoje ela é
// impressa como AVISO e o comando segue; virar recusa é trocar o Print por um
// return de erro em handleExecutorPre — uma linha, quando você decidir.
func checkScopeRequired(cfg *config.Config, required bool, scope core.ProjectScope, flagScope string) string {
	if !required || scope != core.ProjectScopeIAM {
		return "" // leitura, ou produto onde omitir significa o projeto default
	}
	if flagScope == scopeTenant || flagScope == scopeDefault {
		return "" // o usuário declarou o escopo nesta invocação
	}
	if cfg != nil {
		var id string
		if err := cfg.Get(config.ProjectKey, &id); err == nil && id != "" {
			return "" // há escopo configurado
		}
	}
	return scopeRequiredMessage
}
