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
		{"development build without key uses the noop exporter", "", nil, DefaultEndpoint, true},
		{"blank key is treated as missing", "  ", nil, DefaultEndpoint, true},
		{"release build uses posthog on the default endpoint", "phc_test", nil, DefaultEndpoint, false},
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

			switch exporter := cfg.Exporter().(type) {
			case NoopExporter:
				if !tc.wantNoop {
					t.Errorf("want the posthog exporter, got noop")
				}
			case *PostHogExporter:
				if tc.wantNoop {
					t.Errorf("want the noop exporter, got posthog")
				}
				if exporter.endpoint != tc.wantEndpoint {
					t.Errorf("exporter endpoint = %q, want %q", exporter.endpoint, tc.wantEndpoint)
				}
			default:
				t.Errorf("unexpected exporter %T", exporter)
			}
		})
	}
}

func TestBuildConfigDefaultsToNoopWithoutLdflags(t *testing.T) {
	if _, ok := LoadBuildConfig(envFrom(nil)).Exporter().(NoopExporter); !ok {
		t.Errorf("a build without -ldflags must not send anything")
	}
}
