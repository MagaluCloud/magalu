package http

import (
	"bytes"
	"io"
	"net/http"
	"os"
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

func (r *ClientRetryer) cloneRequestBody(req *http.Request) (*bytes.Buffer, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	clonedBody := bytes.NewBuffer(body)
	return clonedBody, nil
}

func (r *ClientRetryer) cloneRequest(req *http.Request) *http.Request {
	body, _ := r.cloneRequestBody(req)
	clonedRequest, err := http.NewRequest(req.Method, req.URL.String(), body)
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
	for i := 0; i < r.attempts; i++ {
		reqCopy := r.cloneRequest(req)
		res, err = r.Transport.RoundTrip(reqCopy)

		if err != nil {
			if os.IsTimeout(err) {
				logger().Debug("Request timeout, retrying...", "attempt", i+1, "")
				time.Sleep(waitBeforeRetry)
				waitBeforeRetry = waitBeforeRetry * 2
				continue
			}
			req.Body.Close()
			return res, err
		}
		if res.StatusCode >= 500 {
			logger().Debug("Server responded with fail, retrying...", "attempt", i+1, "status code", res.StatusCode, "")
			time.Sleep(waitBeforeRetry)
			waitBeforeRetry = waitBeforeRetry * 2
			continue
		}
		req.Body.Close()
		return res, err
	}
	req.Body.Close()
	return res, err
}
