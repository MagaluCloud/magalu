package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"magalu.cloud/core"
	"magalu.cloud/core/config"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/profile_manager"

	"github.com/golang-jwt/jwt/v5"
	"github.com/invopop/yaml"
	"golang.org/x/sync/singleflight"
)

const (
	minRetryWait    = 1 * time.Second
	maxRetryWait    = 10 * time.Second
	maxRetryCount   = 5
	refreshGroupKey = "refreshToken"
	authFilename    = "auth.yaml"
)

// contextKey is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type contextKey string

type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type validationResult struct {
	Active bool `json:"active"`
}

type ConfigResult struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	CurrentEnv      string `json:"current_environment"` // ignored - used just for compatibility
}

type Config struct {
	ClientId         string
	RedirectUri      string
	LoginUrl         string
	TokenUrl         string
	ValidationUrl    string
	RefreshUrl       string
	TenantsListUrl   string
	TokenExchangeUrl string
}

type Auth struct {
	httpClient      *http.Client
	profileManager  *profile_manager.ProfileManager
	configMap       map[string]Config
	accessToken     string
	refreshToken    string
	accessKeyId     string
	secretAccessKey string
	codeVerifier    *codeVerifier
	group           singleflight.Group
	mgcConfig       *config.Config
}

type Tenant struct {
	UUID        string `json:"uuid"`
	Name        string `json:"legal_name"`
	Email       string `json:"email"`
	IsManaged   bool   `json:"is_managed"`
	IsDelegated bool   `json:"is_delegated"`
}

type tenantResult struct {
	AccessToken  string `json:"access_token"`
	CreatedAt    int    `json:"created_at"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"scope_type"`
}

