package cmd

import "testing"

func TestIsAuthenticated(t *testing.T) {
	testCases := []struct {
		name  string
		creds credentials
		env   map[string]string
		want  bool
	}{
		{"no credentials", credentials{}, nil, false},
		{"session before the command", credentials{loggedInBefore: true}, nil, true},
		{"session after the command", credentials{tenantIDAfter: "tenant-1"}, nil, true},
		{"api key flag", credentials{apiKeyFlag: "key"}, nil, true},
		{"saved api key", credentials{savedAPIKey: "key"}, nil, true},
		{"api key env", credentials{}, map[string]string{apiKeyEnvVar: "key"}, true},
		{"saved key pair", credentials{savedKeyID: "id", savedKeySecret: "secret"}, nil, true},
		{"saved key id without secret", credentials{savedKeyID: "id"}, nil, false},
		{"key pair env", credentials{}, map[string]string{objKeyIDEnvVar: "id", objKeySecretEnvVar: "secret"}, true},
		{"key id env without secret", credentials{}, map[string]string{objKeyIDEnvVar: "id"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			getenv := func(k string) string { return tc.env[k] }
			if got := isAuthenticated(tc.creds, getenv); got != tc.want {
				t.Errorf("isAuthenticated = %v, want %v", got, tc.want)
			}
		})
	}
}
