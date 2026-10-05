package telemetry

import (
	"context"
	"errors"
	"testing"
	"time"
)

type exporterFunc func(ctx context.Context, event Event) error

func (f exporterFunc) Export(ctx context.Context, event Event) error {
	return f(ctx, event)
}

func slowExporter(delay time.Duration, honorContext bool) Exporter {
	return exporterFunc(func(ctx context.Context, _ Event) error {
		if !honorContext {
			time.Sleep(delay)
			return nil
		}
		select {
		case <-time.After(delay):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
}

func TestMultiExporter(t *testing.T) {
	errFirst := errors.New("first failed")
	errSecond := errors.New("second failed")

	testCases := []struct {
		name        string
		exporters   MultiExporter
		maxDuration time.Duration
		wantErrs    []error
	}{
		{
			name:        "two slow exporters share the same export budget",
			exporters:   MultiExporter{slowExporter(2*time.Second, true), slowExporter(2*time.Second, true)},
			maxDuration: exportBudget,
			wantErrs:    []error{context.DeadlineExceeded},
		},
		{
			name:        "an exporter that ignores the context does not hold the budget",
			exporters:   MultiExporter{slowExporter(2*time.Second, true), slowExporter(2*time.Second, false)},
			maxDuration: exportBudget,
			wantErrs:    []error{context.DeadlineExceeded},
		},
		{
			name: "every error is reported",
			exporters: MultiExporter{
				exporterFunc(func(context.Context, Event) error { return errFirst }),
				exporterFunc(func(context.Context, Event) error { return errSecond }),
			},
			maxDuration: 100 * time.Millisecond,
			wantErrs:    []error{errFirst, errSecond},
		},
		{
			name:        "every exporter receives the event",
			exporters:   MultiExporter{&fakeExporter{}, &fakeExporter{}},
			maxDuration: 100 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), ExportTimeout)
			defer cancel()

			start := time.Now()
			err := tc.exporters.Export(ctx, Event{Action: "virtualmachine.instances.list"})
			elapsed := time.Since(start)

			if elapsed > tc.maxDuration {
				t.Errorf("Export took %v, want at most %v", elapsed, tc.maxDuration)
			}
			if len(tc.wantErrs) == 0 && err != nil {
				t.Errorf("Export() = %v, want nil", err)
			}
			for _, want := range tc.wantErrs {
				if !errors.Is(err, want) {
					t.Errorf("Export() = %v, want it to include %v", err, want)
				}
			}
			for _, exporter := range tc.exporters {
				if fake, ok := exporter.(*fakeExporter); ok && len(fake.events) != 1 {
					t.Errorf("each exporter must receive the event once, got %d", len(fake.events))
				}
			}
		})
	}
}
