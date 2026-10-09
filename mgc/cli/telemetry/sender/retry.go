package sender

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

const (
	defaultAttempts       = 4
	defaultAttemptTimeout = 10 * time.Second
	defaultBaseDelay      = time.Second

	// jitterFloor faz a espera variar entre 50% e 150% do back-off
	jitterFloor = 0.5
)

type Policy struct {
	Attempts       int
	AttemptTimeout time.Duration
	BaseDelay      time.Duration
	Sleep          func(ctx context.Context, d time.Duration) error
	Rand           func() float64
}

func DefaultPolicy() Policy {
	return Policy{
		Attempts:       defaultAttempts,
		AttemptTimeout: defaultAttemptTimeout,
		BaseDelay:      defaultBaseDelay,
		Sleep:          sleep,
		Rand:           rand.Float64,
	}
}

func (p Policy) Retry(ctx context.Context, exporter telemetry.Exporter, event telemetry.Event, log telemetry.DebugLogger) error {
	var err error
	for attempt := 1; attempt <= p.Attempts; attempt++ {
		start := time.Now()
		attemptCtx, cancel := context.WithTimeout(ctx, p.AttemptTimeout)
		err = exporter.Export(attemptCtx, event)
		cancel()
		elapsed := time.Since(start).Milliseconds()

		if err == nil {
			log("telemetry: sent", "attempt", attempt, "duration_ms", elapsed)
			return nil
		}

		reason := describeFailure(err)
		if attempt == p.Attempts || !retryable(ctx, err) {
			log("telemetry: dropped", "attempt", attempt, "duration_ms", elapsed, "reason", reason)
			return err
		}

		wait := p.backoff(attempt)
		log("telemetry: attempt failed", "attempt", attempt, "duration_ms", elapsed, "reason", reason, "next_wait_ms", wait.Milliseconds())
		if sleepErr := p.Sleep(ctx, wait); sleepErr != nil {
			log("telemetry: dropped", "attempt", attempt, "reason", "sender deadline")
			return err
		}
	}
	return err
}

func (p Policy) backoff(attempt int) time.Duration {
	// deslocar os bits uma posição à esquerda equivale a multiplicar um número inteiro por 2.
	delay := p.BaseDelay << (attempt - 1)
	return time.Duration(float64(delay) * (jitterFloor + p.Rand()))
}

func retryable(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	var statusErr telemetry.ExportStatusError
	if errors.As(err, &statusErr) {
		return statusErr.StatusCode == http.StatusTooManyRequests || statusErr.StatusCode >= http.StatusInternalServerError
	}
	var (
		urlErr *url.Error
		netErr net.Error
	)
	return errors.Is(err, context.DeadlineExceeded) || errors.As(err, &urlErr) || errors.As(err, &netErr)
}

func describeFailure(err error) string {
	var (
		dnsErr    *net.DNSError
		netErr    net.Error
		statusErr telemetry.ExportStatusError
		urlErr    *url.Error
	)
	switch {
	case errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &netErr) && netErr.Timeout()):
		return "timeout"
	case errors.As(err, &dnsErr):
		return "dns"
	case errors.As(err, &statusErr):
		return fmt.Sprintf("status %d", statusErr.StatusCode)
	case errors.As(err, &urlErr) || errors.As(err, &netErr):
		return "network"
	}
	return "unknown"
}

func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
