package tenant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// "net/http/httputil"

	"magalu.cloud/core"
)

type tenantSetParams struct {
	ID string `jsonschema_description:"The UUID of the desired Tenant. To list all possible IDs, run auth tenant list"`
}

type tenantSetPayload struct {
	AccessToken  string `json:"access_token"`
	CreatedAt    int    `json:"created_at"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"scope_type"`
}

func newTenantSelect() core.Executor {
	return core.NewStaticExecute(
		"select",
		"",
		"Set the active Tenant to be used for all subsequential requests",
		func(ctx context.Context, params tenantSetParams, _ struct{}) (string, error) {
			return SelectTenant(ctx, params.ID)
		},
	)
}

func SelectTenant(ctx context.Context, id string) (string, error) {
	auth := core.AuthFromContext(ctx)
	if auth == nil {
		return "Failed", fmt.Errorf("unable to get auth from context")
	}

	at, err := auth.AccessToken()
	if err != nil {
		return "Failed", fmt.Errorf("unable to get current access token. Did you forget to log in?")
	}

	data := map[string]any{"tenant": id}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "Failed", err
	}

	bodyReader := bytes.NewReader(jsonData)
	r, err := http.NewRequest(http.MethodPost, auth.TenantsSelectUrl(), bodyReader)
	r.Header.Set("Authorization", "Bearer "+at)
	r.Header.Set("Content-Type", "application/json")

	if err != nil {
		return "Failed", err
	}

	client := core.HttpClientFromContext(ctx)
	if client == nil {
		return "Failed", fmt.Errorf("unable to get http client from context")
	}

	resp, err := client.Do(r)
	if err != nil {
		return "Failed", err
	}

	defer r.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Failed", core.NewHttpErrorFromResponse(resp)
	}

	payload := &tenantSetPayload{}
	if err = json.NewDecoder(resp.Body).Decode(payload); err != nil {
		return "Failed", err
	}

	err = auth.SetTokens(&core.LoginResult{
		AccessToken:  payload.AccessToken,
		RefreshToken: payload.RefreshToken,
	})
	if err != nil {
		return "Failed", err
	}

	if err := auth.SetCurrentTenantID(id); err != nil {
		return "Failed", err
	}

	return "Success!", nil
}
