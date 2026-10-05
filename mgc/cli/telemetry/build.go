package telemetry

import "strings"

// inserido via -ldflags
// preenchido no release com -X github.com/MagaluCloud/magalu/mgc/cli/telemetry.posthogAPIKey (release.yaml, internal.yaml)
var posthogAPIKey string

const (
	DefaultEndpoint = "https://eu.i.posthog.com/i/v1/logs"
	EnvEndpoint     = "MGC_TELEMETRY_ENDPOINT"
)

type BuildConfig struct {
	APIKey   string
	Endpoint string
}

func LoadBuildConfig(getenv func(string) string) BuildConfig {
	return newBuildConfig(posthogAPIKey, getenv)
}

func newBuildConfig(apiKey string, getenv func(string) string) BuildConfig {
	cfg := BuildConfig{APIKey: strings.TrimSpace(apiKey), Endpoint: DefaultEndpoint}
	if endpoint := strings.TrimSpace(getenv(EnvEndpoint)); endpoint != "" {
		cfg.Endpoint = endpoint
	}
	return cfg
}

func (c BuildConfig) Exporter() Exporter {
	if c.APIKey == "" {
		return NoopExporter{}
	}
	return NewPostHogExporter(c.Endpoint, c.APIKey)
}
