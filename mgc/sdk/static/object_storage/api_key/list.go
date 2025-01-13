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

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Scopes:      core.Scopes{scope_PA},
			Name:        "list",
			Description: "List valid Object Storage credentials",
		},
		list,
	)

	exec = core.NewHumanIdentifiableFieldsExecutor(exec, []string{"name"})

	return exec
})

func list(ctx context.Context) ([]*apiKeysResult, error) {
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
