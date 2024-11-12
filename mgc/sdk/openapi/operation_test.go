package openapi

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"magalu.cloud/core"
	"magalu.cloud/core/auth"
)

type mockAuth struct {
	mock.Mock
}

func (m *mockAuth) ApiKey(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *mockAuth) XTenantID(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *mockAuth) AccessToken(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *mockAuth) CurrentSecurityMethod() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockAuth) SetTokens(token *auth.LoginResult) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockAuth) RefreshAccessToken(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *mockAuth) ValidateAccessToken(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockAuth) CodeChallengeToURL(scopes core.Scopes) (*url.URL, error) {
	args := m.Called(scopes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*url.URL), args.Error(1)
}

func (m *mockAuth) RequestAuthTokenWithAuthorizationCode(ctx context.Context, authCode string) error {
	args := m.Called(ctx, authCode)
	return args.Error(0)
}

func (m *mockAuth) ListTenants(ctx context.Context) ([]*auth.Tenant, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auth.Tenant), args.Error(1)
}

func (m *mockAuth) SelectTenant(ctx context.Context, id string, scopes core.ScopesString) (*auth.TokenExchangeResult, error) {
	args := m.Called(ctx, id, scopes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenExchangeResult), args.Error(1)
}

func (m *mockAuth) CurrentTenant(ctx context.Context) (*auth.Tenant, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Tenant), args.Error(1)
}

func (m *mockAuth) CurrentTenantID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockAuth) SetScopes(ctx context.Context, scopes core.Scopes) (*auth.TokenExchangeResult, error) {
	args := m.Called(ctx, scopes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenExchangeResult), args.Error(1)
}

func (m *mockAuth) CurrentScopes() (core.Scopes, error) {
	args := m.Called()
	return args.Get(0).(core.Scopes), args.Error(1)
}

func (m *mockAuth) CurrentScopesString() (core.ScopesString, error) {
	args := m.Called()
	return args.Get(0).(core.ScopesString), args.Error(1)
}

func (m *mockAuth) SetAPIKey(apiKey string) error {
	args := m.Called(apiKey)
	return args.Error(0)
}

func (m *mockAuth) AccessKeyPair() (string, string) {
	args := m.Called()
	return args.String(0), args.String(1)
}

func (m *mockAuth) SetAccessKey(id string, key string) error {
	args := m.Called(id, key)
	return args.Error(0)
}

func (m *mockAuth) UnsetAccessKey() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockAuth) SetXTenantID(tenantId string) error {
	args := m.Called(tenantId)
	return args.Error(0)
}

