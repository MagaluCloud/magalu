package project

import "testing"

func TestHostForEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		expected string
	}{
		{name: "empty falls back to prod", env: "", expected: prodHost},
		{name: "default falls back to prod", env: "default", expected: prodHost},
		{name: "prod", env: "prod", expected: prodHost},
		{name: "pre-prod", env: "pre-prod", expected: preProdHost},
		{name: "raw pre-prod host", env: preProdHost, expected: preProdHost},
		{name: "unknown falls back to prod", env: "staging", expected: prodHost},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hostForEnv(tt.env); got != tt.expected {
				t.Errorf("hostForEnv(%q) = %q, expected %q", tt.env, got, tt.expected)
			}
		})
	}
}

func TestBuildProjectsURL(t *testing.T) {
	tests := []struct {
		name      string
		serverUrl string
		env       string
		expected  string
	}{
		{
			name:     "prod default",
			env:      "",
			expected: "https://api.magalu.cloud/iam/api/v1/projects",
		},
		{
			name:     "pre-prod",
			env:      "pre-prod",
			expected: "https://api.pre-prod.jaxyendy.com:8443/iam/api/v1/projects",
		},
		{
			name:      "serverUrl replaces the whole base, including the /iam prefix",
			serverUrl: "http://localhost:8080",
			env:       "prod",
			expected:  "http://localhost:8080/api/v1/projects",
		},
		{
			name:      "serverUrl trailing slash",
			serverUrl: "http://localhost:8080/",
			expected:  "http://localhost:8080/api/v1/projects",
		},
		{
			name:      "serverUrl wins over env",
			serverUrl: "https://example.com/iam",
			env:       "pre-prod",
			expected:  "https://example.com/iam/api/v1/projects",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildProjectsURL(tt.serverUrl, tt.env)
			if got != tt.expected {
				t.Errorf("buildProjectsURL(%q, %q) = %q, expected %q", tt.serverUrl, tt.env, got, tt.expected)
			}
		})
	}
}
