package cmdtelemetry

import (
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/spf13/cobra"
)

var _ telemetry.FailureReasoner = InvalidFlagError{}

type InvalidFlagError struct {
	Err error
}

func (e InvalidFlagError) Error() string {
	return e.Err.Error()
}

func (e InvalidFlagError) Unwrap() error {
	return e.Err
}

func (InvalidFlagError) TelemetryFailureReason() telemetry.FailureReason {
	return telemetry.FailureValidation
}

func WrapFlagError(_ *cobra.Command, err error) error {
	if err == nil {
		return nil
	}
	return InvalidFlagError{Err: err}
}
