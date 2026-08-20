package project

// 'mgc iam project ...' mora neste pacote, e não em um static/iam próprio,
// porque é o MESMO domínio: resolve o mesmo endpoint, contra a mesma lista, com
// as mesmas regras de id-ou-nome. Só a chave de config muda. Um pacote separado
// exigiria exportar metade daqui para consumir do outro lado.
//
// O grupo devolvido se chama `iam` e o MergeGroup (core/merge.go) o funde com o
// módulo `iam` que vem da spec — por isso o nome tem de bater exatamente. O
// estático vem primeiro na lista de fontes, então vence em caso de colisão; hoje
// não há colisão, porque a spec do IAM não publica `projects`.

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

// iamUnsetConfirmMessage é mais duro que o do unset da CLI de propósito: sem
// projeto de IAM a ação não cai num default, ela passa a valer para o tenant
// inteiro.
const iamUnsetConfirmMessage = "No IAM project will be selected, and IAM commands will apply to the entire tenant. Proceed?"

var GetIamGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "iam",
			Summary:     "Identity and Access Management",
			Description: "Manage identities, permissions and the project scope that IAM commands apply to",
			GroupID:     "settings",
		},
		func() []core.Descriptor {
			return []core.Descriptor{getIamProjectGroup()}
		},
	)
})

var getIamProjectGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "project",
			Summary: "Choose which project IAM commands apply to",
			Description: `The IAM scope is deliberately separate from the CLI project set by 'mgc project set':
it decides where identities, roles and permissions are created, and getting it
wrong is not a failed request, it is a change in the wrong place`,
			Observations: "With no IAM project selected, IAM commands apply to the entire tenant.",
		},
		func() []core.Descriptor {
			return []core.Descriptor{getIamSet(), getIamCurrent(), getIamUnset()}
		},
	)
})

var getIamSet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:         "set",
			Summary:      "Set the project IAM commands apply to",
			Description:  "Scopes IAM commands to this project. Does not affect the project used by the rest of the CLI, which is set by 'mgc project set'",
			Observations: "Changing the tenant clears the selected project.",
		},
		iamSet,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Success! IAM commands now apply to {{.name}} ({{.id}})\n"
	})
})

var getIamCurrent = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:         "current",
			Summary:      "Show the project IAM commands apply to",
			Description:  "Reads the IAM scope from the local configuration. Does not reach the API, so it answers even without network or login",
			Observations: "An empty id means IAM commands apply to the entire tenant.",
		},
		func(ctx context.Context) (*currentResult, error) {
			return iamCurrent(ctx, struct{}{}, struct{}{})
		},
	)
})

var getIamUnset = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "unset",
			Summary:     "Stop scoping IAM commands to a project",
			Description: "Clears the IAM scope. IAM commands go back to applying to the entire tenant. The CLI project is not affected",
		},
		func(ctx context.Context) (*unsetResult, error) {
			return iamUnset(ctx, struct{}{}, struct{}{})
		},
	)

	return core.NewConfirmableExecutor(executor, func(parameters core.Parameters, configs core.Configs) string {
		return iamUnsetConfirmMessage
	})
})

func iamSet(ctx context.Context, params setParams, cfg projectConfig) (*setResult, error) {
	return setScope(ctx, mgcConfigPkg.IamProjectKey, params.IDOrName, cfg)
}

func iamCurrent(ctx context.Context, _ struct{}, _ struct{}) (*currentResult, error) {
	return currentScope(ctx, mgcConfigPkg.IamProjectKey)
}

func iamUnset(ctx context.Context, _ struct{}, _ struct{}) (*unsetResult, error) {
	return unsetScope(ctx, func(c *mgcConfigPkg.Config) error { return c.UnsetIamProject() })
}
