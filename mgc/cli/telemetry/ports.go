package telemetry

import (
	"context"
	"time"
)

type Exporter interface {
	Export(ctx context.Context, event Event) error
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
