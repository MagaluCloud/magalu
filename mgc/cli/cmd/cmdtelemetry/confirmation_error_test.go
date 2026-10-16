package cmdtelemetry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

func TestConfirmationMismatchErrorKeepsCause(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{name: "retype mismatch message", err: fmt.Errorf("aborted, input didn't match %q", "my-bucket")},
		{name: "generic error", err: errors.New("boom")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wrapped := ConfirmationMismatchError{Err: tc.err}

			if wrapped.Error() != tc.err.Error() {
				t.Errorf("Error() = %q, want %q", wrapped.Error(), tc.err.Error())
			}
			if !errors.Is(error(wrapped), tc.err) {
				t.Errorf("errors.Is(wrapped, %v) = false", tc.err)
			}
			if wrapped.TelemetryFailureReason() != telemetry.FailureUserCancelled {
				t.Errorf("TelemetryFailureReason() = %q, want %q", wrapped.TelemetryFailureReason(), telemetry.FailureUserCancelled)
			}
			if got := telemetry.ClassifyError(wrapped); got != telemetry.FailureUserCancelled {
				t.Errorf("ClassifyError(wrapped) = %q, want %q", got, telemetry.FailureUserCancelled)
			}
		})
	}
}
