package cmd

import (
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
