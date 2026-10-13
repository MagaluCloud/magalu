package cmdtelemetry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/cli/cmd/telemetrycmd"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// NextUnknownFunc espelha getNextUnknownCommand, privada do cmd, injetada aqui para o
// cmdtelemetry não precisar importar o pacote cmd
type NextUnknownFunc func(*cobra.Command, []string) (*string, []string)

var untrackedCommands = map[string]struct{}{
	"help":                          {},
	"completion":                    {},
	cobra.ShellCompRequestCmd:       {},
	cobra.ShellCompNoDescRequestCmd: {},
	"dump-tree":                     {},
	telemetrycmd.CommandName:        {},
}

// BuildCommandInfo descobre qual comando foi executado e monta o CommandInfo.
// O segundo retorno é false quando a execução não deve gerar evento, como em um pedido
// de ajuda, no --version, no root ou em um grupo chamado sem subcomando
// (que só mostram a ajuda) e nos comandos de untrackedCommands.
// Um comando que não existe continua sendo coletado, com UnknownCommand = true.
func BuildCommandInfo(root *cobra.Command, args []string, cmdErr error, cfg *mgcSdk.Config, nextUnknown NextUnknownFunc) (telemetry.CommandInfo, bool) {
	if errors.Is(cmdErr, schema_flags.ErrWantHelp) {
		return untracked()
	}

	cmd, rest, _ := root.Find(args)
	if cmd == nil {
		return telemetry.CommandInfo{UnknownCommand: true}, true
	}

	if flagChanged(cmd, "help") || (cmd == root && flagChanged(cmd, "version")) {
		return untracked()
	}

	if cmd == root || cmd.HasSubCommands() {
		if next, _ := nextUnknown(cmd, rest); next != nil {
			path := append(commandPathSegments(cmd), *next)
			return telemetry.CommandInfo{UnknownCommand: true, Path: path, OptionsSet: changedFlagNames(cmd)}, true
		}
		return untracked()
	}

	path := commandPathSegments(cmd)
	if _, skip := untrackedCommands[path[0]]; skip {
		return untracked()
	}

	return telemetry.CommandInfo{
		Path:       path,
		OptionsSet: changedFlagNames(cmd),
		ResourceID: resourceID(cmd),
		Region:     configValue(cmd, "region", cfg),
		Env:        configValue(cmd, "env", cfg),
	}, true
}

func configValue(cmd *cobra.Command, name string, cfg *mgcSdk.Config) string {
	if cfg == nil {
		return ""
	}
	f := cmd.Flags().Lookup(name)
	if f == nil {
		return ""
	}
	if _, ok := f.Value.(schema_flags.SchemaFlagValue); !ok {
		return ""
	}
	value, err := schema_flags.GetFlagValue(f, cfg)
	if err != nil || value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// commandPathSegments devolve o caminho do comando resolvido, sem o "mgc" na frente.
func commandPathSegments(cmd *cobra.Command) []string {
	fullPath := cmd.CommandPath()
	words := strings.Fields(fullPath)
	return words[1:]
}

func untracked() (telemetry.CommandInfo, bool) {
	return telemetry.CommandInfo{}, false
}

func flagChanged(cmd *cobra.Command, name string) bool {
	f := cmd.Flags().Lookup(name)
	return f != nil && f.Changed
}

func changedFlagNames(cmd *cobra.Command) []string {
	var names []string
	cmd.Flags().Visit(func(f *flag.Flag) {
		if f.Name != "help" {
			names = append(names, f.Name)
		}
	})
	return names
}

func resourceID(cmd *cobra.Command) string {
	f := cmd.Flags().Lookup("id")
	if f == nil {
		return ""
	}
	if value := f.Value.String(); f.Changed || (value != "" && value != f.DefValue) {
		return value
	}
	return ""
}
