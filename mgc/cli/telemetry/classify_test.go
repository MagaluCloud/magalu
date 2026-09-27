package telemetry

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/erikgeiser/promptkit"
)

type reasonedError struct{ reason FailureReason }

func (e reasonedError) Error() string                         { return "reasoned" }
func (e reasonedError) TelemetryFailureReason() FailureReason { return e.reason }

type wrappingReasonedError struct {
	reason FailureReason
	err    error
}

func (e wrappingReasonedError) Error() string                         { return "wrapping" }
func (e wrappingReasonedError) Unwrap() error                         { return e.err }
func (e wrappingReasonedError) TelemetryFailureReason() FailureReason { return e.reason }

func apiError(code int, slug string) error {
	return &mgcHttpPkg.IdentifiableHttpError{
		HttpError: &mgcHttpPkg.HttpError{Code: code, Slug: slug, Message: "details with user data"},
		RequestID: "req-1",
	}
}

type timeoutNetError struct{}

func (timeoutNetError) Error() string   { return "i/o timeout" }
func (timeoutNetError) Timeout() bool   { return true }
func (timeoutNetError) Temporary() bool { return true }

func TestClassifyError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want FailureReason
	}{
		{"nil", nil, ""},
		{"confirmation denied", core.UserDeniedConfirmationError{Prompt: "delete?"}, FailureUserCancelled},
		{"prompt aborted", fmt.Errorf("wrap: %w", promptkit.ErrAborted), FailureUserCancelled},
		{"context deadline", fmt.Errorf("op: %w", context.DeadlineExceeded), FailureTimeout},
		{"net timeout", &url.Error{Op: "Get", URL: "https://api", Err: timeoutNetError{}}, FailureTimeout},
		{"401", apiError(401, "unauthorized"), FailureAuthentication},
		{"403", apiError(403, ""), FailureAuthorization},
		{"404", apiError(404, ""), FailureNotFound},
		{"409", apiError(409, ""), FailureConflict},
		{"429", apiError(429, ""), FailureQuota},
		{"400 with quota slug", apiError(400, "quota_exceeded"), FailureQuota},
		{"422 with limit slug", apiError(422, "InstanceLimitReached"), FailureQuota},
		{"400 with limit_exceeded slug", apiError(400, "limit_exceeded"), FailureQuota},
		{"400 with pagination limit slug", apiError(400, "invalid_limit"), FailureValidation},
		{"401 with limit slug", apiError(401, "token_limit_reached"), FailureAuthentication},
		{"503 with quota slug", apiError(503, "quota_service_down"), FailureAPIServer},
		{"400", apiError(400, "invalid_parameter"), FailureValidation},
		{"500", apiError(500, "unknown"), FailureAPIServer},
		{"503 wrapped", fmt.Errorf("while waiting: %w", apiError(503, "")), FailureAPIServer},
		{"dns", &url.Error{Op: "Get", URL: "https://api", Err: &net.DNSError{Err: "no such host", Name: "api"}}, FailureNetwork},
		{"connection refused", &net.OpError{Op: "dial", Err: errors.New("connection refused")}, FailureNetwork},
		{"typed cli error", reasonedError{FailureAuthentication}, FailureAuthentication},
		{"typed cli error with invalid reason", reasonedError{"bogus"}, FailureUnknown},
		{"typed cli error wrapping http error", wrappingReasonedError{FailureAuthentication, apiError(500, "")}, FailureAuthentication},
		{"typed cli error with invalid reason wrapping http error", wrappingReasonedError{"bogus", apiError(404, "")}, FailureNotFound},
		{"usage error", core.UsageError{Err: errors.New("bad param")}, FailureValidation},
		{"untyped", errors.New("something odd"), FailureUnknown},
	}

	for _, c := range cases {
		if got := ClassifyError(c.err); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
