package cmd

import (
	"fmt"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/MagaluCloud/magalu/mgc/core"
)

func TestCheckScopesErrors(t *testing.T) {
	cases := []struct {
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

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.err.Error(); got != c.message {
				t.Errorf("message changed: got %q, want %q", got, c.message)
			}
			if got := telemetry.ClassifyError(c.err); got != c.reason {
				t.Errorf("ClassifyError = %q, want %q", got, c.reason)
			}
			if got := telemetry.ClassifyError(fmt.Errorf("wrapped: %w", c.err)); got != c.reason {
				t.Errorf("ClassifyError(wrapped) = %q, want %q", got, c.reason)
			}
		})
	}
}
