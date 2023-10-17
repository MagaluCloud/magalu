package auth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/afero"
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var dummyConfig Config = Config{
	ClientId:       "client-id",
	RedirectUri:    "redirect-uri",
	LoginUrl:       "login-url",
	TokenUrl:       "token-url",
	ValidationUrl:  "validation-url",
	RefreshUrl:     "refresh-url",
	TenantsListUrl: "tenant-list-url",
	Scopes:         []string{"test"},
}

var dummyConfigResult *ConfigResult = &ConfigResult{
	AccessToken:     "access-token",
	RefreshToken:    "refresh-token",
	CurrentTenantID: "tenant-id",
	CurrentEnv:      "test",
}

var dummyConfigResultYaml = []byte(`---
access_token: "access-token"
refresh_token: "refresh-token"
current_tenant_id: "tenant-id"
current_environment: "test"
`)

type mockTransport struct {
	statusCode        int
	responseBody      io.ReadCloser
	shouldReturnError bool
}

func (o mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	if o.shouldReturnError {
		return nil, fmt.Errorf("test error")
	}
	return &http.Response{StatusCode: o.statusCode, Body: o.responseBody}, nil
}

func TestNew(t *testing.T) {
	filename, err := utils.BuildMGCFilePath(authFilename)
	if err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	type test struct {
		name           string
		authFileData   []byte
		envAccessToken string
		expected       *ConfigResult
	}

	tests := []test{
		{
			name:           "empty auth file",
			authFileData:   []byte{},
			envAccessToken: "",
			expected:       &ConfigResult{},
		},
		{
			name:           "non empty auth file",
			authFileData:   dummyConfigResultYaml,
			envAccessToken: "",
			expected:       dummyConfigResult,
		},
		{
			name:           "non empty auth file with env var",
			authFileData:   dummyConfigResultYaml,
			envAccessToken: "env-access-token",
			expected: &ConfigResult{
				AccessToken:     "env-access-token",
				RefreshToken:    dummyConfigResult.RefreshToken,
				CurrentTenantID: dummyConfigResult.CurrentTenantID,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			err := afero.WriteFile(fs, filename, tc.authFileData, 0644)
			if err != nil {
				t.Errorf("expected err == nil, found: %v", err)
			}

			t.Setenv("MGC_SDK_ACCESS_TOKEN", tc.envAccessToken)

			auth := New(dummyConfig, &http.Client{Transport: mockTransport{}}, fs)

			if auth.accessToken != tc.expected.AccessToken {
				t.Errorf("expected auth.accessToken == %v, found: %v", tc.expected.AccessToken, auth.accessToken)
			}
			if auth.refreshToken != tc.expected.RefreshToken {
				t.Errorf("expected auth.refreshToken == '', found: %v", auth.refreshToken)
			}
			if auth.currentTenantId != tc.expected.CurrentTenantID {
				t.Errorf("expected auth.currentTenantId == '', found: %v", auth.currentTenantId)
			}
		})
	}

	t.Run("no config file", func(t *testing.T) {
		auth := New(dummyConfig, &http.Client{Transport: mockTransport{}}, afero.NewMemMapFs())
		if auth.accessToken != "" {
			t.Errorf("expected auth.accessToken == '', found: %v", auth.accessToken)
		}
		if auth.refreshToken != "" {
			t.Errorf("expected auth.refreshToken == '', found: %v", auth.refreshToken)
		}
		if auth.currentTenantId != "" {
			t.Errorf("expected auth.currentTenantId == '', found: %v", auth.currentTenantId)
		}
	})
}

