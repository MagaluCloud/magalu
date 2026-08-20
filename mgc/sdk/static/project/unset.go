package project

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type unsetResult struct {
	Unset bool `json:"unset"`
}

// unsetConfirmMessage explica a consequência, não a ação: o usuário já sabe que
// pediu unset; o que ele precisa saber é para onde as requisições passam a ir.
const unsetConfirmMessage = "No project will be selected, and every request will be sent to the tenant's default project. Proceed?"

var getUnset = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "unset",
			Summary:     "Stop scoping the CLI to a project",
			Description: "Clears the selected project. Requests go back to the tenant's default project. The IAM scope is not affected: clear it with 'mgc iam project unset'",
		},
		func(ctx context.Context) (*unsetResult, error) {
			return unset(ctx, struct{}{}, struct{}{})
		},
	)

	return core.NewConfirmableExecutor(executor, func(parameters core.Parameters, configs core.Configs) string {
		return unsetConfirmMessage
	})
})

func unset(ctx context.Context, _ struct{}, _ struct{}) (*unsetResult, error) {
	return unsetScope(ctx, func(c *mgcConfigPkg.Config) error { return c.UnsetProject() })
}

func unsetScope(ctx context.Context, clear func(*mgcConfigPkg.Config) error) (*unsetResult, error) {
	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve config from context")
	}
	if err := clear(config); err != nil {
		return nil, err
	}
	return &unsetResult{Unset: true}, nil
}
