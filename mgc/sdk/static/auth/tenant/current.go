package tenant

import (
	"context"
	"fmt"

	"magalu.cloud/core"
)

type tenantCurrentResult struct {
	ID string
}

func newTenantCurrent() core.Executor {
	return core.NewStaticExecuteSimple(
		"current",
		"",
		"Get the currently active Tenant",
		func(ctx context.Context) (*tenantCurrentResult, error) {
			auth := core.AuthFromContext(ctx)
			if auth == nil {
				return nil, fmt.Errorf("unable to get auth from context")
			}

			id := auth.CurrentTenantID()
			if id == "" {
				return nil, fmt.Errorf("current tenant ID is empty. Try logging in at least once or run the 'auth tenant select' operation")
			}

			return &tenantCurrentResult{ID: id}, nil
		},
	)
}
