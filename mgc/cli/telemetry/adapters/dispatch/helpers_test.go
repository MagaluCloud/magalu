package dispatch

import (
	"fmt"
	"time"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

func testEvent() telemetry.Event {
	ttfv := 1500 * time.Millisecond
	return telemetry.Event{
		Timestamp:        time.Date(2026, 9, 22, 14, 30, 45, 123000000, time.UTC),
		Actor:            &telemetry.Actor{TenantID: "tenant-1"},
		ExecutionContext: telemetry.ExecutionContextInteractive,
		Resource:         &telemetry.Resource{Type: "virtual_machine_instances", ID: "vm-1"},
		Action:           "virtualmachine.instances.get",
		OptionsSet:       []string{"id", "output"},
		Outcome:          telemetry.OutcomeFailure,
		FailureReason:    telemetry.FailureValidation,
		Duration:         842 * time.Millisecond,
		ClientVersion:    "v1.4.2",
		OS:               "linux",
		InstallMethod:    telemetry.InstallMethodHomebrew,
		TimeToFirstValue: &ttfv,
		LastRequestID:    "req-1",
		InstallationID:   "b3f1c2de-0000-4000-8000-000000000000",
		Product:          "virtual machine",
	}
}

type debugEntry struct {
	msg string
	kv  []any
}

func (d debugEntry) value(key string) (string, bool) {
	for i := 0; i+1 < len(d.kv); i += 2 {
		if d.kv[i] == key {
			return fmt.Sprint(d.kv[i+1]), true
		}
	}
	return "", false
}
