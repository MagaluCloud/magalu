package tenant

import (
	"context"
	"fmt"

	"magalu.cloud/core"
)

type tenantSetParams struct {
	ID string `jsonschema_description:"The UUID of the desired Tenant. To list all possible IDs, run auth tenant list"`
}

func newTenantSelect() core.Executor {
	return core.NewStaticExecute(
		"select",
		"",
		"Set the active Tenant to be used for all subsequential requests",
		selectTenant,
	)
}

func selectTenant(ctx context.Context, params tenantSetParams, _ struct{}) (string, error) {
	auth := core.AuthFromContext(ctx)
	if auth == nil {
		return "", fmt.Errorf("Unable to get auth from context")
	}
	err := auth.SelectTenant(params.ID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Success! Current tenant changed to %s\n", params.ID), nil
}
