package telemetry

import (
	"errors"
	"testing"
	"time"
)

func TestRecord(t *testing.T) {
	readOnly := errors.New("read-only home")

	testCases := []struct {
		name       string
		setup      testSetup
		executions []execution
		check      func(t *testing.T, r runResult)
	}{
		{
			name:       "interactive login shows the notice and is collected with a single save",
			setup:      testSetup{terminal: true},
			executions: []execution{{path: authLogin, stderr: expectedNotice, events: 1}},
			check: func(t *testing.T, r runResult) {
				if r.saves != 1 || r.state.InstallationID == "" || !r.state.NoticeShown || r.state.CredentialsSetAt == nil {
					t.Errorf("saves=%d state=%+v, want a single save with id, notice and credentials", r.saves, r.state)
				}
			},
		},
		{
			name:       "already logged in user is not collected in the execution that shows the notice",
			setup:      testSetup{terminal: true},
			executions: []execution{{path: vmList, stderr: expectedNoticeNotCollected}},
			check: func(t *testing.T, r runResult) {
				if r.state.InstallationID != "" || r.state.FirstValueRecorded {
					t.Errorf("only notice_shown may be written, got %+v", r.state)
				}
			},
		},
		{
			name:  "already logged in user is collected from the next execution on",
			setup: testSetup{terminal: true},
			executions: []execution{
				{path: vmList, stderr: expectedNoticeNotCollected},
				{path: vmList, events: 1},
			},
		},
		{
			name:       "failed login is treated as any other command",
			setup:      testSetup{terminal: true},
			executions: []execution{{path: authLogin, cmdErr: errors.New("boom"), stderr: expectedNoticeNotCollected}},
			check: func(t *testing.T, r runResult) {
				if r.state.CredentialsSetAt != nil {
					t.Error("a failed login must not set credentials_set_at")
				}
			},
		},
		{
			name:  "unreadable state shows no notice, collects nothing and keeps the file",
			setup: testSetup{state: State{Disabled: true}, loadErr: errors.New("yaml: invalid"), terminal: true},
			executions: []execution{
				{path: authLogin},
				{path: vmList},
			},
			check: func(t *testing.T, r runResult) {
				if r.saves != 0 || !r.state.Disabled {
					t.Error("the unreadable state file must not be overwritten")
				}
			},
		},
		{
			name:  "unwritable state on login shows no notice and collects nothing, run after run",
			setup: testSetup{terminal: true, saveErr: readOnly},
			executions: []execution{
				{path: authLogin},
				{path: authLogin},
			},
			check: func(t *testing.T, r runResult) {
				if r.saves != 2 || r.state.NoticeShown {
					t.Errorf("saves=%d notice_shown=%v, want one failed save per run and nothing persisted", r.saves, r.state.NoticeShown)
				}
			},
		},
		{
			name:  "unwritable state on another command shows no notice and collects nothing, run after run",
			setup: testSetup{terminal: true, saveErr: readOnly},
			executions: []execution{
				{path: vmList},
				{path: vmList},
			},
			check: func(t *testing.T, r runResult) {
				if r.saves != 2 || r.state.NoticeShown {
					t.Errorf("saves=%d notice_shown=%v, want one failed save per run and nothing persisted", r.saves, r.state.NoticeShown)
				}
			},
		},
		{
			name:       "first value is saved together with the installation id",
			setup:      collecting(State{CredentialsSetAt: &loginAt}),
			executions: []execution{{path: vmList, at: time.Second, events: 1}},
			check: func(t *testing.T, r runResult) {
				if r.saves != 1 || r.state.InstallationID == "" || !r.state.FirstValueRecorded {
					t.Errorf("saves=%d state=%+v, want a single save with id and TTFV", r.saves, r.state)
				}
			},
		},
		{
			name:  "installation id is generated once and reused",
			setup: collecting(State{}),
			executions: []execution{
				{path: vmList, events: 1},
				{path: vmList, events: 1},
			},
			check: func(t *testing.T, r runResult) {
				if r.saves != 1 || r.state.InstallationID != "id-1" {
					t.Errorf("saves=%d state=%+v, want the id generated and saved once", r.saves, r.state)
				}
				for _, e := range r.events {
					if e.InstallationID != "id-1" {
						t.Errorf("event carries id %q", e.InstallationID)
					}
				}
			},
		},
		{
			name: "installation id is still sent when saving fails",
			setup: func() testSetup {
				s := collecting(State{})
				s.saveErr = readOnly
				return s
			}(),
			executions: []execution{{path: vmList, events: 1}},
			check: func(t *testing.T, r runResult) {
				if r.events[0].InstallationID != "id-1" {
					t.Error("an ephemeral id must be used when saving fails")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := runExecutions(t, tc.setup, tc.executions)
			if tc.check != nil {
				tc.check(t, r)
			}
		})
	}
}

func TestServiceTimeToFirstValue(t *testing.T) {
	testCases := []struct {
		name         string
		state        State
		path         []string
		at           time.Duration
		cmdErr       error
		wantTTFV     time.Duration
		wantRecorded bool
	}{
		{"first successful infra command after login", State{CredentialsSetAt: &loginAt}, vmList, 4231 * time.Millisecond, nil, 4231 * time.Millisecond, true},
		{"first infra command fails", State{CredentialsSetAt: &loginAt}, vmList, time.Minute, errors.New("boom"), 0, false},
		{"already recorded", State{CredentialsSetAt: &loginAt, FirstValueRecorded: true}, vmList, time.Hour, nil, 0, true},
		{"never logged in", State{}, vmList, 0, nil, 0, false},
		{"clock moved backwards", State{CredentialsSetAt: &loginAt}, vmList, -time.Hour, nil, 0, false},
		{"non infrastructure command", State{CredentialsSetAt: &loginAt}, []string{"config", "list"}, time.Second, nil, 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := runExecutions(t, collecting(tc.state), []execution{{path: tc.path, at: tc.at, cmdErr: tc.cmdErr, events: 1}})

			var got time.Duration
			if ttfv := r.events[0].TimeToFirstValue; ttfv != nil {
				got = *ttfv
			}
			if got != tc.wantTTFV {
				t.Errorf("TTFV = %v, want %v", got, tc.wantTTFV)
			}
			if r.state.FirstValueRecorded != tc.wantRecorded {
				t.Errorf("first_value_recorded = %v, want %v", r.state.FirstValueRecorded, tc.wantRecorded)
			}
		})
	}
}

