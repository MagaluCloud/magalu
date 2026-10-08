package telemetry

import "testing"

func TestBuildConfig(t *testing.T) {
	testCases := []struct {
		name         string
		apiKey       string
		env          map[string]string
		wantEndpoint string
		wantNoop     bool
	}{
		{
			"development build without key uses the noop exporter",
			"",
			nil,
			DefaultEndpoint,
			true,
		},
		{
			"blank key is treated as missing",
			"  ",
			nil,
			DefaultEndpoint,
			true,
		},
		{
			"release build uses posthog on the default endpoint",
			"phc_test",
			nil,
			DefaultEndpoint,
			false,
		},
		{
			"endpoint override replaces the whole url",
			"phc_test",
			map[string]string{EnvEndpoint: "http://localhost:8080/i/v1/logs"},
			"http://localhost:8080/i/v1/logs",
			false,
		},
		{
			"endpoint override without key still sends nothing",
			"",
			map[string]string{EnvEndpoint: "http://localhost:8080/i/v1/logs"},
			"http://localhost:8080/i/v1/logs",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := newBuildConfig(tc.apiKey, envFrom(tc.env))

			if cfg.Endpoint != tc.wantEndpoint {
				t.Errorf("endpoint = %q, want %q", cfg.Endpoint, tc.wantEndpoint)
			}

			if cfg.Enabled() == tc.wantNoop {
				t.Errorf("Enabled() = %v, want %v", cfg.Enabled(), !tc.wantNoop)
			}

			exporters := cfg.Exporters(nil)
			if tc.wantNoop {
				if len(exporters) != 0 {
					t.Errorf("want no exporters, got %d", len(exporters))
				}
				return
			}
			if len(exporters) != 1 {
				t.Fatalf("want a single exporter, got %d", len(exporters))
			}
			posthog, ok := exporters[0].(*PostHogExporter)
			if !ok {
				t.Fatalf("unexpected exporter %T", exporters[0])
			}
			if posthog.endpoint != tc.wantEndpoint {
				t.Errorf("exporter endpoint = %q, want %q", posthog.endpoint, tc.wantEndpoint)
			}
		})
	}
}

func TestBuildConfigDefaultsToNoopWithoutLdflags(t *testing.T) {
	if cfg := LoadBuildConfig(envFrom(nil)); cfg.Enabled() || len(cfg.Exporters(nil)) != 0 {
		t.Errorf("a build without -ldflags must not send anything")
	}
}
