package cmd

import (
	"io"
	"reflect"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd/cmdtelemetry"
	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/cli/cmd/telemetrycmd"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func fakeSchema() *core.Schema {
	return (*core.Schema)(&openapi3.Schema{})
}

// fakeConfigFlagValue é uma implementação mínima de schema_flags.SchemaFlagValue, só
// pra ter configValue sem montar um *core.Schema de verdade (Desc().Schema fica
// nil, seguro desde que o cfg de teste sempre ache um valor pros nomes consultados).
type fakeConfigFlagValue struct {
	propName string
	raw      string
	changed  bool
}

var _ schema_flags.SchemaFlagValue = (*fakeConfigFlagValue)(nil)

func (f *fakeConfigFlagValue) String() string { return f.raw }
func (f *fakeConfigFlagValue) Set(v string) error {
	f.raw = v
	f.changed = true
	return nil
}
func (f *fakeConfigFlagValue) Type() string { return "string" }
func (f *fakeConfigFlagValue) Desc() schema_flags.SchemaFlagValueDesc {
	return schema_flags.SchemaFlagValueDesc{PropName: f.propName, IsConfig: true, Schema: fakeSchema()}
}
func (f *fakeConfigFlagValue) RawDefaultValue() string { return "" }
func (f *fakeConfigFlagValue) RawNoOptDefVal() string  { return "" }
func (f *fakeConfigFlagValue) Usage() string           { return "" }
func (f *fakeConfigFlagValue) Parse() (any, error)     { return f.raw, nil }
func (f *fakeConfigFlagValue) Changed() bool           { return f.changed }

func addConfigFlag(cmd *cobra.Command, name string) {
	cmd.Flags().Var(&fakeConfigFlagValue{propName: name}, name, "")
}

// fakeConfig devolve um *config.Config com region/env pré-configurados (como se tivessem
// vindo de cli.yaml), pra testar o fallback de configValue quando a flag não é passada.
func fakeConfig() *config.Config {
	pm, _ := profile_manager.NewInMemoryProfileManager()
	cfg := config.New(pm)
	_ = cfg.SetTempConfig("region", "br-se1-config")
	_ = cfg.SetTempConfig("env", "prod-config")
	return cfg
}

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

	vm := &cobra.Command{Use: "virtual-machine", Aliases: []string{"vm"}, RunE: showHelp}
	instances := &cobra.Command{Use: "instances", RunE: showHelp}
	create := &cobra.Command{Use: "create", RunE: noop}
	addConfigFlag(create, "region")
	addConfigFlag(create, "env")
	create.Flags().String("machine-type", "", "")
	create.Flags().String("name", "", "")

	get := &cobra.Command{Use: "get [id]", Args: cobra.MaximumNArgs(1)}
	get.Flags().String("id", "", "")
	get.RunE = func(c *cobra.Command, args []string) error {
		return applyPositionalArgs([]*flag.Flag{c.Flags().Lookup("id")}, args)
	}

	telemetryCmd := &cobra.Command{Use: telemetrycmd.CommandName}
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
	return cmdtelemetry.BuildCommandInfo(root, args, err, fakeConfig(), getNextUnknownCommand)
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
		wantAction     string
		wantProduct    string
		wantRegion     string
		wantEnv        string
	}{
		{
			name:        "command with flags keeps only flag names",
			args:        []string{"virtual-machine", "instances", "create", "--region", "br-se1", "--machine-type", "BV1-1-10", "--name", "vm1", "--debug"},
			wantTrack:   true,
			wantPath:    []string{"virtual-machine", "instances", "create"},
			wantOptions: []string{"debug", "machine-type", "name", "region"},
			// region vem da flag explícita; env não foi passada, cai no cfg de teste.
			wantRegion: "br-se1",
			wantEnv:    "prod-config",
		},
		{
			name:       "region and env fall back to config when the flags are not set",
			args:       []string{"virtual-machine", "instances", "create"},
			wantTrack:  true,
			wantPath:   []string{"virtual-machine", "instances", "create"},
			wantRegion: "br-se1-config",
			wantEnv:    "prod-config",
		},
		{
			name:        "flag names are normalized",
			args:        []string{"virtual-machine", "instances", "create", "--machineType", "BV1-1-10"},
			wantTrack:   true,
			wantPath:    []string{"virtual-machine", "instances", "create"},
			wantOptions: []string{"machine-type"},
			wantRegion:  "br-se1-config",
			wantEnv:     "prod-config",
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
			wantPath:    []string{"foo"},
			wantAction:  "foo",
			wantProduct: telemetry.UnknownProduct,
		},
		{
			name:        "unknown command inside a group",
			args:        []string{"virtual-machine", "--debug", "foo"},
			wantTrack:   true,
			wantUnknown: true,
			wantPath:    []string{"virtual-machine", "foo"},
			wantOptions: []string{"debug"},
			wantAction:  "virtualmachine.foo",
		},
		{
			name:        "unknown subcommand inside a group resolved through an alias",
			args:        []string{"vm", "list"},
			wantTrack:   true,
			wantUnknown: true,
			wantPath:    []string{"virtual-machine", "list"},
			wantAction:  "virtualmachine.list",
			wantProduct: "virtual machine",
		},
		{
			// Só o primeiro argumento não resolvido entra no Path; o que vem depois é ignorado.
			name:        "extra arguments after the unresolved token are ignored",
			args:        []string{"virtual-machine", "foo", "extra"},
			wantTrack:   true,
			wantUnknown: true,
			wantPath:    []string{"virtual-machine", "foo"},
			wantAction:  "virtualmachine.foo",
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
		{name: "telemetry command", args: []string{telemetrycmd.CommandName, "disable"}},
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
			if info.Region != tc.wantRegion {
				t.Errorf("region = %q, want %q", info.Region, tc.wantRegion)
			}
			if info.Env != tc.wantEnv {
				t.Errorf("env = %q, want %q", info.Env, tc.wantEnv)
			}
			if tc.wantAction != "" {
				if got := info.Action(); got != tc.wantAction {
					t.Errorf("action = %q, want %q", got, tc.wantAction)
				}
			}
			if tc.wantProduct != "" {
				if got := info.Product(); got != tc.wantProduct {
					t.Errorf("product = %q, want %q", got, tc.wantProduct)
				}
			}
		})
	}
}