func (m *mockAuth) Logout() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockAuth) GetConfig() auth.Config {
	args := m.Called()
	return args.Get(0).(auth.Config)
}
func TestOperation_SetSecurityHeader(t *testing.T) {
	tests := []struct {
		name           string
		securityMethod string
		forceAuth      bool
		needsAuth      bool
		mockSetup      func(*mockAuth)
		expectedHeader map[string]string
		expectError    bool
	}{
		{
			name:           "ApiKey Auth - Success",
			securityMethod: apiKeyAuthMethod,
			forceAuth:      true,
			needsAuth:      false,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(apiKeyAuthMethod)
				m.On("ApiKey", mock.Anything).Return("test-api-key", nil)
			},
			expectedHeader: map[string]string{"x-api-key": "test-api-key"},
			expectError:    false,
		},
		{
			name:           "XaaS Auth - Success",
			securityMethod: xaasAuthMethod,
			forceAuth:      true,
			needsAuth:      false,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(xaasAuthMethod)
				m.On("XTenantID", mock.Anything).Return("test-tenant", nil)
			},
			expectedHeader: map[string]string{"x-tenant-id": "test-tenant"},
			expectError:    false,
		},
		{
			name:           "Bearer Auth - Success",
			securityMethod: bearerAuthMethod,
			forceAuth:      true,
			needsAuth:      false,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(bearerAuthMethod)
				m.On("AccessToken", mock.Anything).Return("test-token", nil)
			},
			expectedHeader: map[string]string{"Authorization": "Bearer test-token"},
			expectError:    false,
		},
		{
			name:           "OAuth2 - Success",
			securityMethod: oAuth2Method,
			forceAuth:      true,
			needsAuth:      false,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(oAuth2Method)
				m.On("AccessToken", mock.Anything).Return("test-token", nil)
			},
			expectedHeader: map[string]string{"Authorization": "Bearer test-token"},
			expectError:    false,
		},
		{
			name:           "No Auth Required - No Headers Set",
			securityMethod: bearerAuthMethod,
			forceAuth:      false,
			needsAuth:      false,
			mockSetup: func(m *mockAuth) {
			},
			expectedHeader: map[string]string{},
			expectError:    false,
		},
		{
			name:           "Empty ApiKey - Error",
			securityMethod: apiKeyAuthMethod,
			forceAuth:      true,
			needsAuth:      false,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(apiKeyAuthMethod)
				m.On("ApiKey", mock.Anything).Return("", assert.AnError)
			},
			expectedHeader: map[string]string{},
			expectError:    true,
		},
		{
			name:           "Empty XTenantID - Error",
			securityMethod: xaasAuthMethod,
			forceAuth:      true,
			needsAuth:      false,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(xaasAuthMethod)
				m.On("XTenantID", mock.Anything).Return("", assert.AnError)
			},
			expectedHeader: map[string]string{},
			expectError:    true,
		},
		{
			name:           "Empty Access Token - Error",
			securityMethod: bearerAuthMethod,
			forceAuth:      true,
			needsAuth:      false,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(bearerAuthMethod)
				m.On("AccessToken", mock.Anything).Return("", assert.AnError)
			},
			expectedHeader: map[string]string{},
			expectError:    true,
		},
		{
			name:           "Unsupported Auth Method",
			securityMethod: "unknownauth",
			forceAuth:      true,
			needsAuth:      true,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return("unknownauth")
				m.On("AccessToken", mock.Anything).Return("auth-token", nil).Once()
			},
			expectedHeader: map[string]string{},
			expectError:    false,
		},
		{
			name:           "Auth Required but ForceAuth False",
			securityMethod: bearerAuthMethod,
			forceAuth:      false,
			needsAuth:      true,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(bearerAuthMethod)
				m.On("AccessToken", mock.Anything).Return("auth-token", nil)
			},
			expectedHeader: map[string]string{"Authorization": "Bearer auth-token"},
			expectError:    false,
		},
		{
			name:           "Auth Required and ForceAuth True",
			securityMethod: bearerAuthMethod,
			forceAuth:      true,
			needsAuth:      true,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(bearerAuthMethod)
				m.On("AccessToken", mock.Anything).Return("auth-token", nil)
			},
			expectedHeader: map[string]string{"Authorization": "Bearer auth-token"},
			expectError:    false,
		},
		{
			name:           "Auth Not Required but ForceAuth True",
			securityMethod: bearerAuthMethod,
			forceAuth:      true,
			needsAuth:      false,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(bearerAuthMethod)
				m.On("AccessToken", mock.Anything).Return("auth-token", nil)
			},
			expectedHeader: map[string]string{"Authorization": "Bearer auth-token"},
			expectError:    false,
		},
		{
			name:           "Auth Required but AccessToken Error",
			securityMethod: bearerAuthMethod,
			forceAuth:      false,
			needsAuth:      true,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(bearerAuthMethod)
				m.On("AccessToken", mock.Anything).Return("", assert.AnError)
			},
			expectedHeader: map[string]string{},
			expectError:    true,
		},
		{
			name:           "Auth Not Required and ForceAuth False",
			securityMethod: bearerAuthMethod,
			forceAuth:      false,
			needsAuth:      false,
			mockSetup:      func(m *mockAuth) {},
			expectedHeader: map[string]string{},
			expectError:    false,
		},
		{
			name:           "Empty Security Method",
			securityMethod: "",
			forceAuth:      false,
			needsAuth:      false,
			mockSetup:      func(m *mockAuth) {},
			expectedHeader: map[string]string{},
			expectError:    false,
		},
		{
			name:           "Invalid ForceAuth Parameter",
			securityMethod: bearerAuthMethod,
			forceAuth:      false,
			needsAuth:      true,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(bearerAuthMethod)
				m.On("AccessToken", mock.Anything).Return("auth-token", nil)
			},
			expectedHeader: map[string]string{"Authorization": "Bearer auth-token"},
			expectError:    false,
		},
		{
			name:           "API Key Auth with Error",
			securityMethod: apiKeyAuthMethod,
			forceAuth:      true,
			needsAuth:      true,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(apiKeyAuthMethod)
				m.On("ApiKey", mock.Anything).Return("", assert.AnError)
			},
			expectedHeader: map[string]string{},
			expectError:    true,
		},
		{
			name:           "XaaS Auth with Error",
			securityMethod: xaasAuthMethod,
			forceAuth:      true,
			needsAuth:      true,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(xaasAuthMethod)
				m.On("XTenantID", mock.Anything).Return("", assert.AnError)
			},
			expectedHeader: map[string]string{},
			expectError:    true,
		},
		{
			name:           "Empty Access Token when Auth Required",
			securityMethod: bearerAuthMethod,
			forceAuth:      false,
			needsAuth:      true,
			mockSetup: func(m *mockAuth) {
				m.On("CurrentSecurityMethod").Return(bearerAuthMethod)
				m.On("AccessToken", mock.Anything).Return("", assert.AnError)
			},
			expectedHeader: map[string]string{},
			expectError:    true,
		},
		{
			name:           "Bearer Auth with ForceAuth False and NeedsAuth False",
			securityMethod: bearerAuthMethod,
			forceAuth:      false,
			needsAuth:      false,
			mockSetup:      func(m *mockAuth) {},
			expectedHeader: map[string]string{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(mockAuth)
			tt.mockSetup(mockAuth)

			op := &operation{
				operation: &openapi3.Operation{},
			}
			if tt.needsAuth {
				op.operation.Security = &openapi3.SecurityRequirements{
					{tt.securityMethod: []string{}},
				}
			}

			req := &http.Request{
				Header: make(http.Header),
			}

			params := core.Parameters{}
			if tt.forceAuth {
				params[forceAuthParameter] = true
			}

			err := op.setSecurityHeader(context.Background(), params, req, mockAuth)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				for key, value := range tt.expectedHeader {
					assert.Equal(t, value, req.Header.Get(key))
				}
			}

			mockAuth.AssertExpectations(t)
		})
	}
}
