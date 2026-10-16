package project

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type currentResult struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Warning string `json:"warning,omitempty"`
}

const (
	warnNotInTenant  = "project not found in this tenant"
	warnUnresolvable = "could not reach the API to resolve the project name"
	warnNoProject    = "no project selected: requests use the tenant's default project. Run 'mgc alpha project set' to pick one"
)

var getCurrent = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:         "current",
			Summary:      "Show the project the CLI is using",
			Description:  "Resolves the selected project against the API to show its name. The id comes from the local configuration, so it is still reported when the API cannot be reached",
			Observations: "A tenant's default project has the same id as the tenant, and is listed by 'mgc alpha project list' like any other.",
		},
		current,
	)
})

func current(ctx context.Context, _ struct{}, cfg projectConfig) (*currentResult, error) {
	return currentScope(ctx, mgcConfigPkg.ProjectKey, cfg)
}

func currentScope(ctx context.Context, key string, cfg projectConfig) (*currentResult, error) {
	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve config from context")
	}

	var id string
	if err := config.Get(key, &id); err != nil || id == "" {
		return &currentResult{Warning: warnNoProject}, nil
	}

	available, err := list(ctx, listParams{}, cfg)
	return resolveCurrent(id, available, err), nil
}

func resolveCurrent(id string, available []projectResult, lookupErr error) *currentResult {
	if lookupErr != nil {
		return &currentResult{ID: id, Warning: warnUnresolvable}
	}
	for _, p := range available {
		if p.ID == id {
			return &currentResult{ID: id, Name: p.Name, Type: p.Type}
		}
	}
	return &currentResult{ID: id, Warning: warnNotInTenant}
}
