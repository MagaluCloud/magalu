package tenant

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"magalu.cloud/core"
)

type Tenant struct {
	UUID        string `json:"uuid" mapstructure:"uuid"`
	Name        string `json:"legal_name" mapstructure:"legal_name"`
	Email       string `json:"email" mapstructure:"email"`
	IsManaged   bool   `json:"is_managed" mapstructure:"is_managed"`
	IsDelegated bool   `json:"is_delegated" mapstructure:"is_delegated"`
}

func newTenantList() core.Executor {
	return core.NewStaticExecuteSimple(
		"list",
		"",
		"List all available tenants for current login",
		ListTenants,
	)
}

func ListTenants(ctx context.Context) ([]*Tenant, error) {
	auth := core.AuthFromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("unable to get auth from context")
	}

	at, err := auth.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("unable to get current access token. Did you forget to log in?")
	}

	r, err := http.NewRequest(http.MethodGet, auth.TenantsListUrl(), nil)
	r.Header.Set("Authorization", at)

	if err != nil {
		return nil, err
	}

	client := core.HttpClientFromContext(ctx)
	if client == nil {
		return nil, fmt.Errorf("Unable to get http client from context")
	}

	resp, err := client.Do(r)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()
	var result []*Tenant
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
