package sharing

import (
	"context"
	"fmt"
	"net/http"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/utils"
)

type removeParams struct {
	DelegationID string `json:"uuid" jsonschema:"description=Shared uuid" mgc:"positional"`
}

var getRemove = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "remove",
			Description: "Remove access to my Account/Organization",
		},
		remove,
	)

	return executor
})

func remove(ctx context.Context, parameter removeParams, _ struct{}) (*[]Delegation, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve auth configuration from context")
	}

	httpClient := auth.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve HTTP Client from context")
	}

	config := auth.GetConfig()

	r, err := http.NewRequestWithContext(ctx, http.MethodDelete, config.DelegationUrl+"/"+parameter.DelegationID, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()

	var result []Delegation
	return &result, nil
}
