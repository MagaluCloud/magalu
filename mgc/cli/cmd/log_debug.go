package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const logDebugFlag = "debug"
const logDebugDef = "info+:*"

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
			return getDebugLevelFromOS()
		}
	}

	return ""
}

func getDebugLevelFromOS() string {
	if result := os.Getenv("MGC_SDK_LOG_DEBUG"); result != "" {
		return strings.ToLower(result) + "+:*"
	}

	return logDebugDef
}
