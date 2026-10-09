package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry/adapters/debuglog"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry/adapters/dispatch"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry/sender"
)

const (
	senderDeadline = 60 * time.Second

	// senderBackstop encerra o processo mesmo se um socket ignorar o ctx
	senderBackstop = senderDeadline + 5*time.Second
)

func RunTelemetrySender(args []string) {
	backstop := time.AfterFunc(senderBackstop, func() { os.Exit(0) })
	defer backstop.Stop()

	exporters := func(log telemetry.DebugLogger) []telemetry.Exporter {
		return telemetry.LoadBuildConfig(os.Getenv).Exporters(log)
	}
	runTelemetrySender(args, exporters, sender.DefaultPolicy(), senderDeadline)
}

func runTelemetrySender(
	args []string,
	exporters func(telemetry.DebugLogger) []telemetry.Exporter,
	policy sender.Policy,
	deadline time.Duration,
) {
	// O log começa descartando e passa a gravar no arquivo quando o payload diz onde.
	// O defer usa a variável, então um panic depois disso fica registrado no modo debug.
	// Antes do payload não há onde registrar, porque o processo não tem terminal
	log := telemetry.DebugLogger(debuglog.Discard)
	defer func() {
		if recovered := recover(); recovered != nil {
			log("telemetry: dropped",
				"reason", "panic", //kv
				"error", describePanic(recovered),
			)
		}
	}()

	if len(args) != 1 {
		return
	}
	payload, err := dispatch.Take(args[0])
	if err != nil {
		return
	}

	if payload.DebugLogPath != "" {
		log = debuglog.New(payload.DebugLogPath, time.Now, os.Getpid())
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	policy.Send(ctx, payload.Event, exporters(log), log)
}

// describePanic só repete a mensagem de erros do runtime, que não carregam dados do
// evento. Para qualquer outro valor fica só o tipo
func describePanic(recovered any) string {
	if runtimeErr, ok := recovered.(runtime.Error); ok {
		return runtimeErr.Error()
	}
	return fmt.Sprintf("%T", recovered)
}
