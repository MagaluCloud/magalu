package cmd

import "github.com/spf13/cobra"

func addRawOutputFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(
		"raw",
		false,
		"Output raw data, without any formatting or coloring",
	)
}

func getRawOutputFlag(cmd *cobra.Command) bool {
	raw, err := cmd.Root().PersistentFlags().GetBool("raw")
	if err != nil {
		return false
	}
	return raw
}
