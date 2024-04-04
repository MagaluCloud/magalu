package clients

import (
	"strings"
)

var default_access_token_expiraion = 7200

const (
	scope_PA             = "pa:cloud-cli:features"
	readClients_scope_PA = "pa:clients:read"
)

type createClientScopes struct {
	UUID     string `json:"id"`
	Reason   string `json:"request_reason"`
	Optional bool   `json:"optional"`
}

type createClient struct {
	Name                             string               `json:"name" jsonschema:"description=Name of new client,example=Client Name" mgc:"positional"`
	Description                      string               `json:"description" jsonschema:"description=Description of new client,example=Client description" mgc:"positional"`
	Scopes                           []createClientScopes `json:"scopes" jsonschema:"description=List of scopes (separated by space),example=openid profile" mgc:"positional"`
	RedirectURIs                     []string             `json:"redirect_uris" jsonschema:"description=Redirect URIs (separated by space)" mgc:"positional"`
	Icon                             *string              `json:"icon,omitempty" jsonschema:"description=URL for client icon" mgc:"positional"`
	AccessTokenExp                   *int                 `json:"access_token_exp,omitempty" jsonschema:"description=Access token expiration (in seconds),example=7200" mgc:"positional"`
	AlwaysRequireLogin               *bool                `json:"always_require_login,omitempty" jsonschema:"description=Must ignore active Magalu ID session and always require login,example=false" mgc:"positional"`
	ClientPrivacyTermUrl             string               `json:"client_privacy_term_url" jsonschema:"description=URL to privacy term" mgc:"positional"`
	ClientTermUrl                    string               `json:"client_term_url" jsonschema:"description=URL to terms of use" mgc:"positional"`
	Audience                         []string             `json:"audience,omitempty" jsonschema:"description=Client audiences (separated by space),example=public" mgc:"positional"`
	BackchannelLogoutSessionEnabled  *bool                `json:"backchannel_logout_session_required,omitempty" jsonschema:"description=Client requires backchannel logout session,example=false" mgc:"positional"`
	BackchannelLogoutUri             *string              `json:"backchannel_logout_uri,omitempty" jsonschema:"description=Backchannel logout URI" mgc:"positional"`
	OidcAudience                     []string             `json:"oidc_audience,omitempty" jsonschema:"description=Audiences for ID token, should be the Client ID values" mgc:"positional"`
	RefreshTokenCustomExpiresEnabled *bool                `json:"refresh_token_custom_expires_enabled,omitempty" jsonschema:"description=Use custom value for refresh token expiration,example=false" mgc:"positional"`
	RefreshTokenExp                  *int                 `json:"refresh_token_exp,omitempty" jsonschema:"description=Custom refresh token expiration value (in seconds),example=15778476" mgc:"positional"`
	Reason                           string               `json:"request_reason,omitempty" jsonschema:"description=Note to inform the reason for creating the client. Will help with the application approval process" mgc:"positional"`
}

type updateClient struct {
	Name                             *string  `json:"name" jsonschema:"description=Name of new client,example=Client Name" mgc:"positional"`
	Description                      *string  `json:"description" jsonschema:"description=Description of new client,example=Client description" mgc:"positional"`
	RedirectURIs                     []string `json:"redirect_uris" jsonschema:"description=Redirect URIs (separated by space)" mgc:"positional"`
	Icon                             *string  `json:"icon,omitempty" jsonschema:"description=URL for client icon" mgc:"positional"`
	AccessTokenExp                   *int     `json:"access_token_exp,omitempty" jsonschema:"description=Access token expiration (in seconds),example=7200" mgc:"positional"`
	AlwaysRequireLogin               *bool    `json:"always_require_login,omitempty" jsonschema:"description=Must ignore active Magalu ID session and always require login,example=false" mgc:"positional"`
	ClientPrivacyTermUrl             *string  `json:"client_privacy_term_url" jsonschema:"description=URL to privacy term" mgc:"positional"`
	ClientTermUrl                    *string  `json:"client_term_url" jsonschema:"description=URL to terms of use" mgc:"positional"`
	Audience                         []string `json:"audience,omitempty" jsonschema:"description=Client audiences (separated by space),example=public" mgc:"positional"`
	OidcAudience                     []string `json:"oidc_audience,omitempty" jsonschema:"description=Audiences for ID token, should be the Client ID values" mgc:"positional"`
	BackchannelLogoutSessionEnabled  *bool    `json:"backchannel_logout_session_required,omitempty" jsonschema:"description=Client requires backchannel logout session,example=false" mgc:"positional"`
	BackchannelLogoutUri             *string  `json:"backchannel_logout_uri,omitempty" jsonschema:"description=Backchannel logout URI" mgc:"positional"`
	RefreshTokenCustomExpiresEnabled *bool    `json:"refresh_token_custom_expires_enabled,omitempty" jsonschema:"description=Use custom value for refresh token expiration,example=false" mgc:"positional"`
	RefreshTokenExp                  *int     `json:"refresh_token_exp,omitempty" jsonschema:"description=Custom refresh token expiration value (in seconds),example=15778476" mgc:"positional"`
	Reason                           *string  `json:"request_reason,omitempty" jsonschema:"description=Note to inform the reason for creating the client. Will help with the application approval process" mgc:"positional"`
}

