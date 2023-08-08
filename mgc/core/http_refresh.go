package core

import (
	"fmt"
	"net/http"
)

type refreshTokenFn func() (string, error)

type HttpRefreshLogger struct {
	Transport http.RoundTripper
	RefreshFn refreshTokenFn
}

func NewDefaultHttpRefreshLogger(t http.RoundTripper, rFn refreshTokenFn) *HttpRefreshLogger {
	return &HttpRefreshLogger{
		Transport: t,
		RefreshFn: rFn,
	}
}

func (t *HttpRefreshLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	resp, err := transport.RoundTrip(req)
	if resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	token, rErr := t.RefreshFn()
	if rErr != nil {
		return resp, fmt.Errorf("Unauthorized and failed to refresh token. Please, login again: %w", rErr)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	return transport.RoundTrip(req)
}
