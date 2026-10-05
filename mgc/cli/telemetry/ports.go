package telemetry

import (
	"context"
	"time"
)

type Exporter interface {
	Export(ctx context.Context, event Event) error
}

// Warmer é opcional. O exportador abre a conexão enquanto o comando roda, para o
// envio no fim não pagar o handshake dentro do orçamento de envio
type Warmer interface {
	Warm(ctx context.Context)
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
