package cmd

import (
	"fmt"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/MagaluCloud/magalu/mgc/core"
)

func TestCheckScopesErrors(t *testing.T) {
	testCases := []struct {
		name    string
		err     error
		message string
		reason  telemetry.FailureReason
	}{
		{
			name:    "not logged in",
			err:     notLoggedInError{},
			message: "you are not logged in. To authenticate, please run 'mgc auth login'",
			reason:  telemetry.FailureAuthentication,
		},
		{
			name:    "missing scopes",
			err:     missingScopesError{Missing: core.Scopes{"virtual-machine.read", "network.read"}},
			message: "you are missing the following scopes for this operation: [virtual-machine.read network.read]",
			reason:  telemetry.FailureAuthorization,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.err.Error(); got != tc.message {
				t.Errorf("message changed: got %q, want %q", got, tc.message)
			}
			if got := telemetry.ClassifyError(tc.err); got != tc.reason {
				t.Errorf("ClassifyError = %q, want %q", got, tc.reason)
			}
			if got := telemetry.ClassifyError(fmt.Errorf("wrapped: %w", tc.err)); got != tc.reason {
				t.Errorf("ClassifyError(wrapped) = %q, want %q", got, tc.reason)
			}
		})
	}
}
