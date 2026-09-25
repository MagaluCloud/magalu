package telemetry

import (
	"encoding/json"
	"slices"
	"strings"
	"time"
)

type ExecutionContext string

const (
	ExecutionContextInteractive ExecutionContext = "interactive"
	ExecutionContextCI          ExecutionContext = "ci"
	ExecutionContextAgent       ExecutionContext = "agent"
)

type Outcome string

const (
	OutcomeSuccess Outcome = "success"
	OutcomeFailure Outcome = "failure"
)

// FailureReason is a closed list of error categories. No other value may be emitted.
type FailureReason string

const (
	FailureValidation     FailureReason = "validation"
	FailureAuthentication FailureReason = "authentication"
	FailureAuthorization  FailureReason = "authorization"
	FailureNotFound       FailureReason = "not_found"
	FailureConflict       FailureReason = "conflict"
	FailureQuota          FailureReason = "quota"
	FailureAPIServer      FailureReason = "api_server"
	FailureNetwork        FailureReason = "network"
	FailureTimeout        FailureReason = "timeout"
	FailureUserCancelled  FailureReason = "user_cancelled"
	FailureUnknown        FailureReason = "unknown"
)

var failureReasons = map[FailureReason]struct{}{
	FailureValidation:     {},
	FailureAuthentication: {},
	FailureAuthorization:  {},
	FailureNotFound:       {},
	FailureConflict:       {},
	FailureQuota:          {},
	FailureAPIServer:      {},
	FailureNetwork:        {},
	FailureTimeout:        {},
	FailureUserCancelled:  {},
	FailureUnknown:        {},
}

// Valid reports whether r belongs to the closed list of failure reasons.
func (r FailureReason) Valid() bool {
	_, ok := failureReasons[r]
	return ok
}

type LogLevel string

const (
	LogLevelInfo LogLevel = "INFO"
	LogLevelWarn LogLevel = "WARN"
)

type Impact string

const (
	ImpactLow    Impact = "LOW"
	ImpactMedium Impact = "MEDIUM"
	ImpactHigh   Impact = "HIGH"
)

const timestampLayout = "2006-01-02T15:04:05.000Z07:00"

type Actor struct {
	TenantID string `json:"tenant_id"`
}

type Resource struct {
	Type string `json:"type"`
	ID   string `json:"id,omitempty"`
}

// Event is the anonymized record of a single command execution.
// It must never carry flag values, raw error messages, credentials or local paths.
type Event struct {
	Timestamp            time.Time
	Actor                *Actor
	ExecutionContext     ExecutionContext
	ExecutionEnvironment string
	Resource             *Resource
	Action               string
	OptionsSet           []string
	Outcome              Outcome
	FailureReason        FailureReason
	Duration             time.Duration
	ClientVersion        string
	OS                   string
	InstallMethod        string
	TimeToFirstValue     *time.Duration
	LastRequestID        string

	// InstallationID is the anonymous distinct id. It is handed to exporters
	// separately and is not part of the serialized properties.
	InstallationID string
}

// LogLevel is INFO on success and WARN on failure.
func (e Event) LogLevel() LogLevel {
	if e.Outcome == OutcomeFailure {
		return LogLevelWarn
	}
	return LogLevelInfo
}

// Impact is derived from the outcome and, on failure, from the failure reason.
func (e Event) Impact() Impact {
	if e.Outcome != OutcomeFailure {
		return ImpactLow
	}

	switch e.normalizedFailureReason() {
	case FailureValidation, FailureNotFound, FailureUserCancelled:
		return ImpactLow
	case FailureAuthentication, FailureAuthorization, FailureConflict, FailureQuota:
		return ImpactMedium
	default:
		return ImpactHigh
	}
}

func (e Event) normalizedFailureReason() FailureReason {
	if !e.FailureReason.Valid() {
		return FailureUnknown
	}
	return e.FailureReason
}

type eventJSON struct {
	Timestamp            string           `json:"timestamp"`
	Actor                *Actor           `json:"actor,omitempty"`
	ExecutionContext     ExecutionContext `json:"executionContext"`
	ExecutionEnvironment string           `json:"executionEnvironment,omitempty"`
	Resource             *Resource        `json:"resource,omitempty"`
	Action               string           `json:"action"`
	OptionsSet           string           `json:"optionsSet,omitempty"`
	Outcome              Outcome          `json:"outcome"`
	FailureReason        FailureReason    `json:"failure_reason,omitempty"`
	DurationMs           int64            `json:"durationMs"`
	LogLevel             LogLevel         `json:"logLevel"`
	Impact               Impact           `json:"impact"`
	ClientVersion        string           `json:"clientVersion"`
	OS                   string           `json:"os"`
	InstallMethod        string           `json:"installMethod"`
	TimeToFirstValueMs   *int64           `json:"timeToFirstValueMs,omitempty"`
	LastRequestID        string           `json:"lastRequestId,omitempty"`
}

func (e Event) MarshalJSON() ([]byte, error) {
	out := eventJSON{
		Timestamp:            e.Timestamp.UTC().Format(timestampLayout),
		ExecutionContext:     e.ExecutionContext,
		ExecutionEnvironment: e.ExecutionEnvironment,
		Action:               e.Action,
		Outcome:              e.Outcome,
		DurationMs:           nonNegativeMs(e.Duration),
		LogLevel:             e.LogLevel(),
		Impact:               e.Impact(),
		ClientVersion:        e.ClientVersion,
		OS:                   e.OS,
		InstallMethod:        e.InstallMethod,
		LastRequestID:        e.LastRequestID,
	}

	if out.InstallMethod == "" {
		out.InstallMethod = "unknown"
	}
	if e.Actor != nil && e.Actor.TenantID != "" {
		out.Actor = e.Actor
	}
	if e.Resource != nil && e.Resource.Type != "" {
		out.Resource = e.Resource
	}
	if len(e.OptionsSet) > 0 {
		out.OptionsSet = joinOptions(e.OptionsSet)
	}
	if e.Outcome == OutcomeFailure {
		out.FailureReason = e.normalizedFailureReason()
	}
	if e.TimeToFirstValue != nil {
		ms := nonNegativeMs(*e.TimeToFirstValue)
		out.TimeToFirstValueMs = &ms
	}

	return json.Marshal(out)
}

func nonNegativeMs(d time.Duration) int64 {
	if d < 0 {
		return 0
	}
	return d.Milliseconds()
}

// joinOptions returns option names sorted alphabetically, without duplicates, comma separated.
func joinOptions(names []string) string {
	sorted := slices.Clone(names)
	slices.Sort(sorted)
	sorted = slices.Compact(sorted)
	if len(sorted) > 0 && sorted[0] == "" {
		sorted = sorted[1:]
	}
	return strings.Join(sorted, ",")
}
