package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func addRawOutputFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(
		"raw",
		"r",
		false,
		"Output raw data, without any formatting or coloring",
	)
}

// getRawOutputFlag reports whether the CLI prints raw output. That is the case with --raw, and
// also when color.NoColor is set: stdout is not a terminal (a pipe or a file), NO_COLOR is set,
// or TERM is "dumb".
func getRawOutputFlag(cmd *cobra.Command) bool {
	raw, err := cmd.Root().PersistentFlags().GetBool("raw")
	if err != nil {
		return color.NoColor
	}
	return raw || color.NoColor
}
