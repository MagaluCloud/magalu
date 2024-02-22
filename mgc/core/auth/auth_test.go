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

	"magalu.cloud/core"
	"magalu.cloud/core/config"
	"magalu.cloud/core/profile_manager"
	"magalu.cloud/core/utils"
	"magalu.cloud/testing/fs_test_helper"
)

var dummyConfigResult *ConfigResult = &ConfigResult{
	AccessToken:  "access-token",
	RefreshToken: "refresh-token",
	CurrentEnv:   "test",
}

var dummyConfigResultYaml = []byte(`---
access_token: "access-token"
refresh_token: "refresh-token"
current_environment: "test"
`)

var dummyConfigResultYamlRealToken = []byte(`---
access_token: "eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiYWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzA0ODI0MzUzLCJpYXQiOjE3MDQ4MjQzNTN9.Eu1WKIEja4OQc87FH2ku-34Gir2P6RjGATEu-EQPAC8"
refresh_token: "eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiYWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzA0ODI0MzUzLCJpYXQiOjE3MDQ4MjQzNTN9.Eu1WKIEja4OQc87FH2ku-34Gir2P6RjGATEu-EQPAC8"
current_environment: "test"
`)

var dummyConfigMap map[string]Config = map[string]Config{
	"temp": {
		ClientId:       "client-id",
		RedirectUri:    "redirect-uri",
		LoginUrl:       "login-url",
		TokenUrl:       "token-url",
		ValidationUrl:  "validation-url",
		RefreshUrl:     "refresh-url",
		TenantsListUrl: "tenant-list-url",
	},
}

type mockTransport struct {
	statusCode        int
	responseBody      io.ReadCloser
	shouldReturnError bool
}

func (o mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	if o.shouldReturnError {
		return nil, fmt.Errorf("test error")
	}
	return &http.Response{StatusCode: o.statusCode, Body: o.responseBody, Request: &http.Request{Header: http.Header{"X-Request-Id": []string{""}}, Response: &http.Response{}}}, nil
}

type testCaseAuth struct {
	name           string
	transport      mockTransport
	expectedFs     []fs_test_helper.TestFsEntry
	providedFs     []fs_test_helper.TestFsEntry
	run            func(a *Auth) error
	envAccessToken string
}

func setTokens(name string, expected_error error, access_token string, refresh_token string, transport mockTransport, provided []fs_test_helper.TestFsEntry, expected []fs_test_helper.TestFsEntry) testCaseAuth {
	provided = fs_test_helper.AutoMkdirAll(provided)
	expected = fs_test_helper.AutoMkdirAll(expected)
	return testCaseAuth{
		name:       fmt.Sprintf("Auth.Set(%q)", name),
		transport:  transport,
		providedFs: provided,
		expectedFs: expected,
		run: func(auth *Auth) error {
			dummyLoginResult := &LoginResult{
				AccessToken:  access_token,
				RefreshToken: refresh_token,
			}

			err := auth.SetTokens(dummyLoginResult)

			if err != nil {
				return expected_error
			}

			return nil
		},
	}
}

func setAccessKey(name string, expected_error error, acess_key_id string, secret_access_key string, transport mockTransport, provided []fs_test_helper.TestFsEntry, expected []fs_test_helper.TestFsEntry) testCaseAuth {
	provided = fs_test_helper.AutoMkdirAll(provided)
	expected = fs_test_helper.AutoMkdirAll(expected)
	return testCaseAuth{
		name:       fmt.Sprintf("Auth.SetAccessKey(%q)", name),
		providedFs: provided,
		expectedFs: expected,
		transport:  transport,
		run: func(auth *Auth) error {
			err := auth.SetAccessKey(acess_key_id, secret_access_key)
			if err != nil {
				return expected_error
			}
			return nil
		},
	}
}

