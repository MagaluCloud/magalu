package telemetry

import (
	"errors"
	"testing"
	"time"
)

// collecting parte de uma sessão interativa em que o aviso já foi visto, para os
// testes de coleta não serem afetados pela regra do aviso.
func collecting(state State) testSetup {
	state.NoticeShown = true
	return testSetup{state: state, terminal: true}
}

func TestRecordLoginShowsNoticeAndCollects(t *testing.T) {
	svc, store, exp := newTestService(testSetup{terminal: true})

	if stderr := record(svc, commandAt(authLogin, loginAt), nil); stderr == "" {
		t.Error("the notice must be shown after the first interactive login")
	}
	if len(exp.events) != 1 || exp.events[0].Action != "auth.login" {
		t.Fatalf("the login execution must be collected, got %+v", exp.events)
	}
	if store.state.InstallationID == "" {
		t.Error("installation_id must be written when the event is sent")
	}
}

func TestRecordAlreadyLoggedInSkipsOnlyFirstExecution(t *testing.T) {
	svc, store, exp := newTestService(testSetup{terminal: true})

	if stderr := record(svc, commandAt(vmList, loginAt), nil); stderr == "" {
		t.Error("the notice must be shown before the first collection")
	}
	if len(exp.events) != 0 {
		t.Error("the execution that shows the notice must not be collected")
	}
	if store.state.InstallationID != "" || store.state.FirstValueRecorded {
		t.Errorf("only notice_shown may be written, got %+v", store.state)
	}

	next, _, nextExp := newTestService(testSetup{state: store.state, terminal: true})
	record(next, commandAt(vmList, loginAt), nil)

	if len(nextExp.events) != 1 {
		t.Errorf("the next execution must be collected, got %d events", len(nextExp.events))
	}
}

func TestRecordFailedLoginIsTreatedAsAnyOtherCommand(t *testing.T) {
	svc, store, exp := newTestService(testSetup{terminal: true})

	record(svc, commandAt(authLogin, loginAt), errors.New("boom"))

	if len(exp.events) != 0 {
		t.Error("a failed login that shows the notice must not be collected")
	}
	if store.state.CredentialsSetAt != nil {
		t.Error("a failed login must not set credentials_set_at")
	}
}

func TestRecordOutsideInteractiveCollectsWithoutNotice(t *testing.T) {
	svc, store, exp := newTestService(testSetup{env: map[string]string{"CI": "true"}})

	if stderr := record(svc, commandAt(vmList, loginAt), nil); stderr != "" {
		t.Errorf("no notice outside interactive sessions, got %q", stderr)
	}
	if len(exp.events) != 1 || store.state.NoticeShown {
		t.Errorf("events=%d notice_shown=%v; want 1, false", len(exp.events), store.state.NoticeShown)
	}
}

func TestServiceTimeToFirstValue(t *testing.T) {
	t.Run("first successful infra command after login", func(t *testing.T) {
		svc, store, exp := newTestService(collecting(State{CredentialsSetAt: &loginAt}))

		record(svc, commandAt(vmList, loginAt.Add(4231*time.Millisecond)), nil)

		e := exp.events[0]
		if e.TimeToFirstValue == nil || e.TimeToFirstValue.Milliseconds() != 4231 {
			t.Fatalf("expected TTFV 4231ms, got %v", e.TimeToFirstValue)
		}
		if !store.state.FirstValueRecorded {
			t.Error("first_value_recorded must be persisted")
		}
	})

	t.Run("first infra command fails", func(t *testing.T) {
		svc, store, exp := newTestService(collecting(State{CredentialsSetAt: &loginAt}))

		record(svc, commandAt(vmList, loginAt.Add(time.Minute)), errors.New("boom"))

		if exp.events[0].TimeToFirstValue != nil || store.state.FirstValueRecorded {
			t.Error("failed command must not record TTFV")
		}
	})

	t.Run("already recorded", func(t *testing.T) {
		svc, _, exp := newTestService(collecting(State{CredentialsSetAt: &loginAt, FirstValueRecorded: true}))

		record(svc, commandAt(vmList, loginAt.Add(time.Hour)), nil)

		if exp.events[0].TimeToFirstValue != nil {
			t.Error("TTFV must be sent only once")
		}
	})

	t.Run("never logged in", func(t *testing.T) {
		svc, store, exp := newTestService(collecting(State{}))

		record(svc, commandAt(vmList, loginAt), nil)

		if exp.events[0].TimeToFirstValue != nil || store.state.FirstValueRecorded {
			t.Error("TTFV requires a recorded login")
		}
	})

	t.Run("clock moved backwards", func(t *testing.T) {
		svc, store, exp := newTestService(collecting(State{CredentialsSetAt: &loginAt}))

		record(svc, commandAt(vmList, loginAt.Add(-time.Hour)), nil)

		if exp.events[0].TimeToFirstValue != nil || store.state.FirstValueRecorded {
			t.Error("negative TTFV must be dropped and retried later")
		}
	})

	t.Run("non infrastructure command", func(t *testing.T) {
		svc, store, exp := newTestService(collecting(State{CredentialsSetAt: &loginAt}))

		record(svc, commandAt([]string{"config", "list"}, loginAt.Add(time.Second)), nil)

		if exp.events[0].TimeToFirstValue != nil || store.state.FirstValueRecorded {
			t.Error("config commands are not a first value")
		}
	})
}

