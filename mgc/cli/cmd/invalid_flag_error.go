package cmd

import (
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/spf13/cobra"
)

var _ telemetry.FailureReasoner = invalidFlagError{}

type invalidFlagError struct {
	Err error
}

func (e invalidFlagError) Error() string {
	return e.Err.Error()
}

func (e invalidFlagError) Unwrap() error {
	return e.Err
}

func (invalidFlagError) TelemetryFailureReason() telemetry.FailureReason {
	return telemetry.FailureValidation
}

func wrapFlagError(_ *cobra.Command, err error) error {
	if err == nil {
		return nil
	}
	return invalidFlagError{Err: err}
}