func requestAuthTokenWithAuthorizationCode(name string, transport mockTransport, verifier *codeVerifier, expectedErr bool, provided []fs_test_helper.TestFsEntry, expected []fs_test_helper.TestFsEntry) testCaseAuth {
	provided = fs_test_helper.AutoMkdirAll(provided)
	expected = fs_test_helper.AutoMkdirAll(expected)
	return testCaseAuth{
		name:       fmt.Sprintf("Auth.RequestAuth(%q)", name),
		transport:  transport,
		providedFs: provided,
		expectedFs: expected,
		run: func(auth *Auth) error {
			auth.codeVerifier = verifier
			err := auth.RequestAuthTokenWithAuthorizationCode(context.Background(), "")
			hasErr := err != nil

			if hasErr != expectedErr {
				return fmt.Errorf("expected error == %v", expectedErr)
			}

			return nil
		},
	}
}

func doRefreshAccessToken(name string, transport mockTransport, expectedErr bool, expectedResult string, provided []fs_test_helper.TestFsEntry, expected []fs_test_helper.TestFsEntry) testCaseAuth {
	provided = fs_test_helper.AutoMkdirAll(provided)
	expected = fs_test_helper.AutoMkdirAll(expected)
	return testCaseAuth{
		name:       fmt.Sprintf("Auth.DoRefreshAccessToken(%q)", name),
		transport:  transport,
		providedFs: provided,
		expectedFs: expected,
		run: func(auth *Auth) error {
			tk, err := auth.doRefreshAccessToken(context.Background())
			hasErr := err != nil

			if hasErr != expectedErr {
				return fmt.Errorf("expected err == %v", expectedErr)
			}

			if tk != expectedResult {
				return fmt.Errorf("expected tk == %v, found: %v", expectedResult, tk)
			}
			return nil
		},
	}
}

func validateAccessToken(name string, transport mockTransport, expectedErr bool, provided []fs_test_helper.TestFsEntry, expected []fs_test_helper.TestFsEntry) testCaseAuth {
	provided = fs_test_helper.AutoMkdirAll(provided)
	expected = fs_test_helper.AutoMkdirAll(expected)
	return testCaseAuth{
		name:       fmt.Sprintf("Auth.ValidateAccess(%q)", name),
		providedFs: provided,
		expectedFs: expected,
		transport:  transport,
		run: func(auth *Auth) error {
			err := auth.ValidateAccessToken(context.Background())

			hasErr := err != nil
			if hasErr != expectedErr {
				return fmt.Errorf("expected error == %v", expectedErr)
			}

			return nil
		},
	}
}

func selectTenant(name string, transport mockTransport, expectedResult *TokenExchangeResult, expectedErr bool, provided []fs_test_helper.TestFsEntry, expected []fs_test_helper.TestFsEntry) testCaseAuth {
	provided = fs_test_helper.AutoMkdirAll(provided)
	expected = fs_test_helper.AutoMkdirAll(expected)
	return testCaseAuth{
		name:       fmt.Sprintf("Auth.SelectTenant(%q)", name),
		transport:  transport,
		providedFs: provided,
		expectedFs: expected,
		run: func(auth *Auth) error {
			tnt, err := auth.SelectTenant(context.Background(), `qwe123`, "")
			hasErr := err != nil

			if hasErr != expectedErr {
				return fmt.Errorf("expected 'hasErr' == %v, got 'err' == %w", expectedErr, err)
			}
			if !reflect.DeepEqual(tnt, expectedResult) {
				return fmt.Errorf("expected tnt == %v, found: %v", expectedResult, tnt)
			}
			return nil
		},
	}
}

func listTenants(name string, transport mockTransport, expectedTenants []*Tenant, expectedErr bool, provided []fs_test_helper.TestFsEntry, expected []fs_test_helper.TestFsEntry) testCaseAuth {
	provided = fs_test_helper.AutoMkdirAll(provided)
	expected = fs_test_helper.AutoMkdirAll(expected)
	return testCaseAuth{
		name:       fmt.Sprintf("Auth.ListTenants(%q)", name),
		transport:  transport,
		providedFs: provided,
		expectedFs: expected,
		run: func(auth *Auth) error {
			tLst, err := auth.ListTenants(context.Background(), &http.Client{Transport: transport})
			hasErr := err != nil

			if hasErr != expectedErr {
				return fmt.Errorf("expected err == %v", expectedErr)
			}
			if !reflect.DeepEqual(tLst, expectedTenants) {
				return fmt.Errorf("expected tLst == %v, found: %v", expectedTenants, tLst)
			}
			return nil
		},
	}
}

