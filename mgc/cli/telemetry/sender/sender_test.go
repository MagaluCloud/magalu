package sender

import (
	"context"
	"testing"
	"time"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

// independentMargin é a folga do agendador nos testes que medem tempo
const independentMargin = 200 * time.Millisecond

func TestSendRunsExportersIndependently(t *testing.T) {
	const (
		baseDelay     = 100 * time.Millisecond
		independentBy = 50 * time.Millisecond
	)

	policy := DefaultPolicy()
	policy.BaseDelay = baseDelay
	policy.Rand = func() float64 { return 0.5 }

	failing := &scriptedExporter{results: []func(context.Context) error{returns(errConnRefused)}}
	healthy := &timedExporter{}
	logs := &logRecorder{}

	start := time.Now()
	policy.Send(context.Background(), testEvent(), []telemetry.Exporter{failing, healthy}, logs.log)
	elapsed := time.Since(start)

	if len(healthy.received) != 1 {
		t.Fatalf("healthy exporter got %d events, want 1", len(healthy.received))
	}
	if waited := healthy.received[0].Sub(start); waited > independentBy {
		t.Errorf("healthy exporter waited %v for the failing one", waited)
	}
	if failing.callCount() != policy.Attempts {
		t.Errorf("failing exporter got %d attempts, want %d", failing.callCount(), policy.Attempts)
	}
	if minimum := baseDelay + 2*baseDelay + 4*baseDelay; elapsed < minimum {
		t.Errorf("Send returned after %v, it must wait for every exporter (at least %v)", elapsed, minimum)
	}

	exporters := map[string]bool{}
	for _, entry := range logs.entries {
		exporters[entry.value("exporter")] = true
	}
	if !exporters["0"] || !exporters["1"] {
		t.Errorf("every log line must say which exporter it is about, got %v", exporters)
	}
}

func TestSendStopsAtTheSenderDeadline(t *testing.T) {
	const deadline = 100 * time.Millisecond

	stuck := &scriptedExporter{results: []func(context.Context) error{blocksUntilTimeout}}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	start := time.Now()
	DefaultPolicy().Send(ctx, testEvent(), []telemetry.Exporter{stuck}, (&logRecorder{}).log)

	if elapsed := time.Since(start); elapsed > deadline+independentMargin {
		t.Errorf("Send took %v, the sender deadline is %v", elapsed, deadline)
	}
	if stuck.callCount() != 1 {
		t.Errorf("%d attempts after the deadline, want 1", stuck.callCount())
	}
}

func TestSendSurvivesExporterPanic(t *testing.T) {
	healthy := &timedExporter{}
	logs := &logRecorder{}

	DefaultPolicy().Send(context.Background(), testEvent(), []telemetry.Exporter{panickingExporter{}, healthy}, logs.log)

	if len(healthy.received) != 1 {
		t.Errorf("healthy exporter got %d events, want 1", len(healthy.received))
	}
	var dropped bool
	for _, entry := range logs.entries {
		if entry.msg == "telemetry: dropped" && entry.value("reason") == "panic" && entry.value("exporter") == "0" {
			dropped = true
		}
	}
	if !dropped {
		t.Errorf("the panic must be logged as a drop of exporter 0, got %+v", logs.entries)
	}
}
