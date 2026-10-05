package cmd

import (
	"testing"

	"github.com/spf13/cobra"
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

func TestLikelyRecorded(t *testing.T) {
	loggedIn := credentials{loggedInBefore: true}
	vmCreate := []string{"virtual-machine", "instances", "create", "--name", "vm1"}

	testCases := []struct {
		name  string
		args  []string
		creds credentials
		env   map[string]string
		want  bool
	}{
		{"tracked command with session", vmCreate, loggedIn, nil, true},
		{"tracked command without credentials", vmCreate, credentials{}, nil, false},
		{"api key flag", append(vmCreate, "--api-key", "key"), credentials{}, nil, true},
		{"api key flag with value", append(vmCreate, "--api-key=key"), credentials{}, nil, true},
		{"api key env", vmCreate, credentials{}, map[string]string{apiKeyEnvVar: "key"}, true},
		{"login without credentials", []string{"auth", "login"}, credentials{}, nil, true},
		{"help flag", append(vmCreate, "--help"), loggedIn, nil, false},
		{"help shorthand", append(vmCreate, "-h"), loggedIn, nil, false},
		{"version", []string{"--version"}, loggedIn, nil, false},
		{"group without subcommand", []string{"virtual-machine"}, loggedIn, nil, false},
		{"untracked command", []string{telemetryCmdName, "disable"}, loggedIn, nil, false},
		{"unknown command with session", []string{"foo", "bar"}, loggedIn, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			root := newFakeCommandTree()
			auth := &cobra.Command{Use: "auth"}
			auth.AddCommand(&cobra.Command{Use: "login", RunE: func(*cobra.Command, []string) error { return nil }})
			root.AddCommand(auth)
			getenv := func(k string) string { return tc.env[k] }

			if got := likelyRecorded(root, tc.args, tc.creds, getenv); got != tc.want {
				t.Errorf("likelyRecorded = %v, want %v", got, tc.want)
			}
		})
	}
}
