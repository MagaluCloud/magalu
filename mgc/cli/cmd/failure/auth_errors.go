package failure

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/MagaluCloud/magalu/mgc/core"
)

var (
	_ telemetry.FailureReasoner = NotLoggedInError{}
	_ telemetry.FailureReasoner = MissingScopesError{}
)

type NotLoggedInError struct{}

func (NotLoggedInError) Error() string {
	return "you are not logged in. To authenticate, please run 'mgc auth login'"
}

func (NotLoggedInError) TelemetryFailureReason() telemetry.FailureReason {
	return telemetry.FailureAuthentication
}

type MissingScopesError struct {
	Missing core.Scopes
}

func (e MissingScopesError) Error() string {
	return fmt.Sprintf("you are missing the following scopes for this operation: %v", e.Missing)
}

func (MissingScopesError) TelemetryFailureReason() telemetry.FailureReason {
	return telemetry.FailureAuthorization
}
