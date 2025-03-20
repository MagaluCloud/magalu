package cmd

import (
	"github.com/spf13/cobra"
)

func addBaseURLFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		"base-url",
		"",
		"URL to override the default host. Ex. https://api.magalu.com.br or http://localhost/v1/route",
	)
	_ = cmd.PersistentFlags().MarkHidden("base-url")
}

func getBaseURLFlag(cmd *cobra.Command) string {
	host, err := cmd.Root().PersistentFlags().GetString("base-url")
	if err != nil {
		return ""
	}
	return host
}
