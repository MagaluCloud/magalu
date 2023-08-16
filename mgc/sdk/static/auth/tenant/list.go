package tenant

import (
	"context"
	"fmt"

	"magalu.cloud/core"
)

func newTenantList() core.Executor {
	return core.NewStaticExecuteSimple(
		"list",
		"",
		"List all available tenants for current login",
		ListTenants,
	)
}

func ListTenants(ctx context.Context) ([]*core.Tenant, error) {
	auth := core.AuthFromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("unable to get auth from context")
	}
	return auth.ListTenants()
}
