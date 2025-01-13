/*
Executor: create

# Description

Create new client (Oauth Application)

import "github.com/MagaluCloud/magalu/mgc/lib/products/auth/clients"
*/
package clients

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type CreateParameters struct {
	AccessTokenExpiration            *int    `json:"access-token-expiration,omitempty"`
	AlwaysRequireLogin               *bool   `json:"always-require-login,omitempty"`
	Audiences                        *string `json:"audiences,omitempty"`
	BackchannelLogoutSessionEnabled  *bool   `json:"backchannel-logout-session-enabled,omitempty"`
	BackchannelLogoutUri             *string `json:"backchannel-logout-uri,omitempty"`
	ClientPrivacyTermUrl             string  `json:"client-privacy-term-url"`
	Description                      string  `json:"description"`
	Icon                             *string `json:"icon,omitempty"`
	Name                             string  `json:"name"`
	OidcAudiences                    *string `json:"oidc-audiences,omitempty"`
	Reason                           string  `json:"reason"`
	RedirectUris                     string  `json:"redirect-uris"`
	RefreshTokenCustomExpiresEnabled *bool   `json:"refresh-token-custom-expires-enabled,omitempty"`
	RefreshTokenExpiration           *int    `json:"refresh-token-expiration,omitempty"`
	TermsOfUse                       string  `json:"terms-of-use"`
}

type CreateResult struct {
	ClientId     *string `json:"client_id,omitempty"`
	ClientSecret *string `json:"client_secret,omitempty"`
	Uuid         *string `json:"uuid,omitempty"`
}

func (s *service) Create(
	parameters CreateParameters,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/auth/clients/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) CreateContext(
	ctx context.Context,
	parameters CreateParameters,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/auth/clients/create"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