func TestServiceOptOut(t *testing.T) {
	seen := State{NoticeShown: true}
	withEnv := func(v string) map[string]string {
		return map[string]string{"MGC_CLI_TELEMETRY_OPTOUT": v, "DO_NOT_TRACK": v}
	}

	testCases := []struct {
		name           string
		setup          testSetup
		wantDisabledBy string
	}{
		{"config key", testSetup{state: State{NoticeShown: true, Disabled: true}, terminal: true}, DisabledByConfig},
		{"DO_NOT_TRACK", testSetup{state: seen, env: map[string]string{"DO_NOT_TRACK": "1"}, terminal: true}, DisabledByDoNotTrack},
		{"MGC_CLI_TELEMETRY_OPTOUT=1", testSetup{state: seen, env: map[string]string{"MGC_CLI_TELEMETRY_OPTOUT": "1"}, terminal: true}, DisabledByOptOut},
		{"MGC_CLI_TELEMETRY_OPTOUT=true", testSetup{state: seen, env: map[string]string{"MGC_CLI_TELEMETRY_OPTOUT": "true"}, terminal: true}, DisabledByOptOut},
		{"unreadable state", testSetup{state: seen, loadErr: errors.New("yaml: invalid"), terminal: true}, DisabledByUnreadableState},
		{"empty value keeps it enabled", testSetup{state: seen, env: withEnv(""), terminal: true}, ""},
		{"0 keeps it enabled", testSetup{state: seen, env: withEnv("0"), terminal: true}, ""},
		{"false keeps it enabled", testSetup{state: seen, env: withEnv("false"), terminal: true}, ""},
		{"FALSE keeps it enabled", testSetup{state: seen, env: withEnv("FALSE"), terminal: true}, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc, store, exp := newTestService(tc.setup)

			if got := svc.DisabledBy(); got != tc.wantDisabledBy {
				t.Errorf("DisabledBy = %q, want %q", got, tc.wantDisabledBy)
			}

			record(svc, commandAt(vmList, loginAt), nil)

			wantEvents := 1
			if tc.wantDisabledBy != "" {
				wantEvents = 0
				if store.saves != 0 {
					t.Error("nothing may be written when disabled")
				}
			}
			if len(exp.events) != wantEvents {
				t.Errorf("%d events, want %d", len(exp.events), wantEvents)
			}
		})
	}
}

