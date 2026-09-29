package telemetry

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func formatJsonOutput(s string) string {
	return strings.NewReplacer("\n", "", "\t", "").Replace(s)
}

func TestImpactAndLogLevel(t *testing.T) {
	testCases := []struct {
		outcome  Outcome
		reason   FailureReason
		impact   Impact
		logLevel LogLevel
	}{
		{OutcomeSuccess, "", ImpactLow, LogLevelInfo},
		{OutcomeFailure, FailureValidation, ImpactLow, LogLevelWarn},
		{OutcomeFailure, FailureNotFound, ImpactLow, LogLevelWarn},
		{OutcomeFailure, FailureUserCancelled, ImpactLow, LogLevelWarn},

		{OutcomeFailure, FailureAuthentication, ImpactMedium, LogLevelWarn},
		{OutcomeFailure, FailureAuthorization, ImpactMedium, LogLevelWarn},
		{OutcomeFailure, FailureConflict, ImpactMedium, LogLevelWarn},
		{OutcomeFailure, FailureQuota, ImpactMedium, LogLevelWarn},

		{OutcomeFailure, FailureAPIServer, ImpactHigh, LogLevelWarn},
		{OutcomeFailure, FailureNetwork, ImpactHigh, LogLevelWarn},
		{OutcomeFailure, FailureTimeout, ImpactHigh, LogLevelWarn},
		{OutcomeFailure, FailureUnknown, ImpactHigh, LogLevelWarn},

		{OutcomeFailure, FailureReason("not-in-list"), ImpactHigh, LogLevelWarn},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s/%s", tc.outcome, tc.reason), func(t *testing.T) {
			e := Event{Outcome: tc.outcome, FailureReason: tc.reason}
			if got := e.Impact(); got != tc.impact {
				t.Errorf("impact = %s, want %s", got, tc.impact)
			}
			if got := e.LogLevel(); got != tc.logLevel {
				t.Errorf("logLevel = %s, want %s", got, tc.logLevel)
			}
		})
	}
}

func TestMarshalJSONGolden(t *testing.T) {
	ttfv := 4231 * time.Millisecond

	testCases := []struct {
		name  string
		event Event
		want  string
	}{
		{
			name: "success",
			event: Event{
				Timestamp:        time.Date(2026, 9, 22, 11, 30, 45, 0, time.FixedZone("BRT", -3*60*60)),
				Actor:            &Actor{TenantID: "6a7b8c9d"},
				ExecutionContext: ExecutionContextInteractive,
				Resource:         &Resource{Type: "virtual_machine_instances"},
				Action:           "virtualmachine.instances.list",
				OptionsSet:       []string{"region", "machine-type", "region"},
				Outcome:          OutcomeSuccess,
				FailureReason:    FailureAPIServer, // ignorado no sucesso
				Duration:         842 * time.Millisecond,
				ClientVersion:    "v1.4.2",
				OS:               "linux",
				InstallMethod:    "homebrew",
				LastRequestID:    "550e8400-e29b-41d4-a716-446655440000",
				InstallationID:   "should-not-be-serialized",
			},
			want: `{
				"timestamp":"2026-09-22T14:30:45.000Z",
				"actor":{"tenant_id":"6a7b8c9d"},
				"executionContext":"interactive",
				"resource":{"type":"virtual_machine_instances"},
				"action":"virtualmachine.instances.list",
				"optionsSet":"machine-type,region",
				"outcome":"success",
				"durationMs":842,
				"logLevel":"INFO",
				"impact":"LOW",
				"clientVersion":"v1.4.2",
				"os":"linux",
				"installMethod":"homebrew",
				"lastRequestId":"550e8400-e29b-41d4-a716-446655440000"
				}`,
		},
		{
			name: "failure",
			event: Event{
				Timestamp:            time.Date(2026, 9, 22, 14, 30, 45, 123456789, time.UTC),
				ExecutionContext:     ExecutionContextCI,
				ExecutionEnvironment: "github_actions",
				Resource:             &Resource{Type: "virtual_machine_instances", ID: "abc"},
				Action:               "virtualmachine.instances.get",
				Outcome:              OutcomeFailure,
				FailureReason:        FailureQuota,
				Duration:             -5 * time.Millisecond,
				ClientVersion:        "v1.4.2",
				OS:                   "darwin",
				TimeToFirstValue:     &ttfv,
			},
			want: `{
				"timestamp":"2026-09-22T14:30:45.123Z",
				"executionContext":"ci",
				"executionEnvironment":"github_actions",
				"resource":{"type":"virtual_machine_instances","id":"abc"},
				"action":"virtualmachine.instances.get",
				"outcome":"failure",
				"failure_reason":"quota",
				"durationMs":0,
				"logLevel":"WARN",
				"impact":"MEDIUM",
				"clientVersion":"v1.4.2",
				"os":"darwin",
				"installMethod":"manual",
				"timeToFirstValueMs":4231
				}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.event)
			if err != nil {
				t.Fatal(err)
			}
			if want := formatJsonOutput(tc.want); string(got) != want {
				t.Errorf("got  %s\nwant %s", got, want)
			}
		})
	}
}

func TestMarshalJSONPayloadRules(t *testing.T) {
	negativeTTFV := -2 * time.Second

	testCases := []struct {
		name        string
		event       Event
		contains    []string
		notContains []string
	}{
		{
			name: "conditional fields are omitted",
			event: Event{
				Timestamp:        time.Now(),
				Actor:            &Actor{},
				Resource:         &Resource{},
				ExecutionContext: ExecutionContextInteractive,
				Action:           "unknown",
				Outcome:          OutcomeSuccess,
			},
			notContains: []string{`"failure_reason"`, `"actor"`, `"resource"`, `"optionsSet"`, `"timeToFirstValueMs"`,
				`"lastRequestId"`, `"executionEnvironment"`, `"distinct_id"`, `"InstallationID"`},
		},
		{
			name:        "invalid failure reason becomes unknown",
			event:       Event{Outcome: OutcomeFailure, FailureReason: "boom: stack trace here"},
			contains:    []string{`"failure_reason":"unknown"`},
			notContains: []string{"stack trace"},
		},
		{
			name: "payload is anonymized",
			event: Event{
				Timestamp:        time.Now(),
				ExecutionContext: ExecutionContextInteractive,
				Action:           "virtualmachine.instances.create",
				OptionsSet:       []string{"region", "api-key"},
				Outcome:          OutcomeFailure,
				FailureReason:    FailureValidation,
				ClientVersion:    "v1.4.2",
				OS:               "linux",
				InstallationID:   "b3f1c2de-0000-4000-8000-000000000000",
			},
			contains: []string{`"optionsSet":"api-key,region"`},
			notContains: []string{"abc123", "br-ne1", "Bearer eyJhbGci", "user@example.com", "/home/user/.config/mgc",
				`"ip"`, `"source"`, `"message"`, `"error"`, `"stack"`},
		},
		{
			name:        "negative time to first value is omitted",
			event:       Event{Outcome: OutcomeSuccess, TimeToFirstValue: &negativeTTFV},
			notContains: []string{"timeToFirstValueMs"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.event)
			if err != nil {
				t.Fatal(err)
			}
			for _, s := range tc.contains {
				if !strings.Contains(string(got), s) {
					t.Errorf("payload must contain %s: %s", s, got)
				}
			}
			for _, s := range tc.notContains {
				if strings.Contains(string(got), s) {
					t.Errorf("payload must not contain %s: %s", s, got)
				}
			}
		})
	}
}
