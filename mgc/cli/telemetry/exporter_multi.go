package telemetry

import (
	"context"
	"errors"
)

// MultiExporter envia o mesmo evento para todos os exportadores em paralelo, dentro do mesmo context
type MultiExporter []Exporter

func (m MultiExporter) Export(ctx context.Context, event Event) error {
	results := make(chan error, len(m))
	for _, exporter := range m {
		go func() { results <- exporter.Export(ctx, event) }()
	}

	var errs []error
	for range m {
		select {
		case err := <-results:
			errs = append(errs, err)
		case <-ctx.Done():
			return errors.Join(append(errs, ctx.Err())...)
		}
	}
	return errors.Join(errs...)
}
