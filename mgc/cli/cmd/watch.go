package cmd

import (
	"github.com/spf13/cobra"
)

const watchFlag = "cli.watch"

func addWatchFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(
		watchFlag,
		false,
		`Wait until the operation is completed by calling the 'get' link and waiting until termination. Akin to '! get -w'`,
	)
}

func getWatchFlag(cmd *cobra.Command) bool {
	w, err := cmd.PersistentFlags().GetBool(watchFlag)
	if err != nil {
		return false
	}
	return w
}
