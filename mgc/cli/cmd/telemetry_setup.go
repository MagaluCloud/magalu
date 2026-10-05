package cmd

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newTelemetryService(sdk *mgcSdk.Sdk, version string) *telemetry.Service {
	return telemetry.New(telemetry.Options{
		ClientVersion:  version,
		OS:             runtime.GOOS,
		ExecutablePath: telemetry.ResolveExecutablePath(),
		IsTerminal:     term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd())),
		Store:          telemetry.NewViperStateStore(mgcConfigDir(sdk)),
		Exporter:       telemetry.LoadBuildConfig(os.Getenv).Exporter(),
		Debug:          func(msg string, kv ...any) { logger().Debugw(msg, kv...) },
	})
}

func mgcConfigDir(sdk *mgcSdk.Sdk) string {
	return filepath.Dir(sdk.ProfileManager().Current().Dir())
}

// trackRequestIDs coloca o requestIDRecorder na frente do transport usado pelos
// comandos, para saber o X-Request-Id da última chamada feita na API.
func trackRequestIDs(sdk *mgcSdk.Sdk) *requestIDRecorder {
	client := sdk.HttpClient()
	recorder := newRequestIDRecorder(client.Transport)
	client.Transport = recorder
	return recorder
}

// recordTelemetry completa o CommandInfo com o que só se sabe depois do comando
// e deixa o serviço decidir entre aviso, envio ou nada.
func recordTelemetry(
	svc *telemetry.Service,
	sdk *mgcSdk.Sdk,
	root *cobra.Command,
	args []string,
	cmdErr error,
	start, end time.Time,
	requestIDs *requestIDRecorder,
) {
	info, track := telemetryCommandInfo(root, args, cmdErr)
	if !track {
		return
	}

	info.Start = start
	info.End = end
	info.LastRequestID = requestIDs.LastRequestID()
	if tenantID, err := sdk.Auth().CurrentTenantID(); err == nil {
		info.TenantID = tenantID
	}

	svc.Record(context.Background(), info, cmdErr, os.Stderr)
}
