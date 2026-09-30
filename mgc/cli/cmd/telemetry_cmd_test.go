package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

func TestTelemetryCmd(t *testing.T) {
	testCases := []struct {
		name      string
		env       map[string]string
		stateFile string // conteúdo inicial do telemetry.yaml, "" para nenhum
		// brokenDir aponta o store para um caminho dentro de um arquivo, e aí tanto o
		// Load quanto o Save falham, sem depender de permissões que o root ignoraria.
		brokenDir bool
		commands  [][]string
		want      string // saída do último comando
		wantErr   string // trecho esperado no erro do último comando
	}{
		{
			name:     "status when enabled",
			commands: [][]string{{"status"}},
			want:     "Telemetry is enabled.\n",
		},
		{
			name:     "disable",
			commands: [][]string{{"disable"}},
			want:     "Telemetry disabled.\n",
		},
		{
			name:     "status after disable",
			commands: [][]string{{"disable"}, {"status"}},
			want:     "Telemetry is disabled by the CLI settings (mgc telemetry disable).\n",
		},
		{
			name:     "enable after disable",
			commands: [][]string{{"disable"}, {"enable"}},
			want:     "Telemetry enabled.\n",
		},
		{
			name:     "status with DO_NOT_TRACK",
			env:      map[string]string{"DO_NOT_TRACK": "1"},
			commands: [][]string{{"status"}},
			want:     "Telemetry is disabled by the DO_NOT_TRACK environment variable.\n",
		},
		{
			name:     "status with MGC_CLI_TELEMETRY_OPTOUT",
			env:      map[string]string{"MGC_CLI_TELEMETRY_OPTOUT": "1"},
			commands: [][]string{{"status"}},
			want:     "Telemetry is disabled by the MGC_CLI_TELEMETRY_OPTOUT environment variable.\n",
		},
		{
			name:     "enable while DO_NOT_TRACK is set",
			env:      map[string]string{"DO_NOT_TRACK": "1"},
			commands: [][]string{{"enable"}},
			want:     "Telemetry enabled in settings, but still disabled by the DO_NOT_TRACK environment variable.\n",
		},
		{
			name:      "status with an unreadable state file",
			stateFile: "telemetry: [not: valid",
			commands:  [][]string{{"status"}},
			want:      "Telemetry is disabled by an unreadable telemetry.yaml file (run mgc telemetry enable or disable to rewrite it).\n",
		},
		{
			name:      "enable rewrites an unreadable state file",
			stateFile: "telemetry: [not: valid",
			commands:  [][]string{{"enable"}, {"status"}},
			want:      "Telemetry is enabled.\n",
		},
		{
			name:      "disable when the state can be neither read nor saved",
			brokenDir: true,
			commands:  [][]string{{"disable"}},
			wantErr:   "unable to save telemetry settings",
		},
		{
			name:      "enable when the state can be neither read nor saved",
			brokenDir: true,
			commands:  [][]string{{"enable"}},
			wantErr:   "unable to save telemetry settings",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := newTelemetryStateDir(t, tc.stateFile, tc.brokenDir)

			// Cada comando usa um Service novo, como execuções separadas da CLI.
			getService := func() *telemetry.Service {
				return telemetry.New(telemetry.Options{
					Store:  telemetry.NewViperStateStore(dir),
					Getenv: func(k string) string { return tc.env[k] },
				})
			}

			out, err := runTelemetryCommands(t, getService, tc.commands)

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Errorf("err = %v, want it to contain %q", err, tc.wantErr)
			}
			if out != tc.want {
				t.Errorf("got %q, want %q", out, tc.want)
			}
		})
	}
}

func newTelemetryStateDir(t *testing.T, stateFile string, broken bool) string {
	t.Helper()
	dir := t.TempDir()

	if broken {
		blocker := filepath.Join(dir, "file")
		if err := os.WriteFile(blocker, nil, 0o644); err != nil {
			t.Fatal(err)
		}
		return filepath.Join(blocker, "sub")
	}

	if stateFile != "" {
		if err := os.WriteFile(filepath.Join(dir, telemetry.StateFileName), []byte(stateFile), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// runTelemetryCommands roda os comandos em sequência e devolve a saída e o erro do
// último. Um erro num comando intermediário falha o teste.
func runTelemetryCommands(t *testing.T, getService func() *telemetry.Service, commands [][]string) (string, error) {
	t.Helper()

	var out bytes.Buffer
	var err error
	for i, args := range commands {
		out.Reset()
		cmd := newTelemetryCmd(getService)
		cmd.SetArgs(args)
		cmd.SetOut(&out)
		cmd.SetErr(io.Discard)
		cmd.SilenceUsage = true
		if err = cmd.Execute(); err != nil && i < len(commands)-1 {
			t.Fatalf("%v: %v", args, err)
		}
	}
	return out.String(), err
}
