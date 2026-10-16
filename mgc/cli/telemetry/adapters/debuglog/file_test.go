package debuglog

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var at = time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)

func fixedNow() time.Time { return at }

func TestFormatLine(t *testing.T) {
	testCases := []struct {
		name string
		msg  string
		kv   []any
		want string
	}{
		{
			name: "plain values",
			msg:  "telemetry: attempt",
			kv:   []any{"attempt", 2, "status", 503, "next_wait_ms", 1500},
			want: "2026-10-07T12:00:00.000Z pid=42 telemetry: attempt attempt=2 status=503 next_wait_ms=1500\n",
		},
		{
			name: "values with spaces are quoted",
			msg:  "telemetry: attempt",
			kv:   []any{"error", "dial tcp: connection refused"},
			want: "2026-10-07T12:00:00.000Z pid=42 telemetry: attempt error=\"dial tcp: connection refused\"\n",
		},
		{
			name: "empty value is quoted",
			msg:  "telemetry: sent",
			kv:   []any{"reason", ""},
			want: "2026-10-07T12:00:00.000Z pid=42 telemetry: sent reason=\"\"\n",
		},
		{
			name: "odd key without value",
			msg:  "telemetry: dropped",
			kv:   []any{"orphan"},
			want: "2026-10-07T12:00:00.000Z pid=42 telemetry: dropped orphan\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatLine(at, 42, tc.msg, tc.kv); got != tc.want {
				t.Errorf("got  %q\nwant %q", got, tc.want)
			}
		})
	}
}

func TestNewAppendsAndTruncates(t *testing.T) {
	nearLimit := strings.Repeat("x", maxDebugLogBytes-10)

	testCases := []struct {
		name      string
		existing  string // "" significa que o arquivo não existe
		loggers   int
		wantLines int
		wantOld   bool
	}{
		{"missing file is created", "", 1, 1, false},
		{"two senders keep both lines in order", "", 2, 2, false},
		{"existing history is kept", "old line\n", 1, 2, true},
		{"file near the limit is reset first", nearLimit, 1, 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config", FileName)
			if tc.existing != "" {
				if err := os.MkdirAll(filepath.Dir(path), dirMode); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(tc.existing), fileMode); err != nil {
					t.Fatal(err)
				}
			}

			for i := range tc.loggers {
				New(path, fixedNow, 100+i)("telemetry: sent", "attempt", 1)
			}

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			text := string(content)
			if lines := strings.Count(text, "\n"); lines != tc.wantLines {
				t.Errorf("%d lines, want %d", lines, tc.wantLines)
			}
			if strings.Contains(text, "old line") != tc.wantOld {
				t.Errorf("old history present = %v, want %v", !tc.wantOld, tc.wantOld)
			}
			if tc.loggers == 2 && strings.Index(text, "pid=100") > strings.Index(text, "pid=101") {
				t.Errorf("lines out of order: %s", text)
			}
			if len(content) > maxDebugLogBytes {
				t.Errorf("file has %d bytes, more than the %d limit", len(content), maxDebugLogBytes)
			}
			if info, err := os.Stat(path); err == nil && runtime.GOOS != "windows" && info.Mode().Perm() != fileMode && tc.existing == "" {
				t.Errorf("mode = %v, want %v", info.Mode().Perm(), os.FileMode(fileMode))
			}
		})
	}
}

func TestNewIgnoresWriteFailures(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("needs a directory the current user cannot write to")
	}
	dir := t.TempDir()
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, dirMode) })
	path := filepath.Join(dir, FileName)

	New(path, fixedNow, 1)("telemetry: sent", "attempt", 1)

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("no file can be created in a read-only directory, stat err = %v", err)
	}
}
