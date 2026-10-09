package sender

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

var (
	errConnRefused = &url.Error{Op: "Post", URL: "https://eu.i.posthog.com/i/v1/logs", Err: &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}}
	errDNS         = &url.Error{Op: "Post", URL: "https://eu.i.posthog.com/i/v1/logs", Err: &net.DNSError{Err: "no such host", Name: "eu.i.posthog.com"}}
)

// testEvent tem dados que nunca podem aparecer nos logs do processo de envio
func testEvent() telemetry.Event {
	return telemetry.Event{
		Timestamp:      time.Date(2026, 9, 22, 14, 30, 45, 0, time.UTC),
		Actor:          &telemetry.Actor{TenantID: "tenant-secret"},
		Action:         "virtualmachine.instances.list",
		Outcome:        telemetry.OutcomeSuccess,
		InstallationID: "installation-secret",
	}
}

func status(code int) error {
	return telemetry.ExportStatusError{StatusCode: code}
}

type logEntry struct {
	msg string
	kv  []any
}

func (e logEntry) value(key string) string {
	for i := 0; i+1 < len(e.kv); i += 2 {
		if e.kv[i] == key {
			return fmt.Sprint(e.kv[i+1])
		}
	}
	return ""
}

// logRecorder guarda as linhas de log e pode ser usado por várias goroutines
type logRecorder struct {
	mu      sync.Mutex
	entries []logEntry
}

func (r *logRecorder) log(msg string, kv ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, logEntry{msg, kv})
}

func (r *logRecorder) last() logEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.entries) == 0 {
		return logEntry{}
	}
	return r.entries[len(r.entries)-1]
}

// scriptedExporter devolve um resultado por chamada e repete o último quando acabam
type scriptedExporter struct {
	mu      sync.Mutex
	results []func(ctx context.Context) error
	calls   int
}

func (s *scriptedExporter) Export(ctx context.Context, _ telemetry.Event) error {
	s.mu.Lock()
	result := s.results[min(s.calls, len(s.results)-1)]
	s.calls++
	s.mu.Unlock()
	return result(ctx)
}

func (s *scriptedExporter) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

func returns(err error) func(context.Context) error {
	return func(context.Context) error { return err }
}

func blocksUntilTimeout(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

// timedExporter anota quando recebeu o evento, para medir se um exportador esperou o outro
type timedExporter struct {
	mu       sync.Mutex
	received []time.Time
}

func (e *timedExporter) Export(context.Context, telemetry.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.received = append(e.received, time.Now())
	return nil
}

type panickingExporter struct{}

func (panickingExporter) Export(context.Context, telemetry.Event) error {
	panic("exporter bug")
}