func TestServiceSetDisabled(t *testing.T) {
	unreadable := errors.New("yaml: invalid")

	testCases := []struct {
		name           string
		setup          testSetup
		disabled       bool
		wantErr        bool
		wantDisabledBy string
	}{
		{"disable writes the config key", testSetup{}, true, false, DisabledByConfig},
		{"enable replaces an unreadable file", testSetup{loadErr: unreadable}, false, false, ""},
		{"enable that cannot save keeps blocking", testSetup{loadErr: unreadable, saveErr: errors.New("permission denied")}, false, true, DisabledByUnreadableState},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc, _, _ := newTestService(tc.setup)

			err := svc.SetDisabled(tc.disabled)
			if (err != nil) != tc.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if got := svc.DisabledBy(); got != tc.wantDisabledBy {
				t.Errorf("DisabledBy = %q, want %q", got, tc.wantDisabledBy)
			}
		})
	}
}

func TestServiceBuildsEvent(t *testing.T) {
	create := commandAt([]string{"virtual-machine", "instances", "create"}, loginAt)
	create.OptionsSet = []string{"region", "machine-type"}
	create.TenantID = "tenant-1"
	create.LastRequestID = "req-9"

	testCases := []struct {
		name          string
		info          CommandInfo
		cmdErr        error
		wantAction    string
		wantResource  *Resource
		wantOutcome   Outcome
		wantReason    FailureReason
		wantTenant    string
		wantRequestID string
	}{
		{
			name:          "command with flags, tenant and request id",
			info:          create,
			wantAction:    "virtualmachine.instances.create",
			wantResource:  &Resource{Type: "virtual_machine_instances"},
			wantOutcome:   OutcomeSuccess,
			wantTenant:    "tenant-1",
			wantRequestID: "req-9",
		},
		{
			name:        "unknown command with untyped error",
			info:        CommandInfo{UnknownCommand: true, Start: loginAt, End: loginAt},
			cmdErr:      errors.New(`unknown command "foo" for "mgc"`),
			wantAction:  UnknownAction,
			wantOutcome: OutcomeFailure,
			wantReason:  FailureValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setup := collecting(State{})
			setup.env = map[string]string{"GITHUB_ACTIONS": "true"}
			svc, _, exp := newTestService(setup)

			record(svc, tc.info, tc.cmdErr)
			e := exp.events[0]

			if e.Action != tc.wantAction || e.Outcome != tc.wantOutcome || e.FailureReason != tc.wantReason {
				t.Errorf("action/outcome/reason = %s/%s/%s, want %s/%s/%s", e.Action, e.Outcome, e.FailureReason, tc.wantAction, tc.wantOutcome, tc.wantReason)
			}
			if (e.Resource == nil) != (tc.wantResource == nil) || (e.Resource != nil && *e.Resource != *tc.wantResource) {
				t.Errorf("resource = %+v, want %+v", e.Resource, tc.wantResource)
			}
			var tenant string
			if e.Actor != nil {
				tenant = e.Actor.TenantID
			}
			if tenant != tc.wantTenant {
				t.Errorf("tenant = %q, want %q", tenant, tc.wantTenant)
			}
			if e.LastRequestID != tc.wantRequestID || e.Duration != tc.info.End.Sub(tc.info.Start) {
				t.Errorf("request id/duration = %q/%v", e.LastRequestID, e.Duration)
			}
			if e.ExecutionContext != ExecutionContextCI || e.ExecutionEnvironment != "github_actions" {
				t.Errorf("context = %s/%s, want ci/github_actions", e.ExecutionContext, e.ExecutionEnvironment)
			}
		})
	}
}
