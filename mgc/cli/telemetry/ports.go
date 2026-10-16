package telemetry

import (
	"context"
	"fmt"
	"time"
)

type Exporter interface {
	Export(ctx context.Context, event Event) error
}

// Dispatcher entrega o evento para ser enviado fora do comando, sem esperar a rede
type Dispatcher interface {
	Dispatch(event Event) error
}

type State struct {
	Disabled           bool
	NoticeShown        bool
	CredentialsSetAt   *time.Time
	FirstValueRecorded bool
	InstallationID     string
}

type StateStore interface {
	Load() (State, error)
	Save(State) error
}

const StateFileName = "telemetry.yaml"

type DebugLogger func(msg string, keysAndValues ...any)

type ExportStatusError struct {
	StatusCode int
}

func (e ExportStatusError) Error() string {
	return fmt.Sprintf("unexpected status %d", e.StatusCode)
}
