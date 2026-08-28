package project

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type defaultResult struct {
	Default bool `json:"default"`
}

// defaultConfirmMessage explica a consequência, não a ação: quem digitou o
// comando já sabe o que pediu; o que precisa saber é para onde as requisições
// passam a ir.
const defaultConfirmMessage = "No project will be selected, and every request will be sent to the tenant's default project. Proceed?"

var getDefault = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "default",
			Summary:     "Go back to the tenant's default project",
			Description: "Clears the selected project, so requests are scoped to the tenant's default project again",
		},
		func(ctx context.Context) (*defaultResult, error) {
			return useDefault(ctx, struct{}{}, struct{}{})
		},
	)

	return core.NewConfirmableExecutor(executor, func(parameters core.Parameters, configs core.Configs) string {
		return defaultConfirmMessage
	})
})

// useDefault não pode se chamar `default`: é palavra reservada da linguagem.
func useDefault(ctx context.Context, _ struct{}, _ struct{}) (*defaultResult, error) {
	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve config from context")
	}
	if err := config.UnsetProject(); err != nil {
		return nil, err
	}
	return &defaultResult{Default: true}, nil
}
