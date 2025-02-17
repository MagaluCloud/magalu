/*
Executor: list

# Description

# List user clients

import "github.com/MagaluCloud/magalu/mgc/lib/products/auth/clients"
*/
package clients

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type ListResultItem struct {
	AccessTokenExpiration            *int                         `json:"access_token_expiration,omitempty"`
	AlwaysRequireLogin               bool                         `json:"always_require_login"`
	Audiences                        *ListResultItemAudiences     `json:"audiences,omitempty"`
	BackchannelLogoutSessionEnabled  bool                         `json:"backchannel_logout_session_enabled"`
	BackchannelLogoutUri             *string                      `json:"backchannel_logout_uri,omitempty"`
	ClientApprovalStatus             *string                      `json:"client_approval_status,omitempty"`
	ClientId                         *string                      `json:"client_id,omitempty"`
	ClientPrivacyTermUrl             *string                      `json:"client_privacy_term_url,omitempty"`
	Description                      *string                      `json:"description,omitempty"`
	Icon                             *string                      `json:"icon,omitempty"`
	Name                             *string                      `json:"name,omitempty"`
	OidcAudience                     *ListResultItemOidcAudience  `json:"oidc_audience,omitempty"`
	RedirectUris                     *ListResultItemRedirectUris  `json:"redirect_uris,omitempty"`
	RefreshTokenCustomExpiresEnabled bool                         `json:"refresh_token_custom_expires_enabled"`
	RefreshTokenExpiration           *int                         `json:"refresh_token_expiration,omitempty"`
	Scopes                           *ListResultItemScopes        `json:"scopes,omitempty"`
	ScopesDefault                    *ListResultItemScopesDefault `json:"scopes_default,omitempty"`
	TermOfUse                        *string                      `json:"term_of_use,omitempty"`
	Uuid                             *string                      `json:"uuid,omitempty"`
}

type ListResultItemAudiences []string

type ListResultItemOidcAudience []string

type ListResultItemRedirectUris []string

type ListResultItemScopes []string

type ListResultItemScopesDefault []string

type ListResult []ListResultItem

func (s *service) List() (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/auth/clients/list"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ListContext(
	ctx context.Context,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/auth/clients/list"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
