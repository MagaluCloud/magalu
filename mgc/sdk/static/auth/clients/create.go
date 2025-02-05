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
	//required by spec
	Name         string `json:"name" jsonschema:"description=Name of new client,example=Client Name" mgc:"positional"`
	Description  string `json:"description" jsonschema:"description=Description of new client,example=Client description" mgc:"positional"`
	RedirectURIs string `json:"redirect_uris" jsonschema:"description=Redirect URIs (separated by space)" mgc:"positional"`

	//required
	BackchannelLogoutSessionEnabled *bool   `json:"backchannel_logout_session,omitempty" jsonschema:"description=Client requires backchannel logout session,example=false" mgc:"positional"`
	ClientTermsUrl                  string  `json:"client_term_url" jsonschema:"description=URL to terms of use" mgc:"positional"`
	ClientPrivacyTermUrl            string  `json:"client_privacy_term_url" jsonschema:"description=URL to privacy term" mgc:"positional"`
	Audiences                       *string `json:"audiences,omitempty" jsonschema:"description=Client audiences (separated by space),example=public" mgc:"positional"`

	//optional by spec
	Email                            *string `json:"email,omitempty" jsonschema:"description=Email of new client,example=client@example.com" mgc:"positional"`
	Reason                           *string `json:"request_reason,omitempty" jsonschema:"description=Note to inform the reason for creating the client. Will help with the application approval process" mgc:"positional"`
	Icon                             *string `json:"icon,omitempty" jsonschema:"description=URL for client icon" mgc:"positional"`
	AccessTokenExp                   *int    `json:"access_token_expiration,omitempty" jsonschema:"description=Access token expiration (in seconds),example=7200" mgc:"positional"`
	AlwaysRequireLogin               *bool   `json:"always_require_login,omitempty" jsonschema:"description=Must ignore active Magalu ID session and always require login,example=false" mgc:"positional"`
	BackchannelLogoutUri             *string `json:"backchannel_logout_uri,omitempty" jsonschema:"description=Backchannel logout URI" mgc:"positional"`
	OidcAudience                     *string `json:"oidc_audience,omitempty" jsonschema:"description=OIDC audience (separated by space),example=public" mgc:"positional"`
	RefreshTokenCustomExpiresEnabled *bool   `json:"refresh_token_custom_expires_enabled,omitempty" jsonschema:"description=Use custom value for refresh token expiration,example=false" mgc:"positional"`
	RefreshTokenExp                  *int    `json:"refresh_token_exp,omitempty" jsonschema:"description=Custom refresh token expiration value (in seconds),example=15778476" mgc:"positional"`
	SupportUrl                       *string `json:"support_url,omitempty" jsonschema:"description=Support URL" mgc:"positional"`
	GrantTypes                       *string `json:"grant_types,omitempty" jsonschema:"description=Grant types the client can request for token generation (separated by space)" mgc:"positional"`
	// If BackchannelLogoutSessionEnabled is true, BackchannelLogoutUri is required
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

	if parameter.BackchannelLogoutSessionEnabled != nil && *parameter.BackchannelLogoutSessionEnabled {
		if parameter.BackchannelLogoutUri == nil {
			return nil, fmt.Errorf("programming error: BackchannelLogoutUri is required when BackchannelLogoutSessionEnabled is true")
		}
	}

	config := auth.GetConfig()

	// REQUIRED
	clientPayload := createClient{
		Name:        parameter.Name,
		Description: parameter.Description,
	}
	clientPayload.RedirectURIs = stringToSlice(parameter.RedirectURIs, " ", true)

	if parameter.Reason == nil {
		parameter.Reason = new(string)
		*parameter.Reason = "Created by MGCCLI"
	}

	clientPayload.ClientTermUrl = parameter.ClientTermsUrl
	clientPayload.ClientPrivacyTermUrl = parameter.ClientPrivacyTermUrl

	// FIXED SCOPES
	clientPayload.Scopes = []createClientScopes{{
		UUID:     config.PublicClientsScopeIDs["openid"],
		Reason:   *parameter.Reason,
		Optional: true,
	}, {
		UUID:     config.PublicClientsScopeIDs["profile"],
		Reason:   *parameter.Reason,
		Optional: true,
	}}

	//OPTIONAL
	if parameter.BackchannelLogoutSessionEnabled != nil {
		clientPayload.BackchannelLogoutSessionEnabled = parameter.BackchannelLogoutSessionEnabled
		clientPayload.BackchannelLogoutUri = parameter.BackchannelLogoutUri
	}

	if parameter.Audiences != nil {
		clientPayload.Audience = stringToSlice(*parameter.Audiences, " ", true)
	}

	if parameter.Email != nil {
		clientPayload.Email = parameter.Email
	}

	if parameter.Reason != nil {
		clientPayload.Reason = *parameter.Reason
	}

	if parameter.Icon != nil {
		clientPayload.Icon = parameter.Icon
	}

	clientPayload.AccessTokenExp = &default_access_token_expiraion
	if parameter.AccessTokenExp != nil {
		clientPayload.AccessTokenExp = parameter.AccessTokenExp
	}

	if parameter.AlwaysRequireLogin != nil {
		clientPayload.AlwaysRequireLogin = parameter.AlwaysRequireLogin
	}

	if parameter.OidcAudience != nil {
		clientPayload.OidcAudience = stringToSlice(*parameter.OidcAudience, " ", true)
	}

	if parameter.RefreshTokenCustomExpiresEnabled != nil {
		clientPayload.RefreshTokenCustomExpiresEnabled = parameter.RefreshTokenCustomExpiresEnabled
	}

	if parameter.RefreshTokenExp != nil {
		clientPayload.RefreshTokenExp = parameter.RefreshTokenExp
	}

	if parameter.SupportUrl != nil {
		clientPayload.SupportUrl = parameter.SupportUrl
	}

	if parameter.GrantTypes != nil {
		clientPayload.GrantTypes = stringToSlice(*parameter.GrantTypes, " ", true)
	}

	if parameter.Email != nil {
		clientPayload.Email = parameter.Email
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
