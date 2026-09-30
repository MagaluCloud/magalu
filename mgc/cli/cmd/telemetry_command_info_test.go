package cmd

import (
	"io"
	"reflect"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// newFakeCommandTree monta uma árvore do cobra que imita a da CLI (grupos que
// mostram ajuda, id posicional e nomes de flag normalizados) sem carregar o SDK.
func newFakeCommandTree() *cobra.Command {
	noop := func(*cobra.Command, []string) error { return nil }
	showHelp := func(c *cobra.Command, _ []string) error { return c.Help() }

	root := &cobra.Command{
		Use:           "mgc",
		Version:       "v42.42.42",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE:          showHelp,
	}
	root.SetGlobalNormalizationFunc(normalizeFlagName)
	root.PersistentFlags().Bool("debug", false, "")
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)

	vm := &cobra.Command{Use: "virtual-machine", RunE: showHelp}
	instances := &cobra.Command{Use: "instances", RunE: showHelp}
	create := &cobra.Command{Use: "create", RunE: noop}
	create.Flags().String("region", "", "")
	create.Flags().String("machine-type", "", "")
	create.Flags().String("name", "", "")

	get := &cobra.Command{Use: "get [id]", Args: cobra.MaximumNArgs(1)}
	get.Flags().String("id", "", "")
	get.RunE = func(c *cobra.Command, args []string) error {
		return applyPositionalArgs([]*flag.Flag{c.Flags().Lookup("id")}, args)
	}

	telemetryCmd := &cobra.Command{Use: telemetryCmdName}
	telemetryCmd.AddCommand(&cobra.Command{Use: "disable", RunE: noop})

	// Na CLI real, esse erro vem do --help de uma flag e dos links.
	flagHelp := &cobra.Command{Use: "flag-help", RunE: func(*cobra.Command, []string) error {
		return schema_flags.ErrWantHelp
	}}

	instances.AddCommand(create, get)
	vm.AddCommand(instances)
	root.AddCommand(vm, telemetryCmd, flagHelp, &cobra.Command{Use: "dump-tree", RunE: noop})
	return root
}

func runFakeCommandTree(args ...string) (telemetry.CommandInfo, bool) {
	root := newFakeCommandTree()
	root.SetArgs(args)
	err := root.Execute()
	return telemetryCommandInfo(root, args, err)
}

func TestTelemetryCommandInfo(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		wantTrack      bool
		wantUnknown    bool
		wantPath       []string
		wantOptions    []string
		wantResourceID string
	}{
		{
			name:        "command with flags keeps only flag names",
			args:        []string{"virtual-machine", "instances", "create", "--region", "br-se1", "--machine-type", "BV1-1-10", "--name", "vm1", "--debug"},
			wantTrack:   true,
			wantPath:    []string{"virtual-machine", "instances", "create"},
			wantOptions: []string{"debug", "machine-type", "name", "region"},
		},
		{
			name:        "flag names are normalized",
			args:        []string{"virtual-machine", "instances", "create", "--machineType", "BV1-1-10"},
			wantTrack:   true,
			wantPath:    []string{"virtual-machine", "instances", "create"},
			wantOptions: []string{"machine-type"},
		},
		{
			name:           "resource id by flag",
			args:           []string{"virtual-machine", "instances", "get", "--id", "abc"},
			wantTrack:      true,
			wantPath:       []string{"virtual-machine", "instances", "get"},
			wantOptions:    []string{"id"},
			wantResourceID: "abc",
		},
		{
			// O id posicional não marca a flag como alterada, por isso não aparece em
			// wantOptions, mas ainda deve ser capturado.
			name:           "resource id as positional argument",
			args:           []string{"virtual-machine", "instances", "get", "abc"},
			wantTrack:      true,
			wantPath:       []string{"virtual-machine", "instances", "get"},
			wantResourceID: "abc",
		},
		{
			name:        "unknown command at root",
			args:        []string{"foo", "bar"},
			wantTrack:   true,
			wantUnknown: true,
		},
		{
			name:        "unknown command inside a group",
			args:        []string{"virtual-machine", "--debug", "foo"},
			wantTrack:   true,
			wantUnknown: true,
			wantOptions: []string{"debug"},
		},
		{name: "help requested through a flag", args: []string{"flag-help"}},
		{name: "group help flag", args: []string{"virtual-machine", "--help"}},
		{name: "command short help flag", args: []string{"virtual-machine", "instances", "create", "-h"}},
		{name: "version", args: []string{"--version"}},
		{name: "root without arguments", args: []string{}},
		{name: "group without subcommand", args: []string{"virtual-machine"}},
		{name: "root with only a flag", args: []string{"--debug"}},
		{name: "help command", args: []string{"help"}},
		{name: "completion", args: []string{"completion", "bash"}},
		{name: "shell completion request", args: []string{cobra.ShellCompRequestCmd, "virtual-machine", ""}},
		{name: "shell completion request without descriptions", args: []string{cobra.ShellCompNoDescRequestCmd, "virtual-machine", ""}},
		{name: "dump-tree", args: []string{"dump-tree"}},
		{name: "telemetry command", args: []string{telemetryCmdName, "disable"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, track := runFakeCommandTree(tc.args...)

			if track != tc.wantTrack || info.UnknownCommand != tc.wantUnknown {
				t.Fatalf("track=%v unknown=%v, want track=%v unknown=%v", track, info.UnknownCommand, tc.wantTrack, tc.wantUnknown)
			}
			if !reflect.DeepEqual(info.Path, tc.wantPath) {
				t.Errorf("path = %v, want %v", info.Path, tc.wantPath)
			}
			if !reflect.DeepEqual(info.OptionsSet, tc.wantOptions) {
				t.Errorf("options = %v, want %v", info.OptionsSet, tc.wantOptions)
			}
			if info.ResourceID != tc.wantResourceID {
				t.Errorf("resource id = %q, want %q", info.ResourceID, tc.wantResourceID)
			}
		})
	}
}
