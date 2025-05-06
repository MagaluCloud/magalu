package sharing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type createParams struct {
	DelegateToEmail string `json:"email" jsonschema:"description=User email" mgc:"positional"`
}

type Account struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	UUID  string `json:"uuid"`
}

type createPayload struct {
	DelegateToUUID string `json:"delegated_to" jsonschema:"description=User UUID" mgc:"positional"`
	ScopeID        string `json:"scope_id"`
	IsAdmin        bool   `json:"is_admin"`
}

type delegateResult struct {
	DelegationID      string `json:"delegation_id"`
	ScopeDelegationID string `json:"scope_delegation_id"`
}

var getCreate = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "admit",
			Description: "Admit access to my Account/Organization",
		},
		create,
	)

	return executor
})

func create(ctx context.Context, parameter createParams, _ struct{}) (*delegateResult, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve auth configuration from context")
	}

	httpClient := auth.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve HTTP Client from context")
	}

	config := auth.GetConfig()

	//////////////////////////////////////////////////////////////
	// Find Account ID
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, config.AccountsUrl+"?email="+parameter.DelegateToEmail, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response")
	}

	defer resp.Body.Close()
	var account_result []*Account
	if err = json.NewDecoder(resp.Body).Decode(&account_result); err != nil {
		return nil, err
	}

	if len(account_result) < 1 {
		return nil, fmt.Errorf("user account not found")
	}

	DelegateToUUID := account_result[0].UUID
	//////////////////////////////////////////////////////////////
	clientPayload := createPayload{
		DelegateToUUID: DelegateToUUID,
		ScopeID:        config.TotalScopeID,
		IsAdmin:        false,
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(clientPayload)
	if err != nil {
		return nil, err
	}

	r, err = http.NewRequestWithContext(ctx, http.MethodPost, config.DelegationUrl, &buf)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	resp, err = httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result delegateResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
