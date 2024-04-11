package cmd

import "github.com/spf13/cobra"

const bypassConfirmationFlag = "no-confirm"

func addBypassConfirmationFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().Bool(
		bypassConfirmationFlag,
		false,
		"Bypasses confirmation step for commands that ask a confirmation from the user",
	)
}

func getBypassConfirmationFlag(cmd *cobra.Command) bool {
	allow, err := cmd.Root().PersistentFlags().GetBool(bypassConfirmationFlag)
	if err != nil {
		return false
	}
	return allow
}