type createClientResult struct {
	UUID         string `json:"uuid,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

type updateClientResult struct {
	UUID     string `json:"uuid,omitempty"`
	ClientID string `json:"client_id,omitempty"`
}

type listClientResult struct {
	UUID                             string   `json:"uuid,omitempty"`
	ClientID                         string   `json:"client_id,omitempty"`
	Name                             string   `json:"name,omitempty"`
	Description                      string   `json:"description,omitempty"`
	Status                           string   `json:"client_approval_status,omitempty"`
	Scopes                           []string `json:"scopes,omitempty"`
	ScopesDefault                    []string `json:"scopes_default,omitempty"`
	TermOfUse                        string   `json:"term_of_use,omitempty"`
	ClientPrivacyTermUrl             string   `json:"client_privacy_term_url,omitempty"`
	Audiences                        []string `json:"audiences,omitempty"`
	OidcAudiences                    []string `json:"oidc_audience,omitempty"`
	AlwaysRequireLogin               bool     `json:"always_require_login"`
	BackchannelLogoutSessionEnabled  bool     `json:"backchannel_logout_session_enabled"`
	BackchannelLogoutUri             string   `json:"backchannel_logout_uri,omitempty"`
	RefreshTokenCustomExpiresEnabled bool     `json:"refresh_token_custom_expires_enabled"`
	RefreshTokenExp                  int      `json:"refresh_token_expiration,omitempty"`
	AccessTokenExp                   int      `json:"access_token_expiration,omitempty"`
	RedirectURIs                     []string `json:"redirect_uris,omitempty"`
	Icon                             string   `json:"icon,omitempty"`
}

type clients struct {
	UUID        string `json:"uuid,omitempty"`
	ClientID    string `json:"client_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Scopes      []struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"scopes,omitempty"`
	ScopesDefault []struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"scopes_default,omitempty"`
	TermOfUse                        string   `json:"client_term_url,omitempty"`
	ClientPrivacyTermUrl             string   `json:"client_privacy_term_url,omitempty"`
	Audience                         []string `json:"audience,omitempty"`
	OidcAudience                     []string `json:"oidc_audience,omitempty"`
	AlwaysRequireLogin               bool     `json:"always_require_login,omitempty"`
	BackchannelLogoutSessionRequired bool     `json:"backchannel_logout_session_required,omitempty"`
	BackchannelLogoutUri             string   `json:"backchannel_logout_uri,omitempty"`
	RefreshTokenCustomExpiresEnabled bool     `json:"refresh_token_custom_expires_enabled,omitempty"`
	RefreshTokenExp                  int      `json:"refresh_token_exp,omitempty"`
	AccessTokenExp                   int      `json:"access_token_exp,omitempty"`
	RedirectURIs                     []string `json:"redirect_uris,omitempty"`
	Icon                             string   `json:"icon,omitempty"`
}

func stringToSlice(s, sep string, shouldTrim bool) []string {
	entries := strings.Split(s, sep)

	result := make([]string, 0)
	if shouldTrim {
		for _, str := range entries {
			newValue := strings.TrimSpace(str)
			if newValue == "" {
				continue
			}
			result = append(result, newValue)
		}
	}

	return result
}
