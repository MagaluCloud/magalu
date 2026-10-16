package telemetry

import "context"

type NoopExporter struct{}

func (NoopExporter) Export(context.Context, Event) error {
	return nil
}
