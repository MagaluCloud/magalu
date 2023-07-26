package core

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type loginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type validationResult struct {
	Active bool `json:"active"`
}

type AuthConfig struct {
	ClientId      string
	RedirectUri   string
	LoginUrl      string
	TokenUrl      string
	ValidationUrl string
	RefreshUrl    string
	Scopes        []string
}

type Auth struct {
	httpClient   *http.Client
	config       AuthConfig
	accessToken  string
	refreshToken string
	codeVerifier *codeVerifier
}

type authKey string

var key = authKey("auth")

func NewAuthContext(parentCtx context.Context, auth *Auth) context.Context {
	return context.WithValue(parentCtx, key, auth)
}
func AuthFromContext(ctx context.Context) *Auth {
	a, _ := ctx.Value(key).(*Auth)
	return a
}

func NewAuth(config AuthConfig, client *http.Client) *Auth {
	return &Auth{
		httpClient:   client,
		config:       config,
		codeVerifier: nil,
	}
}

func (o *Auth) AccessToken() string {
	// TODO: remove GetEnv and read from file
	if o.accessToken == "" {
		return os.Getenv("MGC_SDK_ACCESS_TOKEN")
	}
	return o.accessToken
}

func (o *Auth) RedirectUri() string {
	return o.config.RedirectUri
}

func (o *Auth) setToken(token *loginResult) error {
	// TODO: Persist the token in the disk, return error if something happens
	o.accessToken = token.AccessToken
	o.refreshToken = token.RefreshToken

	return nil
}

func (o *Auth) CodeChallengeToURL() (*url.URL, error) {
	config := o.config
	loginUrl, err := url.Parse(config.LoginUrl)
	if err != nil {
		return nil, err
	}
	codeVerifier, err := newVerifier()
	o.codeVerifier = codeVerifier
	if err != nil {
		return nil, err
	}

	query := loginUrl.Query()
	query.Add("response_type", "code")
	query.Add("client_id", config.ClientId)
	query.Add("redirect_uri", config.RedirectUri)
	query.Add("code_challenge", o.codeVerifier.CodeChallengeS256())
	query.Add("code_challenge_method", "S256")
	query.Add("scope", strings.Join(config.Scopes, " "))

	loginUrl.RawQuery = query.Encode()

	return loginUrl, nil
}

/** Creates a new request access token from authorization code request, be
 * mindful that the code verifier used in this request come from the last call
 * of `CodeChallengeToUrl` method. */
func (o *Auth) RequestAuthTokeWithAuthorizationCode(authCode string) error {
	if o.codeVerifier == nil {
		return errors.New("no code verification provided, first execute a code challenge request")
	}
	config := o.config
	data := url.Values{}
	data.Set("client_id", config.ClientId)
	data.Set("redirect_uri", config.RedirectUri)
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("code_verifier", o.codeVerifier.value)

	r, err := http.NewRequest(http.MethodPost, config.TokenUrl, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	resp, err := o.httpClient.Do(r)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	var result loginResult
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if err = o.setToken(&result); err != nil {
		return err
	}

	return nil
}

func (o *Auth) ValidateAccessToken() error {
	r, err := o.newValidateAccessTokenRequest()
	if err != nil {
		return err
	}

	resp, err := o.httpClient.Do(r)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	var result validationResult
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Active {
		return o.RefreshAccessToken()
	}

	return nil
}

func (o *Auth) newValidateAccessTokenRequest() (*http.Request, error) {
	config := o.config
	data := url.Values{}
	data.Set("client_id", config.ClientId)
	data.Set("token_hint", "access_token")
	data.Set("token", o.accessToken)

	r, err := http.NewRequest(http.MethodPost, config.ValidationUrl, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return r, err
}

func (o *Auth) RefreshAccessToken() error {
	r, err := o.newRefreshAccessTokenRequest()
	if err != nil {
		return err
	}

	resp, err := o.httpClient.Do(r)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	var result loginResult
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	o.accessToken = result.AccessToken
	o.refreshToken = result.RefreshToken
	return nil
}

func (o *Auth) newRefreshAccessTokenRequest() (*http.Request, error) {
	config := o.config
	data := url.Values{}
	data.Set("client_id", config.ClientId)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", o.refreshToken)

	r, err := http.NewRequest(http.MethodPost, config.RefreshUrl, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return r, err
}
