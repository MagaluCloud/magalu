package api_key

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	_ "embed"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/pterm/pterm"
	"golang.org/x/exp/maps"
)

const scopesURL = "https://api.magalu.cloud/iam/api/v1/scopes"

type createParams struct {
	ApiKeyName        string   `json:"name" jsonschema:"description=Name of new api key,required,example=My MGC Key" mgc:"positional"`
	ApiKeyDescription string   `json:"description,omitempty" jsonschema:"description=Description of new api key,example=created from MGC CLI,default=Created from CLI"`
	ApiKeyExpiration  string   `json:"expiration,omitempty" jsonschema:"description=Date to expire new api (YYYY-MM-DD),example=2024-11-07"`
	Scopes            []string `json:"scopes,omitempty" jsonschema:"description=List of scopes to assign to new api key,example=dbaas.read"`
}

type Scope struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	UUID  string `json:"uuid"`
}

type APIProduct struct {
	Name   string  `json:"name"`
	Scopes []Scope `json:"scopes"`
	UUID   string  `json:"uuid"`
}

type Platform struct {
	APIProducts []APIProduct `json:"api_products"`
	Name        string       `json:"name"`
	UUID        string       `json:"uuid"`
}

type PlatformsResponse []Platform

var getCreate = utils.NewLazyLoader(func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Scopes:      core.Scopes{scope_PA},
			Name:        "create",
			Summary:     "Create a new API Key for your tenant",
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

	scopesRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, scopesURL, nil)
	if err != nil {
		return nil, err
	}

	scopesResponse, err := httpClient.Do(scopesRequest)
	if err != nil {
		return nil, err
	}

	var scopesListFile PlatformsResponse

	defer scopesResponse.Body.Close()
	if err = json.NewDecoder(scopesResponse.Body).Decode(&scopesListFile); err != nil {
		return nil, err
	}

	scopesTitleMap := make(map[string]string)
	scopeNameMap := make(map[string]string)

	for _, company := range scopesListFile {
		if company.Name == "Magalu Cloud" {
			for _, product := range company.APIProducts {
				for _, scope := range product.Scopes {
					scopeName := product.Name + " [" + scope.Name + "]" + " - " + scope.Title
					scopesTitleMap[scopeName] = scope.UUID
					scopeNameMap[strings.ToLower(scope.Name)] = scope.UUID
				}
			}
		}
	}

	var scopesCreateList []scopesCreate
	var invalidScopes []string
	if len(parameter.Scopes) > 0 {
		for _, v := range parameter.Scopes {
			if id, ok := scopeNameMap[strings.ToLower(v)]; ok {
				scopesCreateList = append(scopesCreateList, scopesCreate{
					ID: id,
				})
			} else {
				invalidScopes = append(invalidScopes, v)
			}
		}
		if len(invalidScopes) > 0 {
			return nil, fmt.Errorf("invalid scopes: %s", strings.Join(invalidScopes, ", "))
		}
	} else {
		input := maps.Keys(scopesTitleMap)
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

		for _, v := range op {
			scopesCreateList = append(scopesCreateList, scopesCreate{
				ID: scopesTitleMap[v],
			})
		}
	}

	currentTenantID, err := auth.CurrentTenantID()
	if err != nil {
		return nil, err
	}

	newApi := &createApiKey{
		Name:          parameter.ApiKeyName,
		TenantID:      currentTenantID,
		ScopesList:    scopesCreateList,
		StartValidity: time.Now().Format(time.DateOnly),
		Description:   parameter.ApiKeyDescription,
	}

	if parameter.ApiKeyExpiration != "" {
		if _, err = time.Parse(time.DateOnly, parameter.ApiKeyExpiration); err != nil {
			return nil, fmt.Errorf("invalid date format for expiration, use YYYY-MM-DD")
		}
		newApi.EndValidity = parameter.ApiKeyExpiration
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(newApi)
	if err != nil {
		return nil, err
	}

	config := auth.GetConfig()
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
