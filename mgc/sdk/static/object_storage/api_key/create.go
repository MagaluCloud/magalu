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
			Scopes:      core.Scopes{scope_PA},
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
		return nil, fmt.Errorf("programming error: unable to retrieve auth configuration from context")
	}

	httpClient := auth.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve HTTP Client from context")
	}

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

	const reason = "permission to read and write at object-storage"

	newApi := &createApiKey{
		Name:        parameter.ApiKeyName,
		Description: *parameter.ApiKeyDescription,
		TenantID:    currentTenantID,
		ScopesList: []scopesObjectStorage{
			{ID: config.ObjectStoreScopeIDs[0], RequestReason: reason},
			{ID: config.ObjectStoreScopeIDs[1], RequestReason: reason},
		},
		StartValidity: time.Now().Format(time.DateOnly),
		EndValidity:   *parameter.ApiKeyExpiration,
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(newApi)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, config.ApiKeysUrlV2, &buf)
	if err != nil {
		return nil, err
	}

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
