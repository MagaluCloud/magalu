package sender

import (
	"context"
	"sync"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

func (p Policy) Send(ctx context.Context, event telemetry.Event, exporters []telemetry.Exporter, log telemetry.DebugLogger) {
	var wg sync.WaitGroup
	for i, exporter := range exporters {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exporterLog := withExporter(log, i)
			// Um panic aqui derruba o processo de envio com código diferente de 0
			defer func() {
				if recovered := recover(); recovered != nil {
					exporterLog("telemetry: dropped", "reason", "panic")
				}
			}()
			_ = p.Retry(ctx, exporter, event, exporterLog)
		}()
	}
	wg.Wait()
}

func withExporter(log telemetry.DebugLogger, index int) telemetry.DebugLogger {
	return func(msg string, keysAndValues ...any) {
		exporterField := []any{"exporter", index}
		fields := append(exporterField, keysAndValues...)
		log(msg, fields...)
	}
}
