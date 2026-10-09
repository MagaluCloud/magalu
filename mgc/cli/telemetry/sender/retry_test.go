package sender

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// testPolicy grava as esperas em vez de dormir. Com Rand em 0,5 o jitter some e as
// esperas ficam exatamente em 1 s, 2 s e 4 s
func testPolicy(waits *[]time.Duration, sleepErr error) Policy {
	policy := DefaultPolicy()
	policy.Rand = func() float64 { return 0.5 }
	policy.Sleep = func(_ context.Context, d time.Duration) error {
		*waits = append(*waits, d)
		return sleepErr
	}
	return policy
}

func TestRetry(t *testing.T) {
	testCases := []struct {
		name           string
		results        []func(context.Context) error
		attemptTimeout time.Duration
		sleepErr       error
		wantCalls      int
		wantWaits      []time.Duration
		wantErr        bool
		wantLast       string
		wantReason     string
	}{
		{
			name:      "success on the first attempt",
			results:   []func(context.Context) error{returns(nil)},
			wantCalls: 1,
			wantLast:  "telemetry: sent",
		},
		{
			name:      "503 then success",
			results:   []func(context.Context) error{returns(status(503)), returns(nil)},
			wantCalls: 2,
			wantWaits: []time.Duration{time.Second},
			wantLast:  "telemetry: sent",
		},
		{
			name:       "network down on every attempt",
			results:    []func(context.Context) error{returns(errConnRefused)},
			wantCalls:  4,
			wantWaits:  []time.Duration{time.Second, 2 * time.Second, 4 * time.Second},
			wantErr:    true,
			wantLast:   "telemetry: dropped",
			wantReason: "network",
		},
		{
			name:      "dns failure is retried",
			results:   []func(context.Context) error{returns(errDNS), returns(nil)},
			wantCalls: 2,
			wantWaits: []time.Duration{time.Second},
			wantLast:  "telemetry: sent",
		},
		{
			name:       "401 is not retried",
			results:    []func(context.Context) error{returns(status(401))},
			wantCalls:  1,
			wantErr:    true,
			wantLast:   "telemetry: dropped",
			wantReason: "status 401",
		},
		{
			name:       "400 is not retried",
			results:    []func(context.Context) error{returns(status(400))},
			wantCalls:  1,
			wantErr:    true,
			wantLast:   "telemetry: dropped",
			wantReason: "status 400",
		},
		{
			name:           "stuck attempt is aborted and retried",
			results:        []func(context.Context) error{blocksUntilTimeout, returns(nil)},
			attemptTimeout: 50 * time.Millisecond,
			wantCalls:      2,
			wantWaits:      []time.Duration{time.Second},
			wantLast:       "telemetry: sent",
		},
		{
			name:      "429 waits the normal back-off",
			results:   []func(context.Context) error{returns(status(429)), returns(nil)},
			wantCalls: 2,
			wantWaits: []time.Duration{time.Second},
			wantLast:  "telemetry: sent",
		},
		{
			name:       "sender deadline during the back-off",
			results:    []func(context.Context) error{returns(errConnRefused)},
			sleepErr:   context.DeadlineExceeded,
			wantCalls:  1,
			wantWaits:  []time.Duration{time.Second},
			wantErr:    true,
			wantLast:   "telemetry: dropped",
			wantReason: "sender deadline",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var waits []time.Duration
			policy := testPolicy(&waits, tc.sleepErr)
			if tc.attemptTimeout > 0 {
				policy.AttemptTimeout = tc.attemptTimeout
			}
			exporter := &scriptedExporter{results: tc.results}
			logs := &logRecorder{}

			err := policy.Retry(context.Background(), exporter, testEvent(), logs.log)

			if (err != nil) != tc.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if exporter.callCount() != tc.wantCalls {
				t.Errorf("%d calls, want %d", exporter.callCount(), tc.wantCalls)
			}
			if fmt.Sprint(waits) != fmt.Sprint(tc.wantWaits) {
				t.Errorf("waits = %v, want %v", waits, tc.wantWaits)
			}
			last := logs.last()
			if last.msg != tc.wantLast {
				t.Errorf("last log = %q, want %q", last.msg, tc.wantLast)
			}
			if tc.wantReason != "" && last.value("reason") != tc.wantReason {
				t.Errorf("reason = %q, want %q", last.value("reason"), tc.wantReason)
			}
			wantLines := tc.wantCalls
			if tc.sleepErr != nil {
				wantLines++
			}
			if len(logs.entries) != wantLines {
				t.Errorf("%d log lines, want one per attempt plus the final line", len(logs.entries))
			}
			for _, entry := range logs.entries {
				text := fmt.Sprint(entry.msg, entry.kv)
				for _, eventData := range []string{"tenant-secret", "installation-secret", "virtualmachine"} {
					if strings.Contains(text, eventData) {
						t.Errorf("log must not contain event data %q: %s", eventData, text)
					}
				}
			}
		})
	}
}

func TestRetryLogsEveryAttempt(t *testing.T) {
	var waits []time.Duration
	logs := &logRecorder{}
	exporter := &scriptedExporter{results: []func(context.Context) error{returns(errConnRefused)}}

	_ = testPolicy(&waits, nil).Retry(context.Background(), exporter, testEvent(), logs.log)

	for i, entry := range logs.entries[:len(logs.entries)-1] {
		if entry.msg != "telemetry: attempt failed" {
			t.Errorf("line %d = %q, want telemetry: attempt failed", i, entry.msg)
		}
		if entry.value("attempt") != fmt.Sprint(i+1) || entry.value("reason") != "network" {
			t.Errorf("line %d = %+v, want attempt %d with reason network", i, entry, i+1)
		}
		if entry.value("next_wait_ms") != fmt.Sprint(waits[i].Milliseconds()) {
			t.Errorf("line %d next_wait_ms = %s, want %d", i, entry.value("next_wait_ms"), waits[i].Milliseconds())
		}
		if entry.value("duration_ms") == "" {
			t.Errorf("line %d must carry duration_ms", i)
		}
	}
}

func TestBackoffJitter(t *testing.T) {
	testCases := []struct {
		name string
		rand float64
		want []time.Duration
	}{
		{"lowest jitter", 0, []time.Duration{500 * time.Millisecond, time.Second, 2 * time.Second}},
		{"no jitter", 0.5, []time.Duration{time.Second, 2 * time.Second, 4 * time.Second}},
		{"highest jitter", 0.999, []time.Duration{1499 * time.Millisecond, 2998 * time.Millisecond, 5996 * time.Millisecond}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			policy := DefaultPolicy()
			policy.Rand = func() float64 { return tc.rand }
			for i, want := range tc.want {
				if got := policy.backoff(i + 1).Truncate(time.Millisecond); got != want {
					t.Errorf("wait before retry %d = %v, want %v", i+1, got, want)
				}
			}
		})
	}
}

func TestDescribeFailure(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want string
	}{
		{"attempt timeout", context.DeadlineExceeded, "timeout"},
		{"dns", errDNS, "dns"},
		{"connection refused", errConnRefused, "network"},
		{"status", status(http.StatusBadGateway), "status 502"},
		{"serialization", errors.New("json: unsupported value"), "unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := describeFailure(tc.err); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
