package cmd

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	mgcLoggerPkg "github.com/MagaluCloud/magalu/mgc/core/logger"
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
		Debug:          func(msg string, kv ...any) { mgcLoggerPkg.New[osArgParser]().Infow(msg, kv...) },
	})
}

func mgcConfigDir(sdk *mgcSdk.Sdk) string {
	return filepath.Dir(sdk.ProfileManager().Current().Dir())
}

type telemetryRun struct {
	start, end     time.Time
	requestIDs     *requestIDRecorder
	loggedInBefore bool
}

func startTelemetryRun(svc *telemetry.Service, sdk *mgcSdk.Sdk, root *cobra.Command, args []string) *telemetryRun {
	run := &telemetryRun{
		requestIDs:     trackRequestIDs(sdk),
		loggedInBefore: hasLoginSession(sdk),
	}
	if likelyRecorded(root, args, readCredentials(sdk, root, run.loggedInBefore, ""), os.Getenv) {
		svc.Warm(context.Background())
	}
	return run
}

// likelyRecorded antecipa, antes do Execute, se o comando vai gerar evento, para só
// então abrir a conexão da telemetria.
func likelyRecorded(root *cobra.Command, args []string, creds credentials, getenv func(string) string) bool {
	if slices.ContainsFunc(args, isHelpOrVersionArg) {
		return false
	}
	info, track := telemetryCommandInfo(root, args, nil)
	if !track {
		return false
	}
	return info.IsLogin() || slices.ContainsFunc(args, isAPIKeyArg) || isAuthenticated(creds, getenv)
}

func isHelpOrVersionArg(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "--version"
}

func isAPIKeyArg(arg string) bool {
	return arg == "--"+apiKeyFlag || strings.HasPrefix(arg, "--"+apiKeyFlag+"=")
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
	run *telemetryRun,
) {
	info, track := telemetryCommandInfo(root, args, cmdErr)
	if !track {
		return
	}

	_ = initLogger(sdk, getLogFilterFlag(root))

	info.Start = run.start
	info.End = run.end
	info.LastRequestID = run.requestIDs.LastRequestID()
	if tenantID, err := sdk.Auth().CurrentTenantID(); err == nil {
		info.TenantID = tenantID
	}
	info.Authenticated = isAuthenticated(readCredentials(sdk, root, run.loggedInBefore, info.TenantID), os.Getenv)

	svc.Record(context.Background(), info, cmdErr, os.Stderr)
}

// hasLoginSession só lê o token salvo, sem refresh nem chamada de rede
func hasLoginSession(sdk *mgcSdk.Sdk) bool {
	tenantID, err := sdk.Auth().CurrentTenantID()
	return err == nil && tenantID != ""
}

type credentials struct {
	loggedInBefore bool
	tenantIDAfter  string
	apiKeyFlag     string
	savedAPIKey    string
	savedKeyID     string
	savedKeySecret string
}

func readCredentials(sdk *mgcSdk.Sdk, root *cobra.Command, loggedInBefore bool, tenantIDAfter string) credentials {
	creds := credentials{
		loggedInBefore: loggedInBefore,
		tenantIDAfter:  tenantIDAfter,
		apiKeyFlag:     getApiKeyFlag(root),
	}
	if apiKey, err := sdk.Auth().ApiKey(context.Background()); err == nil {
		creds.savedAPIKey = apiKey
	}
	creds.savedKeyID, creds.savedKeySecret = sdk.Auth().AccessKeyPair()
	return creds
}

func isAuthenticated(c credentials, getenv func(string) string) bool {
	loggedIn := c.loggedInBefore || c.tenantIDAfter != ""
	apiKey := c.apiKeyFlag != "" || c.savedAPIKey != "" || getenv(apiKeyEnvVar) != ""
	keyPair := (c.savedKeyID != "" && c.savedKeySecret != "") ||
		(getenv(objKeyIDEnvVar) != "" && getenv(objKeySecretEnvVar) != "")
	return loggedIn || apiKey || keyPair
}