func TestSetTokens(t *testing.T) {
	var dummyLoginResult *LoginResult = &LoginResult{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	auth := New(dummyConfig, &http.Client{Transport: mockTransport{}}, afero.NewMemMapFs())

	if err := auth.SetTokens(dummyLoginResult); err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	if auth.accessToken != dummyLoginResult.AccessToken {
		t.Errorf("expected auth.accessToken = %s, found: %v", dummyLoginResult.AccessToken, auth.accessToken)
	}
	if auth.refreshToken != dummyLoginResult.RefreshToken {
		t.Errorf("expected auth.refreshToken = %s, found: %v", dummyLoginResult.RefreshToken, auth.refreshToken)
	}
}

func TestSetAccessKey(t *testing.T) {
	accessKeyId := "MyAccessKeyIdTest"
	secretAccessKey := "MySecretAccessKeyTeste"
	fs := afero.NewMemMapFs()
	currentAuth := New(dummyConfig, &http.Client{Transport: mockTransport{}}, fs)

	if err := currentAuth.SetAccessKey(accessKeyId, secretAccessKey); err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	auths := []*Auth{
		// Current auth
		currentAuth,
		// New auth reading from file
		New(dummyConfig, &http.Client{Transport: mockTransport{}}, fs),
	}

	for i, auth := range auths {
		if auth.accessKeyId != accessKeyId {
			t.Errorf("authIndex %v expected auth.accessKeyId = %s, found: %v", i, accessKeyId, auth.accessKeyId)
		}
		if auth.secretAccessKey != secretAccessKey {
			t.Errorf("authIndex %v expected auth.secretAccessKey = %s, found: %v", i, secretAccessKey, auth.secretAccessKey)
		}
	}
}

func TestRequestAuthTokenWithAuthorizationCode(t *testing.T) {
	filename, err := utils.BuildMGCFilePath(authFilename)
	if err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, filename, dummyConfigResultYaml, 0644); err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	type test struct {
		name        string
		transport   mockTransport
		verifier    *codeVerifier
		expected    LoginResult
		expectedErr bool
	}

	tests := []test{
		{
			name:        "code verifier == nil",
			transport:   mockTransport{},
			verifier:    nil,
			expected:    LoginResult{AccessToken: dummyConfigResult.AccessToken, RefreshToken: dummyConfigResult.RefreshToken},
			expectedErr: true,
		},
		{
			name: "invalid login result",
			transport: mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{`))),
			},
			verifier:    &codeVerifier{},
			expected:    LoginResult{AccessToken: dummyConfigResult.AccessToken, RefreshToken: dummyConfigResult.RefreshToken},
			expectedErr: true,
		},
		{
			name:        "request ended with error",
			transport:   mockTransport{shouldReturnError: true},
			verifier:    &codeVerifier{},
			expected:    LoginResult{AccessToken: dummyConfigResult.AccessToken, RefreshToken: dummyConfigResult.RefreshToken},
			expectedErr: true,
		},
		{
			name: "bad request",
			transport: mockTransport{
				statusCode:   http.StatusBadRequest,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte{})),
			},
			verifier:    &codeVerifier{},
			expected:    LoginResult{AccessToken: dummyConfigResult.AccessToken, RefreshToken: dummyConfigResult.RefreshToken},
			expectedErr: false,
		},
		{
			name: "valid login result",
			transport: mockTransport{
				statusCode: http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{
					"access_token": "ac-token",
					"refresh_token": "rf-token"
				}`))),
			},
			verifier:    &codeVerifier{},
			expected:    LoginResult{AccessToken: "ac-token", RefreshToken: "rf-token"},
			expectedErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			auth := New(dummyConfig, &http.Client{Transport: tc.transport}, fs)
			auth.codeVerifier = tc.verifier

			err := auth.RequestAuthTokenWithAuthorizationCode(context.Background(), "")
			hasErr := err != nil

			if hasErr != tc.expectedErr {
				t.Errorf("expected error == %v", tc.expectedErr)
			}
			if auth.accessToken != tc.expected.AccessToken {
				t.Errorf("expected accessToken == %v, found: %v", tc.expected.AccessToken, auth.accessToken)
			}
			if auth.refreshToken != tc.expected.RefreshToken {
				t.Errorf("expected refreshToken == %v, found: %v", tc.expected.RefreshToken, auth.refreshToken)
			}
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	filename, err := utils.BuildMGCFilePath(authFilename)
	if err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, filename, dummyConfigResultYaml, 0644); err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	type test struct {
		name           string
		transport      mockTransport
		expectedResult LoginResult
		expectedErr    bool
	}

	testErr := []test{
		{
			name:           "request ended with error",
			transport:      mockTransport{shouldReturnError: true},
			expectedResult: LoginResult{AccessToken: dummyConfigResult.AccessToken, RefreshToken: dummyConfigResult.RefreshToken},
			expectedErr:    true,
		},
		{
			name: "invalid validation result",
			transport: mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{`))),
			},
			expectedResult: LoginResult{AccessToken: dummyConfigResult.AccessToken, RefreshToken: dummyConfigResult.RefreshToken},
			expectedErr:    true,
		},
		{
			name: "bad request",
			transport: mockTransport{
				statusCode:   http.StatusBadRequest,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte{})),
			},
			expectedResult: LoginResult{AccessToken: dummyConfigResult.AccessToken, RefreshToken: dummyConfigResult.RefreshToken},
			expectedErr:    true,
		},
		{
			name: "active validation result",
			transport: mockTransport{
				statusCode: http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{
					"active": true
				}`))),
			},
			expectedResult: LoginResult{AccessToken: dummyConfigResult.AccessToken, RefreshToken: dummyConfigResult.RefreshToken},
			expectedErr:    false,
		},
	}

	for _, tcErr := range testErr {
		t.Run(tcErr.name, func(t *testing.T) {
			auth := New(dummyConfig, &http.Client{Transport: tcErr.transport}, fs)

			err := auth.ValidateAccessToken(context.Background())
			hasErr := err != nil

			if hasErr != tcErr.expectedErr {
				t.Errorf("expected err == %v", tcErr.expectedErr)
			}
			if auth.accessToken != tcErr.expectedResult.AccessToken {
				t.Errorf("expected auth.accessToken = %v, found: %v", tcErr.expectedResult.AccessToken, auth.accessToken)
			}
			if auth.refreshToken != tcErr.expectedResult.RefreshToken {
				t.Errorf("expected auth.refreshToken = %v, found: %v", tcErr.expectedResult.RefreshToken, auth.refreshToken)
			}
		})
	}
}

