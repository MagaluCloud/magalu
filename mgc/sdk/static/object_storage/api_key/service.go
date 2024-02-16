package api_key

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"magalu.cloud/core/auth"
	mgcHttpPkg "magalu.cloud/core/http"
)

type ApiKeysResult struct {
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	KeyPairID     string  `json:"key_pair_id"`
	KeyPairSecret string  `json:"key_pair_secret"`
	StartValidity string  `json:"start_validity"`
	EndValidity   *string `json:"end_validity,omitempty"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	TenantName    *string `json:"tenant_name,omitempty"`
}
type apiKeys struct {
	ApiKeysResult
	Tenant struct {
		UUID      string `json:"uuid"`
		LegalName string `json:"legal_name"`
	} `json:"tenant"`
	Scopes []struct {
		UUID        string `json:"uuid"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		ConsentText string `json:"consent_text"`
		Icon        string `json:"icon"`
		APIProduct  struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
			Icon string `json:"icon"`
		} `json:"api_product"`
	} `json:"scopes"`
}

type createApiKey struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	TenantID      string   `json:"tenant_id"`
	ScopeIds      []string `json:"scope_ids"`
	StartValidity string   `json:"start_validity"`
	EndValidity   string   `json:"end_validity"`
}
type ApiKeyResult struct {
	UUID string `json:"uuid,omitempty"`
}

func ListApiKeys(ctx context.Context) ([]*ApiKeysResult, error) {
	token, err := getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get current access token. Did you forget to log in?")
	}

	httpClient := mgcHttpPkg.ClientFromContext(ctx)

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, auth.FromContext(ctx).GetConfig().ApiKeys, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result []*apiKeys
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var finallyResult []*ApiKeysResult

	for _, y := range result {

		if y.RevokedAt != nil {
			continue
		}

		if y.EndValidity != nil {
			expDate, _ := time.Parse(time.RFC3339, *y.EndValidity)
			if expDate.After(time.Now()) {
				continue
			}
		}

		for _, s := range y.Scopes {
			if s.Name != "*" && auth.FromContext(ctx).GetConfig().ObjectStoreProductID != s.APIProduct.UUID {
				continue
			}
			tenantName := y.Tenant.LegalName
			y.ApiKeysResult.TenantName = &tenantName
			finallyResult = append(finallyResult, &y.ApiKeysResult)
			break
		}
	}
	return finallyResult, nil

}

func getAccessToken(ctx context.Context) (string, error) {
	auth := auth.FromContext(ctx)
	if auth == nil {
		return "", fmt.Errorf("unable to retrieve authentication configuration")
	}

	err := auth.ValidateAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("could not validate the Access Token: %w", err)
	}

	token, err := auth.AccessToken(ctx)
	if err != nil {
		return "", err
	}

	return token, nil
}
func NewApiKey(ctx context.Context, name string, description *string, expirationDate *string) (*ApiKeyResult, error) {

	token, err := getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get current access token. Did you forget to log in?")
	}

	httpClient := mgcHttpPkg.ClientFromContext(ctx)
	config := auth.FromContext(ctx).GetConfig()

	currentTenantID, err := auth.FromContext(ctx).CurrentTenantID()
	if err != nil {
		return nil, err
	}

	if description == nil {
		description = new(string)
		*description = "created from CLI"
	}

	if expirationDate == nil {
		expirationDate = new(string)
		*expirationDate = ""
	} else {
		if _, err = time.Parse(time.DateOnly, *expirationDate); err != nil {
			*expirationDate = ""
		}
	}

	newApi := &createApiKey{
		Name:          name,
		Description:   *description,
		TenantID:      currentTenantID,
		ScopeIds:      config.ObjectStoreScopeIDs,
		StartValidity: time.Now().Format(time.DateOnly),
		EndValidity:   *expirationDate,
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(newApi)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, config.ApiKeys, &buf)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", "Bearer "+token)
	r.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result ApiKeyResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}

func RevokeApiKey(ctx context.Context, uuid string) error {
	token, err := getAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("unable to get current access token. Did you forget to log in?")
	}
	httpClient := mgcHttpPkg.ClientFromContext(ctx)
	if httpClient == nil {
		err := fmt.Errorf("couldn't get http client from context")
		return err
	}

	url := fmt.Sprintf("%s/%s/revoke", auth.FromContext(ctx).GetConfig().ApiKeys, uuid)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	r.Header.Set("Authorization", "Bearer "+token)
	resp, err := httpClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	return nil
}

func SelectApiKey(ctx context.Context, uuid string) (*ApiKeysResult, error) {
	apiList, err := ListApiKeys(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range apiList {
		if v.UUID == uuid {
			if err = auth.FromContext(ctx).SetAccessKey(v.KeyPairID, v.KeyPairSecret); err != nil {
				return nil, err
			}
			return v, nil
		}
	}

	return nil, fmt.Errorf("the  API key (%s) is no longer valid", uuid)
}