func newAuth(name string, envAccessToken string, expectedConfig *ConfigResult, provided []fs_test_helper.TestFsEntry, expected []fs_test_helper.TestFsEntry) testCaseAuth {
	provided = fs_test_helper.AutoMkdirAll(provided)
	expected = fs_test_helper.AutoMkdirAll(expected)
	return testCaseAuth{
		name:           fmt.Sprintf("Auth.NewAuth(%q)", name),
		providedFs:     provided,
		expectedFs:     expected,
		envAccessToken: envAccessToken,

		run: func(auth *Auth) error {
			if auth.accessToken != expectedConfig.AccessToken {
				return fmt.Errorf("expected auth.accessToken == '', found: %v", auth.accessToken)
			}

			if auth.refreshToken != expectedConfig.RefreshToken {
				return fmt.Errorf("expected auth.refreshToken == '', found: %v", auth.refreshToken)
			}

			return nil
		},
	}
}

func TestAuthManager(t *testing.T) {
	tests := []testCaseAuth{

		setTokens("Valid token", nil, "access-token", "refresh-token", mockTransport{},
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`access_key_id: ""
access_token: access-token
current_environment: ""
refresh_token: refresh-token
secret_access_key: ""
`),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		setTokens("Valid token without auth file", nil, "access-token", "refresh-token", mockTransport{},
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/",
					Mode: utils.DIR_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`access_key_id: ""
access_token: access-token
current_environment: ""
refresh_token: refresh-token
secret_access_key: ""
`),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		setAccessKey("Valid keys", nil, "MyAccessKeyIdTest", "MySecretAccessKeyTeste", mockTransport{},
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`access_key_id: MyAccessKeyIdTest
access_token: ""
current_environment: ""
refresh_token: ""
secret_access_key: MySecretAccessKeyTeste
`),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		setAccessKey("Valid keys without auth file", nil, "MyAccessKeyIdTest", "MySecretAccessKeyTeste", mockTransport{},
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/",
					Mode: utils.DIR_PERMISSION,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`access_key_id: MyAccessKeyIdTest
access_token: ""
current_environment: ""
refresh_token: ""
secret_access_key: MySecretAccessKeyTeste
`),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		requestAuthTokenWithAuthorizationCode("Code verifier == nil", mockTransport{}, nil, true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(``),
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(``),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		requestAuthTokenWithAuthorizationCode("Bad request",
			mockTransport{
				statusCode:   http.StatusBadRequest,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte{})),
			},
			&codeVerifier{},
			false,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		requestAuthTokenWithAuthorizationCode("Valid login result",
			mockTransport{
				statusCode: http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{
													"access_token": "ac-token",
													"refresh_token": "rf-token"
												}`))),
			},
			&codeVerifier{},
			false,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`access_key_id: ""
access_token: ac-token
current_environment: ""
refresh_token: rf-token
secret_access_key: ""
`),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		requestAuthTokenWithAuthorizationCode("Valid login result without file",
			mockTransport{
				statusCode: http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{
													"access_token": "ac-token",
													"refresh_token": "rf-token"
												}`))),
			},
			&codeVerifier{},
			false,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/",
					Mode: utils.DIR_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`access_key_id: ""
access_token: ac-token
current_environment: ""
refresh_token: rf-token
secret_access_key: ""
`),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		requestAuthTokenWithAuthorizationCode("Invalid login result",
			mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{`))),
			},
			&codeVerifier{},
			true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(``),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		requestAuthTokenWithAuthorizationCode("Request with error",
			mockTransport{
				shouldReturnError: true,
			},
			&codeVerifier{},
			true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(``),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		validateAccessToken("Request ended with error",
			mockTransport{
				shouldReturnError: true,
			},
			true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(``),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		validateAccessToken("Invalid validation result",
			mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{`))),
			},
			true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(``),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		validateAccessToken("Bad request",
			mockTransport{
				statusCode:   http.StatusBadRequest,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte{})),
			},
			true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(``),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		validateAccessToken("Active validation result",
			mockTransport{
				statusCode: http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{
								"active": true
							}`))),
			},
			false,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(``),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		doRefreshAccessToken("Valid response json",
			mockTransport{
				shouldReturnError: true,
			}, true, "access-token",
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			},
		),

		doRefreshAccessToken("Valid response json",
			mockTransport{
				statusCode: http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{
											"access_token": "ac-token",
											"refresh_token": "rf-token"
										}`))),
			}, false, "ac-token",
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`access_key_id: ""
access_token: ac-token
current_environment: ""
refresh_token: rf-token
secret_access_key: ""
`),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			},
		),
		doRefreshAccessToken("Bad request",
			mockTransport{
				statusCode:   http.StatusBadRequest,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte{})),
			}, true, "",
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			},
		),
		doRefreshAccessToken("Invalid response json",
			mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{`))),
			}, true, "",
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			},
		),
		selectTenant("Invalid tenant result",
			mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{`))),
			},
			nil,
			true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		selectTenant("Valid tenant result",
			mockTransport{
				statusCode: http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{
									"id": "qwe123",
									"access_token": "abc",
									"created_at": 0,
									"refresh_token": "def",
									"scope": "test"
								}`))),
			},
			&TokenExchangeResult{
				TenantID:     "qwe123",
				CreatedAt:    core.Time(time.Unix(int64(0), 0)),
				AccessToken:  "abc",
				RefreshToken: "def",
				Scope:        []string{"test"},
			},
			false,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYamlRealToken,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`access_key_id: ""
access_token: abc
current_environment: ""
refresh_token: def
secret_access_key: ""
`),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		listTenants("empty tenant list",
			mockTransport{
				statusCode:   http.StatusOK,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte(`[]`))),
			}, []*Tenant{}, false,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		newAuth("empty auth file", "",
			&ConfigResult{},
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(``),
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		newAuth("non empty auth file", "",
			dummyConfigResult,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		newAuth("Not-empty auth file with env var", "env-access-token",
			&ConfigResult{
				AccessToken:  "env-access-token",
				RefreshToken: dummyConfigResult.RefreshToken,
			},
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),

		listTenants("non empty tenant list",
			mockTransport{
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
			}, []*Tenant{
				{UUID: "1", Name: "jon doe", Email: "jon.doe@profusion.mobi", IsManaged: false, IsDelegated: false},
				{UUID: "2", Name: "jon smith", Email: "jon.smith@profusion.mobi", IsManaged: false, IsDelegated: false},
			}, false,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		listTenants("request ended with err", mockTransport{
			shouldReturnError: true,
		}, nil, true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		listTenants("bad request", mockTransport{
			statusCode:   http.StatusBadRequest,
			responseBody: io.NopCloser(bytes.NewBuffer([]byte{})),
		}, nil, true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),

		listTenants("invalid tenant list", mockTransport{
			statusCode:   http.StatusOK,
			responseBody: io.NopCloser(bytes.NewBuffer([]byte(`{`))),
		}, nil, true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		selectTenant("Bad request",
			mockTransport{
				statusCode:   http.StatusBadRequest,
				responseBody: io.NopCloser(bytes.NewBuffer([]byte{})),
			},
			nil,
			true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
		selectTenant("Request ended with error",

			mockTransport{
				shouldReturnError: true,
			},
			nil,
			true,
			[]fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
			}, []fs_test_helper.TestFsEntry{
				{
					Path: "/default/auth.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: dummyConfigResultYaml,
				},
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`env: temp
`),
				},
			}),
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, fs := profile_manager.NewInMemoryProfileManager()
			config := config.New(m)
			err := config.Set("env", "temp")
			if err != nil {
				t.Errorf("could not set env")
			}
			fs_err := fs_test_helper.PrepareFs(fs, tc.providedFs)
			if fs_err != nil {
				t.Errorf("could not prepare provided FS: %s", fs_err.Error())
			}

			// TODO: it's required to NewAuth test. Check how to handle it better
			t.Setenv("MGC_SDK_ACCESS_TOKEN", tc.envAccessToken)

			auth := New(dummyConfigMap, &http.Client{Transport: tc.transport}, m, config)

			run_error := tc.run(auth)

			if run_error != nil {
				t.Errorf("expected err == nil, found: %v", run_error)
			}

			fs_err = fs_test_helper.CheckFs(fs, tc.expectedFs)

			if fs_err != nil {
				t.Errorf("unexpected FS state: %s", fs_err.Error())
			}

		})
	}
}
