package cmd

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type failingTransport struct{ err error }

func (f failingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, f.err
}

func TestRequestIDRecorder(t *testing.T) {
	testCases := []struct {
		name          string
		responseIDs   []string
		requestHeader string
		transportErr  error
		want          string
	}{
		{"keeps the id of the last response", []string{"resp-1", "resp-2", "resp-3"}, "", nil, "resp-3"},
		{"falls back to the request header", []string{""}, "req-42", nil, "req-42"},
		{"response header wins over the request header", []string{"resp-1"}, "req-42", nil, "resp-1"},
		{"keeps the request id when the call fails", []string{""}, "req-42", errors.New("connection refused"), "req-42"},
		{"no calls yield no id", nil, "", nil, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				// A primeira chamada responde com responseIDs[0], a segunda com
				// responseIDs[1] e assim por diante.
				if id := tc.responseIDs[calls.Add(1)-1]; id != "" {
					w.Header().Set(requestIDHeader, id)
				}
			}))
			defer server.Close()

			var next http.RoundTripper
			if tc.transportErr != nil {
				next = failingTransport{tc.transportErr}
			}
			recorder := newRequestIDRecorder(next)

			recorder.tracked = func(*http.Request) bool { return true }
			client := &http.Client{Transport: recorder}

			for range tc.responseIDs {
				req, err := http.NewRequest(http.MethodGet, server.URL, nil)
				if err != nil {
					t.Fatal(err)
				}
				if tc.requestHeader != "" {
					req.Header.Set(requestIDHeader, tc.requestHeader)
				}

				resp, err := client.Do(req)
				if (err != nil) != (tc.transportErr != nil) {
					t.Fatalf("err = %v, want transport error %v", err, tc.transportErr)
				}
				if resp != nil {
					_ = resp.Body.Close()
				}
			}

			if got := recorder.LastRequestID(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// hostIDTransport responde sem rede, com o próprio host como X-Request-Id
type hostIDTransport struct{}

func (hostIDTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Set(requestIDHeader, req.URL.Hostname())
	return &http.Response{StatusCode: http.StatusOK, Header: header, Body: http.NoBody, Request: req}, nil
}

func TestRequestIDRecorderOnlyTracksCommandRequests(t *testing.T) {
	commandCtx := withCommandRequests(context.Background())
	versionCheck := request{"https://github.com/MagaluCloud/mgccli/releases/latest", context.Background()}

	testCases := []struct {
		name     string
		requests []request
		want     string
	}{
		{"command call", []request{{"https://api.magalu.cloud/br-se1/compute/v1/instances", commandCtx}}, "api.magalu.cloud"},
		{"version check after the command call is ignored", []request{{"https://api.magalu.cloud/br-se1/compute/v1/instances", commandCtx}, versionCheck}, "api.magalu.cloud"},
		{"version check before the command call", []request{versionCheck, {"https://api.pre-prod.example.com/compute", commandCtx}}, "api.pre-prod.example.com"},
		{"any host counts when it comes from the command", []request{{"https://storage.new-domain.example/bucket", commandCtx}}, "storage.new-domain.example"},
		{"only the version check", []request{versionCheck}, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := newRequestIDRecorder(hostIDTransport{})
			client := &http.Client{Transport: recorder}

			for _, r := range tc.requests {
				req, err := http.NewRequestWithContext(r.ctx, http.MethodGet, r.url, nil)
				if err != nil {
					t.Fatal(err)
				}
				resp, err := client.Do(req)
				if err != nil {
					t.Fatal(err)
				}
				_ = resp.Body.Close()
			}

			if got := recorder.LastRequestID(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

type request struct {
	url string
	ctx context.Context
}