func TestDoRefreshAccessToken(t *testing.T) {
	filename, err := utils.BuildMGCFilePath(authFilename)
	if err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, filename, dummyConfigResultYaml, 0644); err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	type test struct {
		name           string
		transport      mockTransport
		expectedResult string
		expectedErr    bool
	}

	tests := []test{
		{
			name:           "request ended with error",
			transport:      mockTransport{shouldReturnError: true},
			expectedResult: "access-token",
			expectedErr:    true,
		},
		{
			name: "bad request",
			transport: mockTransport{
				statusCode:   http.StatusBadRequest,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte{})),
			},
			expectedResult: "",
			expectedErr:    true,
		},
		{
			name: "invalid response json",
			transport: mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{`))),
			},
			expectedResult: "",
			expectedErr:    true,
		},
		{
			name: "valid response json",
			transport: mockTransport{
				statusCode: http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{
					"access_token": "ac-token",
					"refresh_token": "rf-token"
				}`))),
			},
			expectedResult: "ac-token",
			expectedErr:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			auth := New(dummyConfig, &http.Client{Transport: tc.transport}, fs)

			tk, err := auth.doRefreshAccessToken(context.Background())
			hasErr := err != nil

			if hasErr != tc.expectedErr {
				t.Errorf("expected err == %v", tc.expectedErr)
			}
			if tk != tc.expectedResult {
				t.Errorf("expected tk == %v, found: %v", tc.expectedResult, tk)
			}
		})
	}
}

