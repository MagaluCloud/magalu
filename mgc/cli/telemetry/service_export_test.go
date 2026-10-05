package telemetry

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type debugEntry struct {
	msg string
	kv  []any
}

func (d debugEntry) value(key string) string {
	for i := 0; i+1 < len(d.kv); i += 2 {
		if d.kv[i] == key {
			return fmt.Sprint(d.kv[i+1])
		}
	}
	return ""
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	stdout, stderr := os.Stdout, os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout, os.Stderr = w, w
	defer func() { os.Stdout, os.Stderr = stdout, stderr }()

	fn()

	_ = w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestServiceExportTimeout(t *testing.T) {
	release := make(chan struct{})
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(5 * time.Second):
		case <-release:
		}
	}))
	defer slow.Close()
	defer close(release)

	failing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failing.Close()

	testCases := []struct {
		name        string
		exporter    Exporter
		maxDuration time.Duration
		wantReason  string
	}{
		{
			name:        "slow backend is aborted after the 1s budget",
			exporter:    NewPostHogExporter(slow.URL, "phc_test"),
			maxDuration: 1100 * time.Millisecond,
			wantReason:  "timeout",
		},
		{
			name:        "slow backends registered together share the 1s budget",
			exporter:    MultiExporter{NewPostHogExporter(slow.URL, "phc_test"), NewPostHogExporter(slow.URL, "phc_test")},
			maxDuration: 1100 * time.Millisecond,
			wantReason:  "timeout",
		},
		{
			name:        "non 2xx response is dropped with its status",
			exporter:    NewPostHogExporter(failing.URL, "phc_test"),
			maxDuration: 500 * time.Millisecond,
			wantReason:  "status 500",
		},
		{
			name:        "unreachable backend is dropped as a network failure",
			exporter:    NewPostHogExporter("http://127.0.0.1:1/i/v1/logs", "phc_test"),
			maxDuration: 1100 * time.Millisecond,
			wantReason:  "network",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var logs []debugEntry
			svc := New(Options{
				Store:    &fakeStore{state: State{NoticeShown: true, InstallationID: "id-1"}},
				Exporter: tc.exporter,
				Getenv:   envFrom(nil),
				Debug:    func(msg string, kv ...any) { logs = append(logs, debugEntry{msg, kv}) },
			})

			info := commandAt(vmList, loginAt)
			info.TenantID = "tenant-secret"
			info.OptionsSet = []string{"name"}

			var notice bytes.Buffer
			var elapsed time.Duration
			output := captureOutput(t, func() {
				start := time.Now()
				svc.Record(context.Background(), info, nil, &notice)
				elapsed = time.Since(start)
			})

			if elapsed > tc.maxDuration {
				t.Errorf("Record took %v, want at most %v", elapsed, tc.maxDuration)
			}
			if output != "" || notice.Len() != 0 {
				t.Errorf("export failures must be silent, got stdout/stderr %q and writer %q", output, notice.String())
			}
			if len(logs) != 1 || logs[0].value("reason") != tc.wantReason {
				t.Fatalf("debug logs = %+v, want a single entry with reason %q", logs, tc.wantReason)
			}
			entry := fmt.Sprint(logs[0].msg, logs[0].kv)
			for _, eventData := range []string{"tenant-secret", "virtualmachine", "id-1", "phc_test"} {
				if strings.Contains(entry, eventData) {
					t.Errorf("debug log must not contain event data %q: %s", eventData, entry)
				}
			}
		})
	}
}
