package cmd

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/MagaluCloud/magalu/mgc/core"
)

var (
	_ telemetry.FailureReasoner = notLoggedInError{}
	_ telemetry.FailureReasoner = missingScopesError{}
)

type notLoggedInError struct{}

func (notLoggedInError) Error() string {
	return "you are not logged in. To authenticate, please run 'mgc auth login'"
}

func (notLoggedInError) TelemetryFailureReason() telemetry.FailureReason {
	return telemetry.FailureAuthentication
}

type missingScopesError struct {
	Missing core.Scopes
}

func (e missingScopesError) Error() string {
	return fmt.Sprintf("you are missing the following scopes for this operation: %v", e.Missing)
}

func (missingScopesError) TelemetryFailureReason() telemetry.FailureReason {
	return telemetry.FailureAuthorization
}
