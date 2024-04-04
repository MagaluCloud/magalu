/*
Executor: update

# Description

Update a client (Oauth Application)

import "magalu.cloud/lib/products/auth/clients"
*/
package clients

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UpdateParameters struct {
	AccessTokenExpiration            int    `json:"access-token-expiration,omitempty"`
	AlwaysRequireLogin               bool   `json:"always-require-login,omitempty"`
	Audiences                        string `json:"audiences,omitempty"`
	BackchannelLogoutSessionEnabled  bool   `json:"backchannel-logout-session-enabled,omitempty"`
	BackchannelLogoutUri             string `json:"backchannel-logout-uri,omitempty"`
	ClientPrivacyTermUrl             string `json:"client-privacy-term-url,omitempty"`
	Description                      string `json:"description,omitempty"`
	Icon                             string `json:"icon,omitempty"`
	Id                               string `json:"id"`
	Name                             string `json:"name,omitempty"`
	OidcAudiences                    string `json:"oidc-audiences,omitempty"`
	Reason                           string `json:"reason,omitempty"`
	RedirectUris                     string `json:"redirect-uris,omitempty"`
	RefreshTokenCustomExpiresEnabled bool   `json:"refresh-token-custom-expires-enabled,omitempty"`
	RefreshTokenExpiration           int    `json:"refresh-token-expiration,omitempty"`
	TermsOfUse                       string `json:"terms-of-use,omitempty"`
}

type UpdateResult struct {
	ClientId string `json:"client_id,omitempty"`
	Uuid     string `json:"uuid,omitempty"`
}

func Update(
	client *mgcClient.Client,
	ctx context.Context,
	parameters UpdateParameters,
) (
	result UpdateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/auth/clients/update"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UpdateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[UpdateResult](r)
}

// TODO: links
// TODO: related
