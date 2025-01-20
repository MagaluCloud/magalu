package api_key

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type listKeysParams struct {
	InvalidKeys bool `json:"invalid-keys" jsonschema:"description=Include Invalid Rekove and Expired Keys,default=false"`
}

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Scopes:      core.Scopes{scope_PA},
			Name:        "list",
			Summary:     "List your account API keys",
			Description: "This APIs Keys are from your account and can be used to authenticate in the Magalu Cloud",
		},
		list,
	)

	return exec
})

func list(ctx context.Context, parameter listKeysParams, _ struct{}) ([]*apiKeysResult, error) {
	keys, err := listFull(ctx, parameter.InvalidKeys)
	if err != nil {
		return nil, err
	}

	var result []*apiKeysResult
	for _, k := range keys {
		result = append(result, k.ToResult())
	}

	return result, nil
}

func listFull(ctx context.Context, invalidKeys bool) ([]*apiKeys, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("programming error: could not get auth configuration from context")
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, auth.GetConfig().ApiKeysUrlV1, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	httpClient := auth.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: could not get HTTP Client from context")
	}

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	var finallyResult []*apiKeys
	if resp.StatusCode == http.StatusNoContent {
		return finallyResult, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result []*apiKeys
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, y := range result {

		if !invalidKeys && y.RevokedAt != nil {
			continue
		}

		if !invalidKeys && y.EndValidity != nil {
			expDate, _ := time.Parse(time.RFC3339, *y.EndValidity)
			if expDate.Before(time.Now()) {
				continue
			}
		}

		finallyResult = append(finallyResult, y)
	}
	return finallyResult, nil
}
