package telemetry

import (
	"bytes"
	"context"
	"testing"
	"time"
)

type fakeStore struct {
	state   State
	saves   int
	saveErr error
	loadErr error
}

func (f *fakeStore) Load() (State, error) {
	if f.loadErr != nil {
		return State{}, f.loadErr
	}
	return f.state, nil
}

func (f *fakeStore) Save(s State) error {
	f.saves++
	if f.saveErr != nil {
		return f.saveErr
	}
	f.state = s
	return nil
}

type fakeExporter struct {
	events []Event
}

func (f *fakeExporter) Export(_ context.Context, e Event) error {
	f.events = append(f.events, e)
	return nil
}

func envFrom(vars map[string]string) func(string) string {
	return func(k string) string { return vars[k] }
}

var (
	loginAt   = time.Date(2026, 9, 22, 14, 0, 0, 0, time.UTC)
	vmList    = []string{"virtual-machine", "instances", "list"}
	authLogin = []string{"auth", "login"}
)

func commandAt(path []string, end time.Time) CommandInfo {
	return CommandInfo{Path: path, Start: end.Add(-842 * time.Millisecond), End: end}
}

// testSetup descreve o ambiente de um Service de teste. terminal = true com env vazio
// resulta em sessão interativa, onde o aviso de primeira execução pode aparecer.
type testSetup struct {
	state    State
	env      map[string]string
	terminal bool
	loadErr  error
	saveErr  error
}

func newTestService(s testSetup) (*Service, *fakeStore, *fakeExporter) {
	store := &fakeStore{state: s.state, loadErr: s.loadErr, saveErr: s.saveErr}
	exporter := &fakeExporter{}
	svc := New(Options{
		ClientVersion: "v1.4.2",
		OS:            "linux",
		IsTerminal:    s.terminal,
		Store:         store,
		Exporter:      exporter,
		Getenv:        envFrom(s.env),
		Now:           func() time.Time { return loginAt },
		NewID:         func() string { return "id-1" },
	})
	return svc, store, exporter
}

// record chama o Record e devolve o que foi escrito em stderr.
func record(svc *Service, info CommandInfo, cmdErr error) string {
	var stderr bytes.Buffer
	svc.Record(context.Background(), info, cmdErr, &stderr)
	return stderr.String()
}

func collecting(state State) testSetup {
	state.NoticeShown = true
	return testSetup{state: state, terminal: true}
}

type execution struct {
	path   []string
	at     time.Duration
	cmdErr error
	stderr string
	events int
}

type runResult struct {
	state  State
	saves  int
	events []Event
}

func runExecutions(t *testing.T, setup testSetup, executions []execution) runResult {
	t.Helper()
	result := runResult{state: setup.state}

	for i, ex := range executions {
		setup.state = result.state
		svc, store, exp := newTestService(setup)

		stderr := record(svc, commandAt(ex.path, loginAt.Add(ex.at)), ex.cmdErr)
		if stderr != ex.stderr {
			t.Errorf("execution %d: stderr = %q, want %q", i+1, stderr, ex.stderr)
		}
		if len(exp.events) != ex.events {
			t.Errorf("execution %d: %d events, want %d", i+1, len(exp.events), ex.events)
		}

		result.state = store.state
		result.saves += store.saves
		result.events = append(result.events, exp.events...)
	}
	return result
}