func TestListTenants(t *testing.T) {
	filename, err := utils.BuildMGCFilePath(authFilename)
	if err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, filename, dummyConfigResultYaml, 0644); err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	type test struct {
		name           string
		transport      mockTransport
		expectedResult []*Tenant
		expectedErr    bool
	}

	tests := []test{
		{
			name:           "request ended with err",
			transport:      mockTransport{shouldReturnError: true},
			expectedResult: nil,
			expectedErr:    true,
		},
		{
			name: "bad request",
			transport: mockTransport{
				statusCode:   http.StatusBadRequest,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte{})),
			},
			expectedResult: nil,
			expectedErr:    true,
		},
		{
			name: "invalid tenant list",
			transport: mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{`))),
			},
			expectedResult: nil,
			expectedErr:    true,
		},
		{
			name: "empty tenant list",
			transport: mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`[]`))),
			},
			expectedResult: []*Tenant{},
			expectedErr:    false,
		},
		{
			name: "non empty tenant list",
			transport: mockTransport{
				statusCode: http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`[
					{
						"uuid": "1",
						"legal_name": "jon doe",
						"email": "jon.doe@profusion.mobi",
						"is_managed": false,
						"is_delegated": false
					},
					{
						"uuid": "2",
						"legal_name": "jon smith",
						"email": "jon.smith@profusion.mobi",
						"is_managed": false,
						"is_delegated": false
					}
				]`))),
			},
			expectedResult: []*Tenant{
				{UUID: "1", Name: "jon doe", Email: "jon.doe@profusion.mobi", IsManaged: false, IsDelegated: false},
				{UUID: "2", Name: "jon smith", Email: "jon.smith@profusion.mobi", IsManaged: false, IsDelegated: false},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			auth := New(dummyConfig, &http.Client{Transport: tc.transport}, fs)

			tLst, err := auth.ListTenants(context.Background())
			hasErr := err != nil

			if hasErr != tc.expectedErr {
				t.Errorf("expected err == %v", tc.expectedErr)
			}
			if !reflect.DeepEqual(tLst, tc.expectedResult) {
				t.Errorf("expected tLst == %v, found: %v", tc.expectedResult, tLst)
			}
		})
	}
}

func TestSelectTenant(t *testing.T) {
	filename, err := utils.BuildMGCFilePath(authFilename)
	if err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, filename, dummyConfigResultYaml, 0644); err != nil {
		t.Errorf("expected err == nil, found: %v", err)
	}

	type test struct {
		name           string
		transport      mockTransport
		expectedResult *TenantAuth
		expectedErr    bool
	}

	tests := []test{
		{
			name:           "request ended with error",
			transport:      mockTransport{shouldReturnError: true},
			expectedResult: nil,
			expectedErr:    true,
		},
		{
			name: "bad request",
			transport: mockTransport{
				statusCode:   http.StatusBadRequest,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte{})),
			},
			expectedResult: nil,
			expectedErr:    true,
		},
		{
			name: "invalid tenant result",
			transport: mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{`))),
			},
			expectedResult: nil,
			expectedErr:    true,
		},
		{
			name: "valid tenant result",
			transport: mockTransport{
				statusCode: http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{
					"id": "abc123",
					"access_token": "abc",
					"created_at": 0,
					"refresh_token": "def",
					"scope": "test"
				}`))),
			},
			expectedResult: &TenantAuth{
				ID:           "abc123",
				CreatedAt:    core.Time(time.Unix(int64(0), 0)),
				AccessToken:  "abc",
				RefreshToken: "def",
				Scope:        []string{"test"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			auth := New(dummyConfig, &http.Client{Transport: tc.transport}, fs)

			tnt, err := auth.SelectTenant(context.Background(), "abc123")
			hasErr := err != nil

			if hasErr != tc.expectedErr {
				t.Errorf("expected err == %v", tc.expectedErr)
			}
			if !reflect.DeepEqual(tnt, tc.expectedResult) {
				t.Errorf("expected tnt == %v, found: %v", tc.expectedResult, tnt)
			}
		})
	}
}
