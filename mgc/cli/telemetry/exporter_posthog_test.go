package telemetry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type capturedRequest struct {
	method  string
	path    string
	headers http.Header
	body    []byte
}

func newCaptureServer(t *testing.T, status int) (*httptest.Server, *[]capturedRequest) {
	t.Helper()
	var requests []capturedRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		requests = append(requests, capturedRequest{r.Method, r.URL.Path, r.Header.Clone(), body})
		w.WriteHeader(status)
	}))
	t.Cleanup(server.Close)
	return server, &requests
}

func exportedEvent() Event {
	ttfv := 1500 * time.Millisecond
	return Event{
		Timestamp:        time.Date(2026, 9, 22, 14, 30, 45, 123456789, time.UTC),
		Actor:            &Actor{TenantID: "tenant-1"},
		ExecutionContext: ExecutionContextInteractive,
		Resource:         &Resource{Type: "virtual_machine_instances", ID: "vm-1"},
		Action:           "virtualmachine.instances.get",
		OptionsSet:       []string{"output", "id"},
		Outcome:          OutcomeFailure,
		FailureReason:    FailureValidation,
		Duration:         842 * time.Millisecond,
		ClientVersion:    "v1.4.2",
		OS:               "linux",
		InstallMethod:    InstallMethodHomebrew,
		TimeToFirstValue: &ttfv,
		LastRequestID:    "req-1",
		InstallationID:   "b3f1c2de-0000-4000-8000-000000000000",
		Product:          "virtual machine",
	}
}

func decodeLogRecord(t *testing.T, body []byte) (otlpLogs, map[string]otlpAnyValue) {
	t.Helper()
	var payload otlpLogs
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid OTLP JSON: %v\n%s", err, body)
	}
	if len(payload.ResourceLogs) != 1 || len(payload.ResourceLogs[0].ScopeLogs) != 1 ||
		len(payload.ResourceLogs[0].ScopeLogs[0].LogRecords) != 1 {
		t.Fatalf("want exactly one log record, got %s", body)
	}

	attributes := map[string]otlpAnyValue{}
	for _, kv := range payload.ResourceLogs[0].ScopeLogs[0].LogRecords[0].Attributes {
		attributes[kv.Key] = kv.Value
	}
	return payload, attributes
}

func TestPostHogExporterRequest(t *testing.T) {
	server, requests := newCaptureServer(t, http.StatusOK)
	endpoint := server.URL + "/i/v1/logs"

	if err := NewPostHogExporter(endpoint, "phc_test").Export(context.Background(), exportedEvent()); err != nil {
		t.Fatalf("Export() = %v", err)
	}

	if len(*requests) != 1 {
		t.Fatalf("want 1 request, got %d", len(*requests))
	}
	req := (*requests)[0]

	if req.method != http.MethodPost || req.path != "/i/v1/logs" {
		t.Errorf("request = %s %s, want POST /i/v1/logs", req.method, req.path)
	}
	if got := req.headers.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q", got)
	}
	if got := req.headers.Values("Authorization"); len(got) != 1 || got[0] != "Bearer phc_test" {
		t.Errorf("Authorization = %q, want only the project token", got)
	}
	for _, header := range []string{"X-Api-Key", "X-Tenant-Id", "X-Request-Id", "Cookie"} {
		if got := req.headers.Get(header); got != "" {
			t.Errorf("header %s = %q, must not be sent", header, got)
		}
	}

	payload, _ := decodeLogRecord(t, req.body)
	resource := payload.ResourceLogs[0].Resource.Attributes
	if len(resource) != 1 || resource[0].Key != "service.name" || *resource[0].Value.StringValue != "mgc-cli" {
		t.Errorf("resource attributes = %+v, want service.name=mgc-cli", resource)
	}

	record := payload.ResourceLogs[0].ScopeLogs[0].LogRecords[0]
	if record.TimeUnixNano != "1790087445123000000" {
		t.Errorf("timeUnixNano = %q, want the event timestamp in ms precision", record.TimeUnixNano)
	}
	if record.Body.StringValue == nil || *record.Body.StringValue != "CLI - virtual machine - failure" {
		t.Errorf("body = %+v, want CLI - virtual machine - failure", record.Body)
	}
}

