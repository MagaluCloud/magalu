package cmd

import (
	"github.com/spf13/cobra"
)

const logDebugFlag = "debug"
const logDebugDef = "debug+:*"

func addLogDebugFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().Bool(
		logDebugFlag,
		false,
		`Display detailed log information at the debug level`,
	)
}

func getLogDebugFlag(cmd *cobra.Command) string {
	if result, ok := cmd.Root().PersistentFlags().GetBool(logDebugFlag); ok == nil {
		if result {
			return logDebugDef
		}
	}

	return ""
}
