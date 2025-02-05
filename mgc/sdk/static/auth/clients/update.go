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

type updateParams struct {
	ID                               string  `json:"id" jsonschema:"description=UUID of client" mgc:"positional"`
	Name                             *string `json:"name,omitempty" jsonschema:"description=Name of new client,example=Client Name" mgc:"positional"`
	Description                      *string `json:"description,omitempty" jsonschema:"description=Description of new client,example=Client description" mgc:"positional"`
	RedirectURIs                     *string `json:"redirect_uris,omitempty" jsonschema:"description=Redirect URIs (separated by space)" mgc:"positional"`
	Icon                             *string `json:"icon,omitempty" jsonschema:"description=URL for client icon" mgc:"positional"`
	AccessTokenExp                   *int    `json:"access_token_expiration,omitempty" jsonschema:"description=Access token expiration (in seconds),example=7200" mgc:"positional"`
	AlwaysRequireLogin               *bool   `json:"always_require_login,omitempty" jsonschema:"description=Must ignore active Magalu ID session and always require login,example=false" mgc:"positional"`
	ClientPrivacyTermUrl             *string `json:"client_privacy_term_url,omitempty" jsonschema:"description=URL to privacy term" mgc:"positional"`
	ClientTermUrl                    *string `json:"client_term_url,omitempty" jsonschema:"description=URL to terms of use" mgc:"positional"`
	Audiences                        *string `json:"audiences,omitempty" jsonschema:"description=Client audiences (separated by space),example=public" mgc:"positional"`
	BackchannelLogoutSessionEnabled  *bool   `json:"backchannel_logout_session,omitempty" jsonschema:"description=Client requires backchannel logout session,example=false" mgc:"positional"`
	BackchannelLogoutUri             *string `json:"backchannel_logout_uri,omitempty" jsonschema:"description=Backchannel logout URI" mgc:"positional"`
	OidcAudience                     *string `json:"oidc_audience,omitempty" jsonschema:"description=Audiences for ID token, should be the Client ID values" mgc:"positional"`
	RefreshTokenCustomExpiresEnabled *bool   `json:"refresh_token_custom_expires_enabled,omitempty" jsonschema:"description=Use custom value for refresh token expiration,example=false" mgc:"positional"`
	RefreshTokenExp                  *int    `json:"refresh_token_exp,omitempty" jsonschema:"description=Custom refresh token expiration value (in seconds),example=15778476" mgc:"positional"`
	Reason                           *string `json:"request_reason,omitempty" jsonschema:"description=Note to inform the reason for creating the client. Will help with the application approval process" mgc:"positional"`
	SupportUrl                       *string `json:"support_url,omitempty" jsonschema:"description=URL for client support" mgc:"positional"`
}

var getUpdate = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Scopes:      core.Scopes{scope_PA},
			Name:        "update",
			Description: "Update a client (Oauth Application)",
		},
		update,
	)

	msg := "This operation may disable your client {{.parameters.id}} until updates are approved by the ID Magalu. Do you wish to continue?"

	cExecutor := core.NewConfirmableExecutor(
		executor,
		core.ConfirmPromptWithTemplate(msg),
	)

	return core.NewExecuteResultOutputOptions(cExecutor, func(exec core.Executor, result core.Result) string {
		return "template=Client updated successfully\nClient ID={{.client_id}}\n"
	})
})

func update(ctx context.Context, parameter updateParams, _ struct{}) (*updateClientResult, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve auth configuration from context")
	}

	httpClient := auth.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve HTTP Client from context")
	}

	config := auth.GetConfig()

	clientPayload := updateClient{
		Name:                             parameter.Name,
		Description:                      parameter.Description,
		Icon:                             parameter.Icon,
		ClientTermUrl:                    parameter.ClientTermUrl,
		ClientPrivacyTermUrl:             parameter.ClientPrivacyTermUrl,
		AlwaysRequireLogin:               parameter.AlwaysRequireLogin,
		BackchannelLogoutSessionEnabled:  parameter.BackchannelLogoutSessionEnabled,
		BackchannelLogoutUri:             parameter.BackchannelLogoutUri,
		AccessTokenExp:                   parameter.AccessTokenExp,
		RefreshTokenCustomExpiresEnabled: parameter.RefreshTokenCustomExpiresEnabled,
		RefreshTokenExp:                  parameter.RefreshTokenExp,
		Reason:                           parameter.Reason,
		SupportUrl:                       parameter.SupportUrl,
	}

	if parameter.RedirectURIs != nil {
		clientPayload.RedirectURIs = stringToSlice(*parameter.RedirectURIs, " ", true)
	}

	if parameter.Audiences != nil {
		clientPayload.Audience = stringToSlice(*parameter.Audiences, " ", true)
	}

	if parameter.OidcAudience != nil {
		clientPayload.OidcAudience = stringToSlice(*parameter.OidcAudience, " ", true)
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(clientPayload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", config.PublicClientsUrl, parameter.ID)

	r, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, &buf)
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
	var result updateClientResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
