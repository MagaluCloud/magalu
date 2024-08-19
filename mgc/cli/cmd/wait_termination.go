package cmd

import (
	"github.com/spf13/cobra"
)

const (
	waitTerminationFlag string = "cli.wait-termination"
)

func addWaitTerminationFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().BoolP(
		waitTerminationFlag,
		"w",
		false,
		`Wait for any asynchronous actions to transition to their final state. Note that not all actions
implement ExecuteUntilTermination(), if that is not available, then regular Execute()
is used and no wait will be done`,
	)

	f := cmd.Root().PersistentFlags().Lookup(waitTerminationFlag)
	f.NoOptDefVal = "true"
	f.Hidden = true
}

func getWaitTerminationFlag(cmd *cobra.Command) bool {
	v, err := cmd.Root().PersistentFlags().GetBool(waitTerminationFlag)
	if err != nil {
		return false
	}

	return v
}
