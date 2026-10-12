package cmd

import (
	"errors"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var untrackedCommands = map[string]struct{}{
	"help":                          {},
	"completion":                    {},
	cobra.ShellCompRequestCmd:       {},
	cobra.ShellCompNoDescRequestCmd: {},
	"dump-tree":                     {},
	telemetryCmdName:                {},
}

// telemetryCommandInfo descobre qual comando foi executado e monta o CommandInfo.
// O segundo retorno é false quando a execução não deve gerar evento, como em um pedido
// de ajuda, no --version, no root ou em um grupo chamado sem subcomando
// (que só mostram a ajuda) e nos comandos de untrackedCommands.
// Um comando que não existe continua sendo coletado, com UnknownCommand = true.
func telemetryCommandInfo(root *cobra.Command, args []string, cmdErr error) (telemetry.CommandInfo, bool) {
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
		if next, _ := getNextUnknownCommand(cmd, rest); next != nil {
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
	}, true
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
