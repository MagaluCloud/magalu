package project

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type currentResult struct {
	ID string `json:"id"`
}

var getCurrent = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:         "current",
			Summary:      "Show the project the CLI is using",
			Description:  "Reads the selected project from the local configuration. Does not reach the API, so it answers even without network or login",
			Observations: "An empty id means no project is selected: requests go to the tenant's default project.",
		},
		func(ctx context.Context) (*currentResult, error) {
			return current(ctx, struct{}{}, struct{}{})
		},
	)
})

func current(ctx context.Context, _ struct{}, _ struct{}) (*currentResult, error) {
	return currentScope(ctx, mgcConfigPkg.ProjectKey)
}

func currentScope(ctx context.Context, key string) (*currentResult, error) {
	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve config from context")
	}

	var id string
	if err := config.Get(key, &id); err != nil {
		return &currentResult{}, nil // ausência é estado válido, não erro
	}
	return &currentResult{ID: id}, nil
}