type TokenExchangeResult struct {
	TenantID     string    `json:"uuid"`
	CreatedAt    core.Time `json:"created_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Scope        []string  `json:"scope"`
}

type accessTokenClaims struct {
	jwt.RegisteredClaims
	TenantIDWithType string            `json:"tenant"`
	ScopesStr        core.ScopesString `json:"scope"`
}

type FailedRefreshAccessToken struct {
	Message string
}

func (e FailedRefreshAccessToken) Error() string {
	return e.Message
}

var authKey contextKey = "magalu.cloud/core/Authentication"

func NewContext(parentCtx context.Context, auth *Auth) context.Context {
	return context.WithValue(parentCtx, authKey, auth)
}
func FromContext(ctx context.Context) *Auth {
	a, _ := ctx.Value(authKey).(*Auth)
	return a
}

func New(configMap map[string]Config, client *http.Client, profileManager *profile_manager.ProfileManager, mgcConfig *config.Config) *Auth {
	newAuth := Auth{
		httpClient:     client,
		configMap:      configMap,
		codeVerifier:   nil,
		profileManager: profileManager,
		mgcConfig:      mgcConfig,
	}
	newAuth.InitTokensFromFile()

	return &newAuth
}

func (a *Auth) getConfig() Config {
	var env string
	err := a.mgcConfig.Get("env", &env)
	if err != nil {
		logger().Debugw(
			"getConfig couldn't get 'env' from config",
		)
		return a.configMap["default"]
	}

	c, ok := a.configMap[env]
	if !ok {
		logger().Debugw("getConfig couldn't find a valid config to the env", "env", env)
		return a.configMap["default"]
	}
	return c
}

/*
Returns the current user access token.
If token is empty, we might still have refresh token, try getting a new one.
It will either fail with error or return a valid non-empty access token
*/
func (o *Auth) AccessToken(ctx context.Context) (string, error) {
	if o.accessToken == "" {
		if _, err := o.RefreshAccessToken(ctx); err != nil {
			return "", err
		}
	}
	return o.accessToken, nil
}

func (o *Auth) BuiltInScopes() core.Scopes {
	return core.Scopes{
		"openid",
		"cpo:read",
		"cpo:write",
	}
}

func (o *Auth) RedirectUri() string {
	return o.getConfig().RedirectUri
}

func (o *Auth) TenantsListUrl() string {
	return o.getConfig().TenantsListUrl
}

func (o *Auth) TokenExchangeUrl() string {
	return o.getConfig().TokenExchangeUrl
}

func (o *Auth) currentAccessTokenClaims() (*accessTokenClaims, error) {
	if o.accessToken == "" {
		return &accessTokenClaims{}, nil
	}

	tokenClaims := &accessTokenClaims{}
	tokenParser := jwt.NewParser()

	_, _, err := tokenParser.ParseUnverified(o.accessToken, tokenClaims)
	if err != nil {
		return nil, err
	}

	return tokenClaims, nil
}

func (o *Auth) CurrentTenantID() (string, error) {
	claims, err := o.currentAccessTokenClaims()
	if err != nil {
		return "", err
	}

	tenantId := claims.TenantIDWithType
	// Dot is a separator, Tenant will be <TenantType>.<ID>. We only want the ID
	if dotIndex := strings.Index(tenantId, "."); dotIndex != -1 {
		tenantId = tenantId[dotIndex+1:]
	}
	return tenantId, nil
}

func (o *Auth) CurrentTenant(ctx context.Context, httpClient *http.Client) (*Tenant, error) {
	currentTenantId, err := o.CurrentTenantID()
	if err != nil {
		return nil, err
	}

	tenants, err := o.ListTenants(ctx, httpClient)
	if err != nil || len(tenants) == 0 {
		logger().Warnw("Failed to get detailed info about Tenant, returning only ID", "err", err)
		return &Tenant{UUID: currentTenantId}, nil
	}

	for _, tenant := range tenants {
		if tenant.UUID == currentTenantId {
			return tenant, nil
		}
	}

	logger().Warnw(
		"Tenant ID present in Access Token is not present in the list of Tenants",
		"id", currentTenantId,
		"allTenants", tenants,
	)
	return nil, fmt.Errorf("unable to find Tenant in Tenant list that matches the current Tenant ID - %s", currentTenantId)
}

func (o *Auth) CurrentScopesString() (core.ScopesString, error) {
	claims, err := o.currentAccessTokenClaims()
	if err != nil {
		return "", err
	}

	return claims.ScopesStr, nil
}

func (o *Auth) CurrentScopes() (core.Scopes, error) {
	scopesStr, err := o.CurrentScopesString()
	if err != nil {
		return nil, err
	}

	return scopesStr.AsScopes(), nil
}

func (o *Auth) AccessKeyPair() (accessKeyId, secretAccessKey string) {
	return o.accessKeyId, o.secretAccessKey
}

func (o *Auth) SetTokens(token *LoginResult) error {
	// Always update the tokens, this way the user can assume the Auth object is
	// up-to-date after this function, even in case of a persistance error
	o.accessToken = token.AccessToken
	o.refreshToken = token.RefreshToken

	return o.writeCurrentConfig()
}

func (o *Auth) SetAccessKey(id string, key string) error {
	o.accessKeyId = id
	o.secretAccessKey = key
	return o.writeCurrentConfig()
}

func (o *Auth) writeCurrentConfig() error {
	authResult := &ConfigResult{}
	authResult.AccessToken = o.accessToken
	authResult.RefreshToken = o.refreshToken
	authResult.AccessKeyId = o.accessKeyId
	authResult.SecretAccessKey = o.secretAccessKey
	return o.writeConfigFile(authResult)
}

func (o *Auth) InitTokensFromFile() {
	authResult, _ := o.readConfigFile()
	if authResult != nil {
		o.accessToken = authResult.AccessToken
		o.refreshToken = authResult.RefreshToken
		o.accessKeyId = authResult.AccessKeyId
		o.secretAccessKey = authResult.SecretAccessKey
	}

	if envVal := os.Getenv("MGC_SDK_ACCESS_TOKEN"); envVal != "" {
		o.accessToken = envVal
	}
}

func (o *Auth) CodeChallengeToURL(scopes core.Scopes) (*url.URL, error) {
	config := o.getConfig()
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
	query.Add("scope", string(scopes.AsScopesString()))
	query.Add("choose_tenants", "true")

	loginUrl.RawQuery = query.Encode()

	return loginUrl, nil
}

/** Creates a new request access token from authorization code request, be
 * mindful that the code verifier used in this request come from the last call
 * of `CodeChallengeToUrl` method. */
func (o *Auth) RequestAuthTokenWithAuthorizationCode(ctx context.Context, authCode string) error {
	if o.codeVerifier == nil {
		logger().Errorw("no code verification provided")
		return fmt.Errorf("no code verification provided, first execute a code challenge request")
	}
	config := o.getConfig()
	data := url.Values{}
	data.Set("client_id", config.ClientId)
	data.Set("redirect_uri", config.RedirectUri)
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("code_verifier", o.codeVerifier.value)

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, config.TokenUrl, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	logger().Infow("Will send request for Auth Code", "authCode", authCode)
	resp, err := o.httpClient.Do(r)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	var result LoginResult
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if err = o.SetTokens(&result); err != nil {
		return err
	}

	return nil
}

func (o *Auth) ValidateAccessToken(ctx context.Context) error {
	r, err := o.newValidateAccessTokenRequest(ctx)
	if err != nil {
		return err
	}

	resp, err := o.httpClient.Do(r)
	if err != nil {
		return fmt.Errorf("Could not validate Access Token: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	var result validationResult
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Active {
		_, err := o.RefreshAccessToken(ctx)
		return err
	}

	return nil
}

func (o *Auth) newValidateAccessTokenRequest(ctx context.Context) (*http.Request, error) {
	config := o.getConfig()
	data := url.Values{}
	data.Set("client_id", config.ClientId)
	data.Set("token_hint", "access_token")
	data.Set("token", o.accessToken)

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, config.ValidationUrl, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return r, err
}

func (o *Auth) RefreshAccessToken(ctx context.Context) (string, error) {
	_, err, _ := o.group.Do(refreshGroupKey, func() (any, error) {
		return o.doRefreshAccessToken(ctx)
	})
	if err != nil {
		return "", err
	}
	return o.accessToken, nil
}

func (o *Auth) doRefreshAccessToken(ctx context.Context) (string, error) {
	var err error
	var resp *http.Response

	r, err := o.newRefreshAccessTokenRequest(ctx)
	if err != nil {
		return "", err
	}

	for i := 0; i < maxRetryCount; i++ {
		resp, err = o.httpClient.Do(r)
		if err != nil {
			wait := mgcHttpPkg.DefaultBackoff(minRetryWait, maxRetryCount, i, resp)
			fmt.Printf("Refresh access token failed, retrying in %s\n", wait)
			time.Sleep(wait)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return "", mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
		}

		var result LoginResult
		defer resp.Body.Close()
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", err
		}
		if err = o.SetTokens(&result); err != nil {
			return "", err
		} else {
			return o.accessToken, nil
		}
	}

	msg := fmt.Sprintf("failed to refresh access token: %v", err)
	return o.accessToken, FailedRefreshAccessToken{Message: msg}
}

func (o *Auth) newRefreshAccessTokenRequest(ctx context.Context) (*http.Request, error) {
	if o.refreshToken == "" {
		return nil, fmt.Errorf("RefreshToken is not set")
	}

	config := o.getConfig()
	data := url.Values{}
	data.Set("client_id", config.ClientId)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", o.refreshToken)

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, config.RefreshUrl, strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return r, err
}

func (o *Auth) readConfigFile() (*ConfigResult, error) {
	var result ConfigResult
	authFile, err := o.profileManager.Current().Read(authFilename)
	if err != nil {
		logger().Debugw("unable to read from auth configuration file", "error", err)
		return nil, err
	}

	err = yaml.Unmarshal(authFile, &result)
	if err != nil {
		logger().Warnw("bad format auth configuration file", "error", err)
		return nil, err
	}

	return &result, nil
}

func (o *Auth) writeConfigFile(result *ConfigResult) error {
	yamlData, err := yaml.Marshal(result)
	if err != nil {
		logger().Warn("unable to persist auth data", "error", err)
		return err
	}

	return o.profileManager.Current().Write(authFilename, yamlData)
}

func (o *Auth) ListTenants(ctx context.Context, httpClient *http.Client) ([]*Tenant, error) {
	at, err := o.AccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get current access token. Did you forget to log in?")
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, o.getConfig().TenantsListUrl, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "Bearer "+at)
	r.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result []*Tenant
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (o *Auth) SelectTenant(ctx context.Context, id string, scopes core.ScopesString) (*TokenExchangeResult, error) {
	at, err := o.AccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get current access token: %w. Did you forget to log in?", err)
	}

	return o.runTokenExchange(ctx, at, id, scopes, o.httpClient)
}

func (o *Auth) SetScopes(ctx context.Context, scopes core.Scopes, client *http.Client) (*TokenExchangeResult, error) {
	at, err := o.AccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get current access token: %w. Did you forget to log in?", err)
	}

	currentTenantId, err := o.CurrentTenantID()
	if err != nil {
		return nil, fmt.Errorf("unable to get current tenant ID: %w", err)
	}
	return o.runTokenExchange(ctx, at, currentTenantId, scopes.AsScopesString(), client)
}

func (o *Auth) runTokenExchange(ctx context.Context, currentAt string, tenantId string, scopes core.ScopesString, client *http.Client) (*TokenExchangeResult, error) {
	data := map[string]any{
		"tenant": tenantId,
		"scopes": scopes,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonData)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, o.TokenExchangeUrl(), bodyReader)
	r.Header.Set("Authorization", "Bearer "+currentAt)
	r.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	payload := &tenantResult{}
	if err = json.NewDecoder(resp.Body).Decode(payload); err != nil {
		return nil, err
	}

	err = o.SetTokens(&LoginResult{
		AccessToken:  payload.AccessToken,
		RefreshToken: payload.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	createdAt := core.Time(time.Unix(int64(payload.CreatedAt), 0))

	return &TokenExchangeResult{
		AccessToken:  payload.AccessToken,
		CreatedAt:    createdAt,
		TenantID:     tenantId,
		RefreshToken: payload.RefreshToken,
		Scope:        strings.Split(payload.Scope, " "),
	}, nil
}
