package api_key

import (
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

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Scopes:      core.Scopes{"pa:api-keys:read"},
			Name:        "list",
			Description: "List valid Object Storage credentials",
		},
		list,
	)
})

func list(ctx context.Context) ([]*apiKeysResult, error) {

	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("could not get Auth from context")
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, auth.GetConfig().ApiKeysUrl, nil)
	if err != nil {
		return nil, err
	}
	token, err := auth.AccessToken(ctx)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", "Bearer "+token)
	r.Header.Set("Content-Type", "application/json")

	httpClient := mgcHttpPkg.ClientFromContext(ctx)
	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	var finallyResult []*apiKeysResult
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

		if y.RevokedAt != nil {
			continue
		}

		if y.EndValidity != nil {
			expDate, _ := time.Parse(time.RFC3339, *y.EndValidity)
			if expDate.Before(time.Now()) {
				continue
			}
		}

		for _, s := range y.Scopes {
			if s.Name != "*" && name_ObjectStorage != s.APIProduct.Name {
				continue
			}
			tenantName := y.Tenant.LegalName
			y.apiKeysResult.TenantName = &tenantName
			finallyResult = append(finallyResult, &y.apiKeysResult)
			break
		}
	}
	return finallyResult, nil

}
