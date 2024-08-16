package cmd

import (
	"github.com/spf13/cobra"
)

const showCliGlobalFlags = "cli.show-cli-globals"

func addShowCliGlobalFlags(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().Bool(
		showCliGlobalFlags,
		false,
		"Show all CLI global flags on usage text",
	)
	f := cmd.Root().PersistentFlags().Lookup(showCliGlobalFlags)
	f.NoOptDefVal = "true"
}

func getShowCliGlobalFlags(cmd *cobra.Command) bool {
	value, err := cmd.Root().PersistentFlags().GetBool(showCliGlobalFlags)
	if err != nil {
		return false
	}
	return value
}
