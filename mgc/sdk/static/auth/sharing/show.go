package sharing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/utils"
)

type DelegationAccount struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DelegationTenant struct {
	LegalName string `json:"legal_name"`
}

type Delegation struct {
	UUID        string            `json:"uuid"`
	DelegatedAt string            `json:"delegated_at"`
	DelegatedBy DelegationAccount `json:"delegated_by"`
	DelegatedTo DelegationAccount `json:"delegated_to"`
	Tentant     DelegationTenant  `json:"tenant"`
}

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "show",
			Description: "Show people with access to my Account/Organization",
		},
		listDelegations,
	)

	return exec
})

func listDelegations(ctx context.Context) ([]*Delegation, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("programming error: unable to get auth from context")
	}
	res, err := ListDelegations(ctx)
	return res, err
}

func ListDelegations(ctx context.Context) ([]*Delegation, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	httpClient := auth.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to get HTTP Client from context")
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, auth.GetConfig().DelegationUrl, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		var result []*Delegation
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		return result, nil
	}

	if resp.StatusCode == http.StatusNoContent {
		var result []*Delegation
		return result, nil
	}

	return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)

}
