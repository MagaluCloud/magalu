package telemetry

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ExportTimeout = time.Second

	EnvDoNotTrack = "DO_NOT_TRACK"
	EnvOptOut     = "MGC_CLI_TELEMETRY_OPTOUT"

	DisabledByConfig          = "telemetry.disabled"
	DisabledByDoNotTrack      = EnvDoNotTrack
	DisabledByOptOut          = EnvOptOut
	DisabledByUnreadableState = "unreadable telemetry.yaml"
)

type DebugLogger func(msg string, keysAndValues ...any)

type Options struct {
	ClientVersion  string
	OS             string
	ExecutablePath string
	IsTerminal     bool
	Store          StateStore
	Exporter       Exporter
	Getenv         func(string) string
	Now            func() time.Time
	NewID          func() string
	Debug          DebugLogger
	Timeout        time.Duration
}

type Service struct {
	opts  Options
	state State

	executionContext     ExecutionContext
	executionEnvironment string

	stateUnavailable bool
}

func New(opts Options) *Service {
	if opts.Getenv == nil {
		opts.Getenv = os.Getenv
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	if opts.NewID == nil {
		opts.NewID = func() string { return uuid.NewString() }
	}
	if opts.Debug == nil {
		opts.Debug = func(string, ...any) {}
	}
	if opts.Timeout <= 0 {
		opts.Timeout = ExportTimeout
	}
	if opts.Exporter == nil {
		opts.Exporter = NoopExporter{}
	}

	s := &Service{opts: opts}

	if opts.Store != nil {
		state, err := opts.Store.Load()
		if err != nil {
			opts.Debug("telemetry: unable to load state", "error", err)
			s.stateUnavailable = true
		}
		s.state = state
	}

	s.executionContext, s.executionEnvironment = DetectExecutionContext(opts.Getenv, opts.IsTerminal)
	return s
}

func (s *Service) DisabledBy() string {
	switch {
	case isTruthy(s.opts.Getenv(EnvDoNotTrack)):
		return DisabledByDoNotTrack
	case isTruthy(s.opts.Getenv(EnvOptOut)):
		return DisabledByOptOut
	case s.stateUnavailable:
		return DisabledByUnreadableState
	case s.state.Disabled:
		return DisabledByConfig
	}
	return ""
}

func (s *Service) Enabled() bool {
	return s.DisabledBy() == ""
}

func (s *Service) ExecutionContext() ExecutionContext {
	return s.executionContext
}

func (s *Service) SetDisabled(disabled bool) error {
	s.state.Disabled = disabled
	if err := s.save(); err != nil {
		return err
	}
	s.stateUnavailable = false
	return nil
}

func (s *Service) Record(ctx context.Context, info CommandInfo, cmdErr error, w io.Writer) {
	if !s.Enabled() {
		return
	}

	login := info.IsLogin() && cmdErr == nil
	notice := s.noticePending()

	// O login é coletado mesmo quando mostra o aviso, mas qualquer outro comando que
	// mostra o aviso fica de fora e a coleta começa na execução seguinte
	collect := login || !notice

	changed := login && s.markCredentialsSet()
	if notice {
		s.state.NoticeShown = true
		changed = true
	}

	var event Event
	if collect {
		event = s.buildEvent(info, cmdErr)
		changed = s.ensureInstallationID(&event) || changed
		changed = s.applyTimeToFirstValue(&event, info) || changed
	}

	if changed {
		if err := s.save(); err != nil {
			s.opts.Debug("telemetry: unable to save state", "error", err)
			if notice {
				return
			}
		}
	}

	if notice {
		printNotice(w, collect)
	}
	if collect {
		s.export(ctx, event)
	}
}

func (s *Service) export(ctx context.Context, event Event) {
	ctx, cancel := context.WithTimeout(ctx, s.opts.Timeout)
	defer cancel()

	if err := s.opts.Exporter.Export(ctx, event); err != nil {
		s.opts.Debug("telemetry: event dropped", "reason", describeExportError(ctx, err), "error", err.Error())
	}
}

func describeExportError(ctx context.Context, err error) string {
	var (
		dnsErr    *net.DNSError
		statusErr ExportStatusError
	)
	switch {
	case errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded):
		return "timeout"
	case errors.As(err, &dnsErr):
		return "dns"
	case errors.As(err, &statusErr):
		return fmt.Sprintf("status %d", statusErr.StatusCode)
	case isNetworkError(err):
		return "network"
	}
	return "unknown"
}

func (s *Service) buildEvent(info CommandInfo, cmdErr error) Event {
	event := Event{
		Timestamp:            info.End,
		ExecutionContext:     s.executionContext,
		ExecutionEnvironment: s.executionEnvironment,
		Action:               info.Action(),
		Product:              info.Product(),
		OptionsSet:           info.OptionsSet,
		Outcome:              OutcomeSuccess,
		Duration:             info.End.Sub(info.Start),
		ClientVersion:        s.opts.ClientVersion,
		OS:                   s.opts.OS,
		InstallMethod:        DetectInstallMethod(s.opts.ExecutablePath),
		LastRequestID:        info.LastRequestID,
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = s.opts.Now()
	}
	if info.TenantID != "" {
		event.Actor = &Actor{TenantID: info.TenantID}
	}
	if resourceType := info.ResourceType(); resourceType != "" {
		event.Resource = &Resource{Type: resourceType, ID: info.ResourceID}
	}

	if cmdErr != nil {
		event.Outcome = OutcomeFailure
		event.FailureReason = ClassifyError(cmdErr)
		if info.UnknownCommand && event.FailureReason == FailureUnknown {
			event.FailureReason = FailureValidation
		}
	}

	return event
}

func (s *Service) markCredentialsSet() bool {
	if s.state.CredentialsSetAt != nil {
		return false
	}
	now := s.opts.Now().UTC()
	s.state.CredentialsSetAt = &now
	return true
}

func (s *Service) ensureInstallationID(event *Event) bool {
	changed := false
	if s.state.InstallationID == "" {
		s.state.InstallationID = s.opts.NewID()
		changed = true
	}
	event.InstallationID = s.state.InstallationID
	return changed
}

func (s *Service) applyTimeToFirstValue(event *Event, info CommandInfo) bool {
	if event.Outcome != OutcomeSuccess || s.state.FirstValueRecorded || s.state.CredentialsSetAt == nil || !info.IsInfrastructure() {
		return false
	}

	ttfv := event.Timestamp.Sub(*s.state.CredentialsSetAt)
	if ttfv < 0 {
		return false
	}

	event.TimeToFirstValue = &ttfv
	s.state.FirstValueRecorded = true
	return true
}

func (s *Service) save() error {
	if s.opts.Store == nil {
		return nil
	}
	return s.opts.Store.Save(s.state)
}

func isTruthy(v string) bool {
	v = strings.TrimSpace(v)
	return v != "" && v != "0" && !strings.EqualFold(v, "false")
}
