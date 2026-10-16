package cmdtelemetry

import (
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

var _ telemetry.FailureReasoner = ConfirmationMismatchError{}

// Classifica a confirmação por retype negada como user_cancelled, igual à confirmação sim/não
type ConfirmationMismatchError struct {
	Err error
}

func (e ConfirmationMismatchError) Error() string {
	return e.Err.Error()
}

func (e ConfirmationMismatchError) Unwrap() error {
	return e.Err
}

func (ConfirmationMismatchError) TelemetryFailureReason() telemetry.FailureReason {
	return telemetry.FailureUserCancelled
}
