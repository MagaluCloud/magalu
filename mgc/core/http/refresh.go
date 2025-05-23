package http

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

type refreshTokenFn func(ctx context.Context) (string, error)

type RefreshLogger struct {
	Transport http.RoundTripper
	RefreshFn refreshTokenFn
}

func NewDefaultRefreshLogger(t http.RoundTripper, rFn refreshTokenFn) *RefreshLogger {
	return &RefreshLogger{
		Transport: t,
		RefreshFn: rFn,
	}
}

// DefaultBackoff provides a default callback for Client.Backoff which
// will perform exponential backoff based on the attempt number and limited
// by the provided minimum and maximum durations.
//
// It also tries to parse Retry-After response header when a http.StatusTooManyRequests
// (HTTP Code 429) is found in the resp parameter. Hence it will return the number of
// seconds the server states it may be ready to process more requests from this client.
func DefaultBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if s, ok := resp.Header["Retry-After"]; ok {
				if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
					return time.Second * time.Duration(sleep)
				}
			}
		}
	}

	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		sleep = max
	}
	return sleep
}

func (t *RefreshLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	resp, err := transport.RoundTrip(req)
	if req.Header.Get("Authorization") == "" {
		return resp, err
	}
	if err != nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	token, rErr := t.RefreshFn(req.Context())
	if rErr != nil {
		return resp, fmt.Errorf("Unauthorized and failed to refresh token. Please, login again: %w", rErr)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("tried to recalculate request body for retrying after token refresh but failed: %w", err)
		}
		req.Body = body
	}

	return transport.RoundTrip(req)
}
