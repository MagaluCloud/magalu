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

var getUnset = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:         "unset",
			Summary:      "Stop scoping requests to a project",
			Description:  "Removes the selected project. Requests go back to the tenant's default project",
			Observations: "The MGC_PROJECT environment variable, when set, still selects a project.",
		},
		unset,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Project selection removed. Requests now use the tenant's default project\n"
	})
})

func unset(ctx context.Context) (*unsetResult, error) {
	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve config from context")
	}
	if err := config.UnsetProject(); err != nil {
		return nil, err
	}
	return &unsetResult{Unset: true}, nil
}
