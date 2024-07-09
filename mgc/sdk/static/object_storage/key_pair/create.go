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
	KeyPairName        string  `json:"name" jsonschema:"description=Name of new key pair" mgc:"positional"`
	KeyPairDescription *string `json:"description,omitempty" jsonschema:"description=Description of new key pair" mgc:"positional"`
	KeyPairExpiration  *string `json:"expiration,omitempty" jsonschema:"description=Date to expire new key pair,example=2024-11-07 (YYYY-MM-DD)" mgc:"positional"`
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
		return "template={{if .used}}Key created and used successfully{{else}}Key created successfully{{end}} Uuid={{.uuid}}\n"
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

	if parameter.KeyPairDescription == nil {
		parameter.KeyPairDescription = new(string)
		*parameter.KeyPairDescription = "created from CLI"
	}

	if parameter.KeyPairExpiration == nil {
		parameter.KeyPairExpiration = new(string)
		*parameter.KeyPairExpiration = ""
	} else {
		if _, err = time.Parse(time.DateOnly, *parameter.KeyPairExpiration); err != nil {
			*parameter.KeyPairExpiration = ""
		}
	}

	const reason = "permission to read and write at object-storage"

	newApi := &createApiKey{
		Name:        parameter.KeyPairName,
		Description: *parameter.KeyPairDescription,
		TenantID:    currentTenantID,
		ScopesList: []scopesObjectStorage{
			{ID: config.ObjectStoreScopeIDs[0], RequestReason: reason},
			{ID: config.ObjectStoreScopeIDs[1], RequestReason: reason},
		},
		StartValidity: time.Now().Format(time.DateOnly),
		EndValidity:   *parameter.KeyPairExpiration,
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

	id, _ := auth.AccessKeyPair()
	if id == "" {
		_, err = setCurrent(ctx, selectParams{UUID: result.UUID}, struct{}{})
		if err == nil {
			result.Used = true
		}
	}

	return &result, nil
}