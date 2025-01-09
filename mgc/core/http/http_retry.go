package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"syscall"
	"time"
)

type ClientRetryer struct {
	Transport http.RoundTripper
	attempts  int
}

func NewDefaultClientRetryer(transport http.RoundTripper) *ClientRetryer {
	return &ClientRetryer{
		Transport: transport,
		attempts:  5,
	}
}

func NewClientRetryerWithAttempts(transport http.RoundTripper, attempts int) *ClientRetryer {
	if attempts <= 0 {
		return NewDefaultClientRetryer(transport)
	}
	return &ClientRetryer{
		Transport: transport,
		attempts:  attempts,
	}
}

func (r *ClientRetryer) cloneRequestBody(req *http.Request) (io.Reader, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger().Error(err)
		return nil, nil
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))
	return bytes.NewReader(body), nil
}

func (r *ClientRetryer) cloneRequest(req *http.Request) *http.Request {
	var body io.Reader
	var err error

	if req.Body != nil {
		body, err = r.cloneRequestBody(req)
		if err != nil {
			logger().Error(err)
			return req
		}
	}
	clonedRequest, err := http.NewRequestWithContext(req.Context(), req.Method, req.URL.String(), body)
	if err != nil {
		return req
	}
	clonedRequest.Header = req.Header
	return clonedRequest
}

func (r *ClientRetryer) RoundTrip(req *http.Request) (*http.Response, error) {
	waitBeforeRetry := 100 * time.Millisecond
	var res *http.Response
	var err error
	if req.Body != nil {
		defer req.Body.Close()
	}
	for i := 0; i < r.attempts; i++ {
		reqCopy := r.cloneRequest(req)
		res, err = r.Transport.RoundTrip(reqCopy)

		if err != nil {
			var sysErr *os.SyscallError

			if os.IsTimeout(err) {
				logger().Infow("\n\n\nRequest timeout, retrying...\n\n\n", "attempt", i+1, "")
				time.Sleep(waitBeforeRetry)
				waitBeforeRetry = waitBeforeRetry * 2
				continue
			}

			if errors.As(err, &sysErr) {
				if sysErr.Err == syscall.ECONNRESET {
					logger().Infow("\n\n\nConn reset by peer! THIS IS A SERVER PROBLEM!!!\n\n\n", "attempt", i+1, "")
					time.Sleep(waitBeforeRetry)
					waitBeforeRetry = waitBeforeRetry * 2
					continue
				}
			}
			return res, err
		}
		if res.StatusCode >= 500 {
			logger().Infow("\n\n\nServer responded with fail, retrying...\n\n\n", "attempt", i+1, "status code", res.StatusCode, "")
			time.Sleep(waitBeforeRetry)
			waitBeforeRetry = waitBeforeRetry * 2
			continue
		}

		return res, err
	}

	return res, err
}
