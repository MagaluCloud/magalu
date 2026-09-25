package telemetry

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func formatJsonOutput(s string) string {
	return strings.NewReplacer("\n", "", "\t", "").Replace(s)
}

func TestImpactAndLogLevel(t *testing.T) {
	cases := []struct {
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

	for _, c := range cases {
		e := Event{Outcome: c.outcome, FailureReason: c.reason}
		if got := e.Impact(); got != c.impact {
			t.Errorf("%s/%s: impact = %s, want %s", c.outcome, c.reason, got, c.impact)
		}
		if got := e.LogLevel(); got != c.logLevel {
			t.Errorf("%s/%s: logLevel = %s, want %s", c.outcome, c.reason, got, c.logLevel)
		}
	}
}

func TestMarshalJSONSuccessGolden(t *testing.T) {
	e := Event{
		Timestamp:        time.Date(2026, 9, 22, 11, 30, 45, 0, time.FixedZone("BRT", -3*60*60)),
		Actor:            &Actor{TenantID: "6a7b8c9d"},
		ExecutionContext: ExecutionContextInteractive,
		Resource:         &Resource{Type: "virtual_machine_instances"},
		Action:           "virtualmachine.instances.list",
		OptionsSet:       []string{"region", "machine-type", "region"},
		Outcome:          OutcomeSuccess,
		FailureReason:    FailureAPIServer, // ignored on success
		Duration:         842 * time.Millisecond,
		ClientVersion:    "v1.4.2",
		OS:               "linux",
		InstallMethod:    "homebrew",
		LastRequestID:    "550e8400-e29b-41d4-a716-446655440000",
		InstallationID:   "should-not-be-serialized",
	}

	got, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	want := `{
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
		}`

	formatedWant := formatJsonOutput(want)
	if string(got) != formatedWant {
		t.Errorf("got  %s\nwant %s", got, formatedWant)
	}
}

func TestMarshalJSONFailureGolden(t *testing.T) {
	ttfv := 4231 * time.Millisecond
	e := Event{
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
	}

	got, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	want := `{
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
		"installMethod":"unknown",
		"timeToFirstValueMs":4231
		}`

	formatedWant := formatJsonOutput(want)
	if string(got) != formatedWant {
		t.Errorf("got  %s\nwant %s", got, formatedWant)
	}
}

func TestMarshalJSONOmitsConditionalFields(t *testing.T) {
	e := Event{
		Timestamp:        time.Now(),
		Actor:            &Actor{},
		Resource:         &Resource{},
		ExecutionContext: ExecutionContextInteractive,
		Action:           "unknown",
		Outcome:          OutcomeSuccess,
	}

	got, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any
	if err := json.Unmarshal(got, &m); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"failure_reason", "actor", "resource", "optionsSet", "timeToFirstValueMs", "lastRequestId", "executionEnvironment", "distinct_id", "InstallationID"} {
		if _, ok := m[key]; ok {
			t.Errorf("key %q should be omitted, got %s", key, got)
		}
	}
}

func TestMarshalJSONInvalidFailureReasonBecomesUnknown(t *testing.T) {
	e := Event{Outcome: OutcomeFailure, FailureReason: "boom: stack trace here"}
	got, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), `"failure_reason":"unknown"`) || strings.Contains(string(got), "stack trace") {
		t.Errorf("unexpected payload %s", got)
	}
}

func TestAnonymizedPayload(t *testing.T) {
	// Simulates: mgc ... --api-key=abc123 --region=br-ne1, failing with a detailed API message.
	// Only flag names reach the event; values and messages never do.
	secrets := []string{"abc123", "br-ne1", "Bearer eyJhbGci", "user@example.com", "/home/user/.config/mgc"}

	e := Event{
		Timestamp:        time.Now(),
		ExecutionContext: ExecutionContextInteractive,
		Action:           "virtualmachine.instances.create",
		OptionsSet:       []string{"region", "api-key"},
		Outcome:          OutcomeFailure,
		FailureReason:    FailureValidation,
		ClientVersion:    "v1.4.2",
		OS:               "linux",
		InstallationID:   "b3f1c2de-0000-4000-8000-000000000000",
	}

	got, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range secrets {
		if strings.Contains(string(got), s) {
			t.Errorf("payload leaks %q: %s", s, got)
		}
	}
	if !strings.Contains(string(got), `"optionsSet":"api-key,region"`) {
		t.Errorf("optionsSet should carry only sorted flag names: %s", got)
	}
	for _, key := range []string{"ip", "source", "message", "error", "stack"} {
		if strings.Contains(string(got), `"`+key+`"`) {
			t.Errorf("payload must not contain key %q: %s", key, got)
		}
	}
}

func TestMarshalJSONOmitsNegativeTimeToFirstValue(t *testing.T) {
	ttfv := -2 * time.Second
	e := Event{Outcome: OutcomeSuccess, TimeToFirstValue: &ttfv}
	got, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(got), "timeToFirstValueMs") {
		t.Errorf("negative TTFV must be omitted: %s", got)
	}
}