func TestPostHogExporterAttributes(t *testing.T) {
	success := exportedEvent()
	success.Outcome = OutcomeSuccess
	success.FailureReason = ""
	success.Actor = nil
	success.Resource = &Resource{Type: "virtual_machine_instances"}
	success.OptionsSet = nil
	success.TimeToFirstValue = nil
	success.LastRequestID = ""

	testCases := []struct {
		name         string
		event        Event
		wantSeverity string
		wantStrings  map[string]string
		wantInts     map[string]string
		wantAbsent   []string
	}{
		{
			name:         "failure is WARN and nested fields are flattened",
			event:        exportedEvent(),
			wantSeverity: "WARN",
			wantStrings: map[string]string{
				"distinct_id":      "b3f1c2de-0000-4000-8000-000000000000",
				"timestamp":        "2026-09-22T14:30:45.123Z",
				"actor.tenant_id":  "tenant-1",
				"resource.type":    "virtual_machine_instances",
				"resource.id":      "vm-1",
				"action":           "virtualmachine.instances.get",
				"optionsSet":       "id,output",
				"executionContext": "interactive",
				"outcome":          "failure",
				"failure_reason":   "validation",
				"logLevel":         "WARN",
				"impact":           "LOW",
				"clientVersion":    "v1.4.2",
				"os":               "linux",
				"installMethod":    "homebrew",
				"lastRequestId":    "req-1",
			},
			wantInts:   map[string]string{"durationMs": "842", "timeToFirstValueMs": "1500"},
			wantAbsent: []string{"actor", "resource", "executionEnvironment"},
		},
		{
			name:         "success is INFO and omitted model fields are omitted too",
			event:        success,
			wantSeverity: "INFO",
			wantStrings:  map[string]string{"outcome": "success", "resource.type": "virtual_machine_instances"},
			wantInts:     map[string]string{"durationMs": "842"},
			wantAbsent: []string{
				"failure_reason", "actor.tenant_id", "resource.id", "optionsSet",
				"timeToFirstValueMs", "lastRequestId", "executionEnvironment",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := newOTLPLogs(tc.event)
			if err != nil {
				t.Fatal(err)
			}
			payload, attributes := decodeLogRecord(t, body)

			if got := payload.ResourceLogs[0].ScopeLogs[0].LogRecords[0].SeverityText; got != tc.wantSeverity {
				t.Errorf("severityText = %q, want %q", got, tc.wantSeverity)
			}
			for key, want := range tc.wantStrings {
				if got := attributes[key].StringValue; got == nil || *got != want {
					t.Errorf("attribute %s = %v, want stringValue %q", key, got, want)
				}
			}
			for key, want := range tc.wantInts {
				if got := attributes[key].IntValue; got == nil || *got != want {
					t.Errorf("attribute %s = %v, want intValue %q", key, got, want)
				}
			}
			for _, key := range tc.wantAbsent {
				if _, ok := attributes[key]; ok {
					t.Errorf("attribute %s must be omitted", key)
				}
			}
		})
	}
}

func TestPostHogExporterErrors(t *testing.T) {
	testCases := []struct {
		name       string
		status     int
		wantStatus int
	}{
		{"2xx is accepted", http.StatusAccepted, 0},
		{"500 becomes an error", http.StatusInternalServerError, 500},
		{"401 becomes an error", http.StatusUnauthorized, 401},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, _ := newCaptureServer(t, tc.status)

			err := NewPostHogExporter(server.URL, "phc_test").Export(context.Background(), exportedEvent())

			var statusErr ExportStatusError
			switch {
			case tc.wantStatus == 0 && err != nil:
				t.Errorf("Export() = %v, want nil", err)
			case tc.wantStatus != 0 && (!errors.As(err, &statusErr) || statusErr.StatusCode != tc.wantStatus):
				t.Errorf("Export() = %v, want status error %d", err, tc.wantStatus)
			}
		})
	}
}

func TestPostHogExporterDebugLog(t *testing.T) {
	testCases := []struct {
		name        string
		status      int
		wantStatus  string
		wantSuccess string
	}{
		{"accepted call is logged as success", http.StatusOK, "200", "true"},
		{"rejected call is logged with its status", http.StatusInternalServerError, "500", "false"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, _ := newCaptureServer(t, tc.status)
			var logs []debugEntry
			exporter := NewPostHogExporter(server.URL, "phc_test").
				WithDebug(func(msg string, kv ...any) { logs = append(logs, debugEntry{msg, kv}) })

			_ = exporter.Export(context.Background(), exportedEvent())

			if len(logs) != 1 {
				t.Fatalf("debug logs = %+v, want a single entry", logs)
			}
			entry := logs[0]
			if entry.value("method") != http.MethodPost || entry.value("status") != tc.wantStatus || entry.value("success") != tc.wantSuccess {
				t.Errorf("debug log = %+v, want POST with status %s and success %s", entry, tc.wantStatus, tc.wantSuccess)
			}
			if entry.value("duration_ms") == "" {
				t.Errorf("debug log must carry duration_ms, got %+v", entry)
			}
			text := fmt.Sprint(entry.msg, entry.kv)
			for _, eventData := range []string{"tenant-1", "virtualmachine", "phc_test"} {
				if strings.Contains(text, eventData) {
					t.Errorf("debug log must not contain %q: %s", eventData, text)
				}
			}
		})
	}
}
