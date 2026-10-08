package dispatch

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type dispatchRun struct {
	cmd     *exec.Cmd
	payload []byte
	mode    os.FileMode
	files   []string
	logs    []debugEntry
	err     error
}

// helper para nao rodar um subprocesso real.
// runDispatch monta um Subprocess com fakes para tudo que sai do processo e deixa o
// teste trocar o que precisar antes do Dispatch
func runDispatch(t *testing.T, debugLogPath string, override func(*Subprocess)) dispatchRun {
	t.Helper()
	dir := t.TempDir()
	var run dispatchRun

	s := NewSubprocess(
		func() string { return debugLogPath },
		func(msg string, kv ...any) { run.logs = append(run.logs, debugEntry{msg, kv}) },
	)
	s.tempDir = func() string { return dir }
	s.executable = func() (string, error) { return "/usr/local/bin/mgc", nil }
	s.start = func(cmd *exec.Cmd) error {
		run.cmd = cmd
		path := cmd.Args[len(cmd.Args)-1]
		run.payload, _ = os.ReadFile(path)
		if info, err := os.Stat(path); err == nil {
			run.mode = info.Mode().Perm()
		}
		return nil
	}
	if override != nil {
		override(s)
	}

	run.err = s.Dispatch(testEvent())

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		run.files = append(run.files, e.Name())
	}
	return run
}

// nao testa a execucao mas como o comando é montado
func TestSubprocessDispatch(t *testing.T) {
	testCases := []struct {
		name         string
		debugLogPath string
	}{
		{"without debug", ""},
		{"with debug", "/home/user/.config/mgc/telemetry-debug.txt"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run := runDispatch(t, tc.debugLogPath, nil)
			if run.err != nil {
				t.Fatalf("Dispatch() = %v", run.err)
			}

			cmd := run.cmd
			if len(cmd.Args) != 3 || cmd.Args[0] != "/usr/local/bin/mgc" || cmd.Args[1] != SenderArg {
				t.Errorf("args = %v, want [/usr/local/bin/mgc %s <payload>]", cmd.Args, SenderArg)
			}
			if !strings.HasPrefix(filepath.Base(cmd.Args[2]), "mgc-telemetry-") {
				t.Errorf("payload file = %q, want the mgc-telemetry- prefix", cmd.Args[2])
			}
			if cmd.Dir == "" {
				t.Errorf("the sender must not inherit the user's working directory")
			}
			if cmd.Stdin != nil || cmd.Stdout != nil || cmd.Stderr != nil {
				t.Errorf("stdin, stdout and stderr must stay detached from the terminal")
			}
			if cmd.SysProcAttr == nil {
				t.Errorf("the sender must be started detached")
			}
			if runtime.GOOS != "windows" && run.mode != payloadFileMode {
				t.Errorf("payload mode = %v, want %v", run.mode, os.FileMode(payloadFileMode))
			}

			payload, err := Decode(run.payload)
			if err != nil {
				t.Fatalf("payload must decode: %v", err)
			}
			if payload.DebugLogPath != tc.debugLogPath {
				t.Errorf("debug log path = %q, want %q", payload.DebugLogPath, tc.debugLogPath)
			}
			if len(run.files) != 1 {
				t.Errorf("the payload file must stay for the sender, got %v", run.files)
			}
		})
	}
}

func TestSubprocessDispatchFailures(t *testing.T) {
	errStart := errors.New("fork/exec: permission denied")
	errExecutable := errors.New("executable not found")

	testCases := []struct {
		name     string
		override func(*Subprocess)
		wantErr  error
	}{
		{
			name:     "start refused",
			override: func(s *Subprocess) { s.start = func(*exec.Cmd) error { return errStart } },
			wantErr:  errStart,
		},
		{
			name:     "executable not found",
			override: func(s *Subprocess) { s.executable = func() (string, error) { return "", errExecutable } },
			wantErr:  errExecutable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run := runDispatch(t, "", tc.override)

			if !errors.Is(run.err, tc.wantErr) {
				t.Errorf("Dispatch() = %v, want %v", run.err, tc.wantErr)
			}
			if len(run.files) != 0 {
				t.Errorf("no payload file may be left behind, got %v", run.files)
			}
			if len(run.logs) != 0 {
				t.Errorf("a failed dispatch is logged by the service, got %+v", run.logs)
			}
		})
	}
}

func TestSubprocessDispatchLog(t *testing.T) {
	testCases := []struct {
		name         string
		debugLogPath string
		wantDebugLog bool
	}{
		{"without debug", "", false},
		{"with debug", "/home/user/.config/mgc/telemetry-debug.txt", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run := runDispatch(t, tc.debugLogPath, nil)

			if len(run.logs) != 1 || run.logs[0].msg != "telemetry: dispatched" {
				t.Fatalf("logs = %+v, want a single telemetry: dispatched", run.logs)
			}
			entry := run.logs[0]
			if _, ok := entry.value("pid"); !ok {
				t.Errorf("the log must carry the pid: %+v", entry)
			}
			if size, _ := entry.value("bytes"); size != fmt.Sprint(len(run.payload)) {
				t.Errorf("bytes = %s, want %d", size, len(run.payload))
			}
			if path, ok := entry.value("debug_log"); ok != tc.wantDebugLog || path != tc.debugLogPath {
				t.Errorf("debug_log = %q (present %v), want %q", path, ok, tc.debugLogPath)
			}

			text := fmt.Sprint(entry.msg, entry.kv)
			for _, eventData := range []string{"tenant-1", "virtualmachine", "b3f1c2de", "vm-1"} {
				if strings.Contains(text, eventData) {
					t.Errorf("the log must not contain event data %q: %s", eventData, text)
				}
			}
		})
	}
}
