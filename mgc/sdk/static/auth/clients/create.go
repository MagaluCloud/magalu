package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"

	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
)

type createParams struct {
	Name                             string  `json:"name" jsonschema:"description=Name of new client,example=Client Name" mgc:"positional"`
	Description                      string  `json:"description" jsonschema:"description=Description of new client,example=Client description" mgc:"positional"`
	RedirectURIs                     string  `json:"redirect-uris" jsonschema:"description=Redirect URIs (separated by space)" mgc:"positional"`
	Icon                             *string `json:"icon,omitempty" jsonschema:"description=URL for client icon" mgc:"positional"`
	AccessTokenExp                   *int    `json:"access-token-expiration,omitempty" jsonschema:"description=Access token expiration (in seconds),example=7200" mgc:"positional"`
	AlwaysRequireLogin               *bool   `json:"always-require-login,omitempty" jsonschema:"description=Must ignore active Magalu ID session and always require login,example=false" mgc:"positional"`
	ClientPrivacyTermUrl             string  `json:"client-privacy-term-url" jsonschema:"description=URL to privacy term" mgc:"positional"`
	TermsOfUse                       string  `json:"terms-of-use" jsonschema:"description=URL to terms of use" mgc:"positional"`
	Audience                         string  `json:"audiences,omitempty" jsonschema:"description=Client audiences (separated by space),example=public" mgc:"positional"`
	BackchannelLogoutSessionEnabled  *bool   `json:"backchannel-logout-session-enabled,omitempty" jsonschema:"description=Client requires backchannel logout session,example=false" mgc:"positional"`
	BackchannelLogoutUri             *string `json:"backchannel-logout-uri,omitempty" jsonschema:"description=Backchannel logout URI" mgc:"positional"`
	OidcAudience                     *string `json:"oidc-audiences,omitempty" jsonschema:"description=Audiences for ID token, should be the Client ID values" mgc:"positional"`
	RefreshTokenCustomExpiresEnabled *bool   `json:"refresh-token-custom-expires-enabled,omitempty" jsonschema:"description=Use custom value for refresh token expiration,example=false" mgc:"positional"`
	RefreshTokenExp                  *int    `json:"refresh-token-expiration,omitempty" jsonschema:"description=Custom refresh token expiration value (in seconds),example=15778476" mgc:"positional"`
	Reason                           string  `json:"reason" jsonschema:"description=Note to inform the reason for creating the client. Will help with the application approval process" mgc:"positional"`
}

var getCreate = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Scopes:      core.Scopes{scope_PA},
			Name:        "create",
			Description: "Create new client (Oauth Application)",
		},
		create,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Client created successfully! We'll analise your requisition and approve your client. You can check the approval status using client list command.\nClient ID: {{.client_id}}\nClient Secret: {{.client_secret}}\n"
	})
})

func create(ctx context.Context, parameter createParams, _ struct{}) (*createClientResult, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve auth configuration from context")
	}

	httpClient := auth.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve HTTP Client from context")
	}

	config := auth.GetConfig()

	clientPayload := createClient{
		Name:                             parameter.Name,
		Description:                      parameter.Description,
		Icon:                             parameter.Icon,
		ClientTermUrl:                    parameter.TermsOfUse,
		ClientPrivacyTermUrl:             parameter.ClientPrivacyTermUrl,
		AlwaysRequireLogin:               parameter.AlwaysRequireLogin,
		BackchannelLogoutSessionEnabled:  parameter.BackchannelLogoutSessionEnabled,
		BackchannelLogoutUri:             parameter.BackchannelLogoutUri,
		RefreshTokenCustomExpiresEnabled: parameter.RefreshTokenCustomExpiresEnabled,
		RefreshTokenExp:                  parameter.RefreshTokenExp,
		Reason:                           parameter.Reason,
	}

	clientPayload.Scopes = []createClientScopes{{
		UUID:     config.PublicClientsScopeIDs["openid"],
		Reason:   parameter.Reason,
		Optional: true,
	}, {
		UUID:     config.PublicClientsScopeIDs["profile"],
		Reason:   parameter.Reason,
		Optional: true,
	}}

	clientPayload.RedirectURIs = stringToSlice(parameter.RedirectURIs, " ", true)

	if len(parameter.Audience) > 0 {
		clientPayload.Audience = stringToSlice(parameter.Audience, " ", true)
	}

	if parameter.OidcAudience != nil {
		clientPayload.OidcAudience = stringToSlice(*parameter.OidcAudience, " ", true)
	}

	if parameter.AccessTokenExp == nil {
		clientPayload.AccessTokenExp = &default_access_token_expiraion
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(clientPayload)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, config.PublicClientsUrl, &buf)
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
	var result createClientResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
