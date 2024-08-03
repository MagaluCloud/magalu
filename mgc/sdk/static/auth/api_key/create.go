package api_key

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	_ "embed"

	"github.com/pterm/pterm"
	"golang.org/x/exp/maps"
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

//go:embed scopes.json
var scopesFile []byte

var getCreate = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Scopes:      core.Scopes{scope_PA},
			Name:        "create",
			Summary:     "Create a new API Key",
			Description: "Select the scopes that the new API Key will have access to and set an expiration date",
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

	var scopesList map[string]string

	err := json.Unmarshal(scopesFile, &scopesList)
	if err != nil {
		return nil, err
	}

	input := maps.Keys(scopesList)
	slices.Sort(input)
	op, err := pterm.DefaultInteractiveMultiselect.
		WithDefaultText("Select scopes").
		WithMaxHeight(14).
		WithOptions(input).
		Show()
	if err != nil {
		return nil, err
	}

	if len(op) == 0 {
		return nil, fmt.Errorf("no scopes selected")
	}

	var scopesCreateList []scopesCreate

	for _, v := range op {
		scopesCreateList = append(scopesCreateList, scopesCreate{
			ID: scopesList[v],
		})
	}

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

	config := auth.GetConfig()

	newApi := &createApiKey{
		Name:          parameter.ApiKeyName,
		Description:   *parameter.ApiKeyDescription,
		TenantID:      currentTenantID,
		ScopesList:    scopesCreateList,
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
