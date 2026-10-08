package dispatch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

// SenderArg é o argumento interno que faz o binário da CLI rodar como processo de envio
const SenderArg = "__telemetry-send"

const payloadVersion = 1

var ErrUnknownPayloadVersion = errors.New("unknown telemetry payload version")

// Payload é o que a CLI deixa no arquivo temporário para o processo de envio.
// InstallationID e Product viajam fora do Event porque não fazem parte do JSON do modelo
type Payload struct {
	Version        int             `json:"version"`
	Event          telemetry.Event `json:"event"`
	InstallationID string          `json:"installationId,omitempty"`
	Product        string          `json:"product,omitempty"`
	DebugLogPath   string          `json:"debugLogPath,omitempty"`
}

func Encode(event telemetry.Event, debugLogPath string) ([]byte, error) {
	return json.Marshal(Payload{
		Version:        payloadVersion,
		Event:          event,
		InstallationID: event.InstallationID,
		Product:        event.Product,
		DebugLogPath:   debugLogPath,
	})
}

func Decode(data []byte) (Payload, error) {
	var payload Payload
	if err := json.Unmarshal(data, &payload); err != nil {
		return Payload{}, err
	}
	if payload.Version != payloadVersion {
		return Payload{}, fmt.Errorf("%w %d", ErrUnknownPayloadVersion, payload.Version)
	}
	payload.Event.InstallationID = payload.InstallationID
	payload.Event.Product = payload.Product
	return payload, nil
}
