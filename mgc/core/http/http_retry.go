package http

import (
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

func (r *ClientRetryer) RoundTrip(req *http.Request) (*http.Response, error) {
	waitBeforeRetry := 100 * time.Millisecond
	var res *http.Response
	var err error
	for i := 0; i < r.attempts; i++ {
		res, err = r.Transport.RoundTrip(req)
		if err != nil {
			if os.IsTimeout(err) {
				logger().Debug("Request timeout, retrying...", "attempt", i+1, "")
				time.Sleep(waitBeforeRetry)
				waitBeforeRetry = waitBeforeRetry * 2
				continue
			}
			return res, err
		}
		if res.StatusCode >= 400 && res.StatusCode < 600 {
			logger().Debug("Server responded with fail, retrying...", "attempt", i+1, "status code", res.StatusCode, "")
			time.Sleep(waitBeforeRetry)
			waitBeforeRetry = waitBeforeRetry * 2
			continue
		}
		return res, err
	}
	return res, err
}