func TestServiceOptOut(t *testing.T) {
	cases := []struct {
		name       string
		state      State
		env        map[string]string
		disabledBy string
	}{
		{"config key", State{Disabled: true}, nil, DisabledByConfig},
		{"DO_NOT_TRACK", State{}, map[string]string{"DO_NOT_TRACK": "1"}, DisabledByDoNotTrack},
		{"MGC_CLI_TELEMETRY_OPTOUT", State{}, map[string]string{"MGC_CLI_TELEMETRY_OPTOUT": "1"}, DisabledByOptOut},
		{"OPTOUT=true", State{}, map[string]string{"MGC_CLI_TELEMETRY_OPTOUT": "true"}, DisabledByOptOut},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Sessão interativa sem aviso visto: se a telemetria estivesse ativa,
			// o login mostraria o aviso e seria coletado.
			svc, store, exp := newTestService(testSetup{state: c.state, env: c.env, terminal: true})

			if got := svc.DisabledBy(); got != c.disabledBy {
				t.Errorf("DisabledBy = %q, want %q", got, c.disabledBy)
			}
			stderr := record(svc, commandAt(authLogin, loginAt), nil)
			stderr += record(svc, commandAt(vmList, loginAt), nil)

			if stderr != "" || len(exp.events) != 0 {
				t.Errorf("disabled telemetry must neither print nor export, got %q events=%d", stderr, len(exp.events))
			}
			if store.saves != 0 {
				t.Error("nothing may be written when disabled")
			}
		})
	}

	for _, v := range []string{"", "0", "false", "FALSE"} {
		setup := collecting(State{})
		setup.env = map[string]string{"MGC_CLI_TELEMETRY_OPTOUT": v, "DO_NOT_TRACK": v}
		svc, _, exp := newTestService(setup)

		record(svc, commandAt(vmList, loginAt), nil)

		if !svc.Enabled() || len(exp.events) != 1 {
			t.Errorf("value %q must keep telemetry enabled", v)
		}
	}
}

func TestServiceUnreadableState(t *testing.T) {
	unreadable := testSetup{state: State{Disabled: true}, loadErr: errors.New("yaml: invalid"), terminal: true}

	t.Run("sends nothing and keeps the file untouched", func(t *testing.T) {
		svc, store, exp := newTestService(unreadable)

		if got := svc.DisabledBy(); got != DisabledByUnreadableState {
			t.Errorf("DisabledBy = %q, want %q", got, DisabledByUnreadableState)
		}
		stderr := record(svc, commandAt(authLogin, loginAt), nil)
		stderr += record(svc, commandAt(vmList, loginAt), nil)

		if stderr != "" || len(exp.events) != 0 {
			t.Errorf("no notice or event without a readable state, got %q events=%d", stderr, len(exp.events))
		}
		if store.saves != 0 || !store.state.Disabled {
			t.Error("the unreadable state file must not be overwritten")
		}
	})

	t.Run("explicit enable replaces the unreadable file", func(t *testing.T) {
		svc, store, _ := newTestService(unreadable)

		if err := svc.SetDisabled(false); err != nil {
			t.Fatal(err)
		}
		if !svc.Enabled() || store.saves != 1 {
			t.Errorf("Enabled = %v, saves = %d; want true, 1", svc.Enabled(), store.saves)
		}
	})

	t.Run("failed enable keeps blocking", func(t *testing.T) {
		svc, store, _ := newTestService(unreadable)
		store.saveErr = errors.New("permission denied")

		if err := svc.SetDisabled(false); err == nil {
			t.Fatal("expected save error")
		}
		if svc.Enabled() {
			t.Error("telemetry must stay disabled while the state is unreadable")
		}
	})
}

