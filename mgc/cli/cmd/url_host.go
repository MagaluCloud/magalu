package cmd

import (
	"github.com/spf13/cobra"
)

func addHostFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		"host",
		"",
		"URL to override the default host. Ex. https://api.magalu.com.br or http://localhost/v1/route",
	)
	// _ = cmd.PersistentFlags().MarkHidden("host")
}

func getHostFlag(cmd *cobra.Command) string {
	host, err := cmd.Root().PersistentFlags().GetString("host")
	if err != nil {
		return ""
	}
	return host
}
