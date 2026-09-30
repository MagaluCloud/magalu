package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
	"golang.org/x/term"
)

func newTelemetryService(sdk *mgcSdk.Sdk, version string) *telemetry.Service {
	return telemetry.New(telemetry.Options{
		ClientVersion:  version,
		OS:             runtime.GOOS,
		ExecutablePath: telemetry.ResolveExecutablePath(),
		IsTerminal:     term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd())),
		Store:          telemetry.NewViperStateStore(mgcConfigDir(sdk)),
		Debug:          func(msg string, kv ...any) { logger().Debugw(msg, kv...) },
	})
}

func mgcConfigDir(sdk *mgcSdk.Sdk) string {
	return filepath.Dir(sdk.ProfileManager().Current().Dir())
}
