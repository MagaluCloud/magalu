package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry/adapters/debuglog"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry/adapters/dispatch"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry/sender"
)

type senderExporter struct {
	events []telemetry.Event
	panics bool
}

func (e *senderExporter) Export(_ context.Context, event telemetry.Event) error {
	if e.panics {
		panic("exporter bug")
	}
	e.events = append(e.events, event)
	return nil
}

func TestRunTelemetrySender(t *testing.T) {
	event := telemetry.Event{
		Timestamp:      time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC),
		Actor:          &telemetry.Actor{TenantID: "tenant-secret"},
		Action:         "virtualmachine.instances.list",
		Outcome:        telemetry.OutcomeSuccess,
		InstallationID: "installation-1",
		Product:        "virtual machine",
	}

	testCases := []struct {
		name          string
		debug         bool
		noPayload     bool
		extraArg      bool
		panics        bool
		factoryPanics bool
		wantEvents    int
		wantLogged    []string
	}{
		{name: "sends the event", wantEvents: 1},
		{name: "debug mode writes the log file", debug: true, wantEvents: 1, wantLogged: []string{"telemetry: sent"}},
		{name: "missing payload file", noPayload: true},
		{name: "unexpected arguments", extraArg: true},
		{name: "exporter panic is contained", panics: true, debug: true, wantLogged: []string{"telemetry: dropped", "reason=panic"}},
		{name: "sender panic is logged in debug mode", factoryPanics: true, debug: true, wantLogged: []string{"telemetry: dropped", "reason=panic", "error=\"runtime error: slice bounds out of range"}},
		{name: "sender panic without debug writes nothing", factoryPanics: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			payloadPath := filepath.Join(dir, "mgc-telemetry-test.json")
			logFile := filepath.Join(dir, "config", debuglog.FileName)
			debugLogPath := ""
			if tc.debug {
				debugLogPath = logFile
			}
			if !tc.noPayload {
				data, err := dispatch.Encode(event, debugLogPath)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(payloadPath, data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			args := []string{payloadPath}
			if tc.extraArg {
				args = append(args, "extra")
			}

			exporter := &senderExporter{panics: tc.panics}
			exporters := func(telemetry.DebugLogger) []telemetry.Exporter {
				if tc.factoryPanics {
					var none []telemetry.Exporter
					return none[:1]
				}
				return []telemetry.Exporter{exporter}
			}

			// Um panic que escapasse derrubaria o teste, então chegar aqui já prova o recover
			runTelemetrySender(args, exporters, sender.DefaultPolicy(), time.Second)

			if len(exporter.events) != tc.wantEvents {
				t.Fatalf("%d events, want %d", len(exporter.events), tc.wantEvents)
			}
			if tc.wantEvents == 1 {
				got := exporter.events[0]
				if got.Action != event.Action || got.InstallationID != event.InstallationID || got.Product != event.Product {
					t.Errorf("event = %+v, want action, installation id and product preserved", got)
				}
			}
			if _, err := os.Stat(payloadPath); !tc.extraArg && !os.IsNotExist(err) {
				t.Errorf("the payload file must be deleted, stat err = %v", err)
			}

			content, err := os.ReadFile(logFile)
			if !tc.debug && err == nil {
				t.Errorf("without debug no log file may be written, got %q", content)
			}
			for _, want := range tc.wantLogged {
				if !strings.Contains(string(content), want) {
					t.Errorf("debug log must contain %q, got %q (err %v)", want, content, err)
				}
			}
			if strings.Contains(string(content), "tenant-secret") {
				t.Errorf("debug log must not contain event data: %s", content)
			}
		})
	}
}
