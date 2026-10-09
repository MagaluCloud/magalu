package cmd

import (
	"path/filepath"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry/adapters/debuglog"
)

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

func TestTelemetryDebugLogPath(t *testing.T) {
	configDir := filepath.Join("home", "user", ".config", "mgc")

	testCases := []struct {
		name      string
		debugMode bool
		want      string
	}{
		{"debug mode logs next to telemetry.yaml", true, filepath.Join(configDir, debuglog.FileName)},
		{"no log outside debug mode", false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := telemetryDebugLogPath(tc.debugMode, configDir); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
