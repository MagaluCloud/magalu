package common

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newTestAuth(t *testing.T, accessKeyId, secretAccessKey string) *auth.Auth {
	t.Helper()
	pm, _ := profile_manager.NewInMemoryProfileManager()
	cfg := config.New(pm)
	a := auth.New(map[string]auth.Config{}, &http.Client{}, pm, cfg)
	if accessKeyId != "" || secretAccessKey != "" {
		if err := a.SetAccessKey(accessKeyId, secretAccessKey); err != nil {
			t.Fatalf("SetAccessKey: %v", err)
		}
	}
	return a
}

func TestSendRequestWithIgnoredHeaders(t *testing.T) {
	errTransport := errors.New("sentinel transport error")

	testCases := []struct {
		name            string
		accessKeyId     string
		secretAccessKey string
		wantAuthErr     bool
	}{
		{"both empty", "", "", true},
		{"only id set", "id", "", true},
		{"only secret set", "", "secret", true},
		{"key pair present", "id", "secret", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testAuth := newTestAuth(t, tc.accessKeyId, tc.secretAccessKey)

			ctx := auth.NewContext(context.Background(), testAuth)
			ctx = mgcHttpPkg.NewClientContext(ctx, mgcHttpPkg.NewClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return nil, errTransport
			})))

			req, err := http.NewRequest(http.MethodGet, "https://example.com/bucket", nil)
			if err != nil {
				t.Fatalf("NewRequest: %v", err)
			}

			_, err = SendRequestWithIgnoredHeaders(ctx, req, Config{Region: "br-se1"}, excludedHeaders)

			var authErr core.AuthenticationError
			isAuthErr := errors.As(err, &authErr)
			if isAuthErr != tc.wantAuthErr {
				t.Fatalf("errors.As(err, &core.AuthenticationError{}) = %v, want %v (err: %v)", isAuthErr, tc.wantAuthErr, err)
			}

			if tc.wantAuthErr {
				if !strings.Contains(err.Error(), "api-key not set") {
					t.Fatalf("error message = %q, want it to contain %q", err.Error(), "api-key not set")
				}
				return
			}

			// Key pair present: the function must not return at the access-key
			// check, so the signed request has to actually reach the transport.
			if !errors.Is(err, errTransport) {
				t.Fatalf("expected the request to reach the HTTP transport and fail with the sentinel error, got: %v", err)
			}
		})
	}
}