func TestRecordUnwritableStateShowsNoNoticeAndCollectsNothing(t *testing.T) {
	readOnly := errors.New("read-only home")

	for name, path := range map[string][]string{"login": authLogin, "other command": vmList} {
		first, store, exp := newTestService(testSetup{terminal: true, saveErr: readOnly})

		if stderr := record(first, commandAt(path, loginAt), nil); stderr != "" || len(exp.events) != 0 {
			t.Errorf("%s: got %q events=%d; want no notice and no event", name, stderr, len(exp.events))
		}
		if store.saves != 1 || store.state.NoticeShown {
			t.Errorf("%s: saves=%d notice_shown=%v; want one failed save and nothing persisted", name, store.saves, store.state.NoticeShown)
		}

		next, _, nextExp := newTestService(testSetup{state: store.state, terminal: true, saveErr: readOnly})

		if stderr := record(next, commandAt(path, loginAt), nil); stderr != "" || len(nextExp.events) != 0 {
			t.Errorf("%s, next execution: got %q events=%d; want the notice not to repeat", name, stderr, len(nextExp.events))
		}
	}
}

func TestServiceSavesOncePerCommand(t *testing.T) {
	t.Run("first value", func(t *testing.T) {
		svc, store, _ := newTestService(collecting(State{CredentialsSetAt: &loginAt}))

		record(svc, commandAt(vmList, loginAt.Add(time.Second)), nil)

		if store.saves != 1 || store.state.InstallationID == "" || !store.state.FirstValueRecorded {
			t.Errorf("saves = %d, state = %+v; want a single save with id and TTFV", store.saves, store.state)
		}
	})

	t.Run("first interactive login", func(t *testing.T) {
		svc, store, _ := newTestService(testSetup{terminal: true})

		record(svc, commandAt(authLogin, loginAt), nil)

		if store.saves != 1 || store.state.InstallationID == "" || !store.state.NoticeShown || store.state.CredentialsSetAt == nil {
			t.Errorf("saves = %d, state = %+v; want a single save with id, notice and credentials", store.saves, store.state)
		}
	})
}

func TestServiceInstallationID(t *testing.T) {
	svc, store, exp := newTestService(collecting(State{}))

	record(svc, commandAt(vmList, loginAt), nil)
	record(svc, commandAt(vmList, loginAt), nil)

	if store.state.InstallationID != "id-1" || store.saves != 1 {
		t.Errorf("id must be generated and saved once, state=%+v saves=%d", store.state, store.saves)
	}
	for _, e := range exp.events {
		if e.InstallationID != "id-1" {
			t.Errorf("event carries id %q", e.InstallationID)
		}
	}

	failing := collecting(State{})
	failing.saveErr = errors.New("read-only")
	svc, _, exp = newTestService(failing)

	record(svc, commandAt(vmList, loginAt), nil)

	if exp.events[0].InstallationID != "id-1" {
		t.Error("an ephemeral id must still be used when saving fails")
	}
}

func TestServiceBuildsEvent(t *testing.T) {
	setup := collecting(State{})
	setup.env = map[string]string{"GITHUB_ACTIONS": "true"}
	svc, _, exp := newTestService(setup)

	info := commandAt([]string{"virtual-machine", "instances", "create"}, loginAt)
	info.OptionsSet = []string{"region", "machine-type"}
	info.TenantID = "tenant-1"
	info.LastRequestID = "req-9"
	record(svc, info, nil)

	e := exp.events[0]
	if e.Action != "virtualmachine.instances.create" || e.Resource == nil || e.Resource.Type != "virtual_machine_instances" {
		t.Errorf("unexpected action/resource: %s %+v", e.Action, e.Resource)
	}
	if e.Duration != 842*time.Millisecond || e.Outcome != OutcomeSuccess || e.Actor.TenantID != "tenant-1" || e.LastRequestID != "req-9" {
		t.Errorf("unexpected event %+v", e)
	}
	if e.ExecutionContext != ExecutionContextCI || e.ExecutionEnvironment != "github_actions" {
		t.Errorf("unexpected context %s/%s", e.ExecutionContext, e.ExecutionEnvironment)
	}

	record(svc, CommandInfo{UnknownCommand: true, Start: loginAt, End: loginAt}, errors.New(`unknown command "foo" for "mgc"`))

	e = exp.events[1]
	if e.Action != UnknownAction || e.Resource != nil || e.FailureReason != FailureValidation {
		t.Errorf("unknown command: %+v", e)
	}
}
