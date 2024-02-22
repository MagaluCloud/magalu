package api_key

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/utils"
)

type createParams struct {
	ApiKeyName        string  `json:"name" jsonschema:"description=Name of new api key" mgc:"positional"`
	ApiKeyDescription *string `json:"description,omitempty" jsonschema:"description=Description of new api key" mgc:"positional"`
	ApiKeyExpiration  *string `json:"expiration,omitempty" jsonschema:"description=Date to expire new api,example=2024-11-07 (YYYY-MM-DD)" mgc:"positional"`
}

var getCreate = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Scopes:      core.Scopes{"pa:api-keys:create"},
			Name:        "create",
			Description: "Create new credentials used for Object Storage requests",
		},
		create,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Key created successfully\nUuid={{.uuid}}\n"
	})
})

func create(ctx context.Context, parameter createParams, _ struct{}) (*apiKeyResult, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("unable to retrieve authentication configuration")
	}

	httpClient := mgcHttpPkg.ClientFromContext(ctx)
	config := auth.GetConfig()

	currentTenantID, err := auth.CurrentTenantID()
	if err != nil {
		return nil, err
	}

	if parameter.ApiKeyDescription == nil {
		parameter.ApiKeyDescription = new(string)
		*parameter.ApiKeyDescription = "created from CLI"
	}

	if parameter.ApiKeyExpiration == nil {
		parameter.ApiKeyExpiration = new(string)
		*parameter.ApiKeyExpiration = ""
	} else {
		if _, err = time.Parse(time.DateOnly, *parameter.ApiKeyExpiration); err != nil {
			*parameter.ApiKeyExpiration = ""
		}
	}

	newApi := &createApiKey{
		Name:          parameter.ApiKeyName,
		Description:   *parameter.ApiKeyDescription,
		TenantID:      currentTenantID,
		ScopeIds:      config.ObjectStoreScopeIDs,
		StartValidity: time.Now().Format(time.DateOnly),
		EndValidity:   *parameter.ApiKeyExpiration,
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(newApi)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, config.ApiKeysUrl, &buf)
	if err != nil {
		return nil, err
	}

	token, err := auth.AccessToken(ctx)
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
	var result apiKeyResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
