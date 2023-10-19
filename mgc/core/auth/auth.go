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

	"github.com/spf13/afero"
	"magalu.cloud/core"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/utils"

	"golang.org/x/sync/singleflight"
	"gopkg.in/yaml.v3"
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
	AccessToken     string `yaml:"access_token"`
	RefreshToken    string `yaml:"refresh_token"`
	CurrentTenantID string `yaml:"current_tenant_id"`
	AccessKeyId     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	CurrentEnv      string `yaml:"current_environment"` // ignored - used just for compatibility
}

type Config struct {
	ClientId         string
	RedirectUri      string
	LoginUrl         string
	TokenUrl         string
	ValidationUrl    string
	RefreshUrl       string
	TenantsListUrl   string
	TenantsSelectUrl string
	Scopes           []string
}

type Auth struct {
	httpClient      *http.Client
	config          Config
	configFile      string
	accessToken     string
	refreshToken    string
	currentTenantId string
	accessKeyId     string
	secretAccessKey string
	codeVerifier    *codeVerifier
	group           singleflight.Group
	fs              afero.Fs
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

type TenantAuth struct {
	ID           string    `json:"id"`
	CreatedAt    core.Time `json:"created_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Scope        []string  `json:"scope"`
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

func New(config Config, client *http.Client, fs afero.Fs) *Auth {
	filePath, err := utils.BuildMGCFilePath(authFilename)
	if err != nil {
		logger().Warnw("unable to locate auth configuration file", "error", err)
		return nil
	}

	newAuth := Auth{
		httpClient:   client,
		config:       config,
		configFile:   filePath,
		codeVerifier: nil,
		fs:           fs,
	}
	newAuth.InitTokensFromFile()

	return &newAuth
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

func (o *Auth) RedirectUri() string {
	return o.config.RedirectUri
}

func (o *Auth) TenantsListUrl() string {
	return o.config.TenantsListUrl
}

func (o *Auth) TenantsSelectUrl() string {
	return o.config.TenantsSelectUrl
}

func (o *Auth) CurrentTenantID() string {
	return o.currentTenantId
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

func (o *Auth) SetCurrentTenantID(id string) error {
	o.currentTenantId = id
	return o.writeCurrentConfig()
}

func (o *Auth) writeCurrentConfig() error {
	authResult := &ConfigResult{}
	authResult.AccessToken = o.accessToken
	authResult.RefreshToken = o.refreshToken
	authResult.CurrentTenantID = o.currentTenantId
	authResult.AccessKeyId = o.accessKeyId
	authResult.SecretAccessKey = o.secretAccessKey
	return o.writeConfigFile(authResult)
}

func (o *Auth) InitTokensFromFile() {
	authResult, _ := o.readConfigFile()
	if authResult != nil {
		o.accessToken = authResult.AccessToken
		o.refreshToken = authResult.RefreshToken
		o.currentTenantId = authResult.CurrentTenantID
		o.accessKeyId = authResult.AccessKeyId
		o.secretAccessKey = authResult.SecretAccessKey
	}

	if envVal := os.Getenv("MGC_SDK_ACCESS_TOKEN"); envVal != "" {
		o.accessToken = envVal
	}
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
	config := o.config
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
		return mgcHttpPkg.NewHttpErrorFromResponse(resp)
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
	config := o.config
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

func (o *Auth) doRefreshAccessToken(ctx context.Context) (core.Value, error) {
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
			return "", mgcHttpPkg.NewHttpErrorFromResponse(resp)
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

	config := o.config
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
	authFile, err := afero.ReadFile(o.fs, o.configFile)
	if err != nil {
		logger().Warnw("unable to read from auth configuration file", "error", err)
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

	err = afero.WriteFile(o.fs, o.configFile, yamlData, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (o *Auth) ListTenants(ctx context.Context) ([]*Tenant, error) {
	at, err := o.AccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get current access token. Did you forget to log in?")
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, o.config.TenantsListUrl, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "Bearer "+at)
	r.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp)
	}

	defer resp.Body.Close()
	var result []*Tenant
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (o *Auth) SelectTenant(ctx context.Context, id string) (*TenantAuth, error) {
	at, err := o.AccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get current access token. Did you forget to log in?")
	}

	data := map[string]any{"tenant": id}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonData)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, o.config.TenantsSelectUrl, bodyReader)
	r.Header.Set("Authorization", "Bearer "+at)
	r.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	resp, err := o.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp)
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

	if err := o.SetCurrentTenantID(id); err != nil {
		return nil, err
	}

	createdAt := core.Time(time.Unix(int64(payload.CreatedAt), 0))

	return &TenantAuth{
		AccessToken:  payload.AccessToken,
		CreatedAt:    createdAt,
		ID:           id,
		RefreshToken: payload.RefreshToken,
		Scope:        strings.Split(payload.Scope, " "),
	}, nil
}
