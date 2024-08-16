package cmd

import "github.com/spf13/cobra"

const showHiddenFlag = "cli.show-hidden-flag"

func addShowHiddenFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().Bool(
		showHiddenFlag,
		false,
		"Show hidden flags",
	)
	cmd.Root().PersistentFlags().Lookup(showHiddenFlag).Hidden = true
}

func getShowHiddenFlag(cmd *cobra.Command) bool {
	show, err := cmd.Root().PersistentFlags().GetBool(showHiddenFlag)
	if err != nil {
		return false
	}
	return show
}
