package cmd

import "github.com/spf13/cobra"

const showInternalFlag = "cli.show-internal"

func addShowInternalFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().Bool(
		showInternalFlag,
		false,
		"Show internal groups and executors",
	)
}

func getShowInternalFlag(cmd *cobra.Command) bool {
	show, err := cmd.Root().PersistentFlags().GetBool(showInternalFlag)
	if err != nil {
		return false
	}
	return show
}
