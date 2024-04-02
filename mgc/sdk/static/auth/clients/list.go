package clients

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

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Scopes:      core.Scopes{readClients_scope_PA},
			Name:        "list",
			Description: "List user clients",
		},
		list,
	)

	return exec
})

func list(ctx context.Context) ([]*listClientResult, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("programming error: could not get auth configuration from context")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, auth.GetConfig().ClientsV2Url, nil)
	if err != nil {
		return nil, err
	}

	httpClient := auth.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: could not get HTTP Client from context")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var clientsOfUser []*listClientResult
	if resp.StatusCode == http.StatusNoContent {
		return clientsOfUser, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, req)
	}

	defer resp.Body.Close()
	var clientsApiResponse []*clients
	if err = json.NewDecoder(resp.Body).Decode(&clientsApiResponse); err != nil {
		return nil, err
	}

	for _, client := range clientsApiResponse {
		clientResult := &listClientResult{
			UUID:                             client.UUID,
			ClientID:                         client.ClientID,
			Name:                             client.Name,
			Description:                      client.Description,
			Status:                           client.Status,
			TermOfUse:                        client.TermOfUse,
			ClientPrivacyTermUrl:             client.ClientPrivacyTermUrl,
			Audiences:                        client.Audience,
			AlwaysRequireLogin:               client.AlwaysRequireLogin,
			OidcAudiences:                    client.OidcAudience,
			BackchannelLogoutSessionEnabled:  client.BackchannelLogoutSessionRequired,
			BackchannelLogoutUri:             client.BackchannelLogoutUri,
			RefreshTokenCustomExpiresEnabled: client.RefreshTokenCustomExpiresEnabled,
			RefreshTokenExp:                  client.RefreshTokenExp,
			AccessTokenExp:                   client.AccessTokenExp,
			RedirectURIs:                     client.RedirectURIs,
			Icon:                             client.Icon,
		}

		for _, scope := range client.Scopes {
			clientResult.Scopes = append(clientResult.Scopes, scope.Name)
		}
		for _, scope := range client.ScopesDefault {
			clientResult.ScopesDefault = append(clientResult.ScopesDefault, scope.Name)
		}

		clientsOfUser = append(clientsOfUser, clientResult)
	}

	return clientsOfUser, nil
}
