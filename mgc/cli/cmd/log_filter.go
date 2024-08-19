package cmd

import (
	"github.com/spf13/cobra"
	mgcSdk "magalu.cloud/sdk"
)

const logFilterFlag = "cli.log"

func addLogFilterFlag(cmd *cobra.Command, def string) {
	if def == "" {
		def = "warn+:*"
	}
	cmd.Root().PersistentFlags().StringP(
		logFilterFlag,
		"l",
		def,
		`Format is 'levels:namespaces'. Use 'info+:*' to show info for all levels, use '*:*' to show all logs.
See more details about the filter syntax at https://github.com/moul/zapfilter`,
	)
	cmd.Root().PersistentFlags().Lookup(logFilterFlag).Hidden = true
}

func getLogFilterFlag(cmd *cobra.Command) string {
	// debug take precedence over cli.logs
	dbgFlagResult := getLogDebugFlag(cmd)
	if dbgFlagResult != "" {
		return dbgFlagResult
	}
	return cmd.Root().PersistentFlags().Lookup(logFilterFlag).Value.String()
}

// TODO: Bind config to PFlag. Investigate how to make it work correctly
func getLogFilterConfig(sdk *mgcSdk.Sdk) string {
	var logfilter string
	err := sdk.Config().Get("logfilter", &logfilter)
	if err != nil {
		return ""
	}
	return logfilter
}
