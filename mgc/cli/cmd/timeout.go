package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

const timeoutFlag = "cli.timeout"

func addTimeoutFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().DurationP(
		timeoutFlag,
		"t",
		0.0,
		`If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s`,
	)
	cmd.Root().PersistentFlags().Lookup(timeoutFlag).Hidden = false
}

func getTimeoutFlag(cmd *cobra.Command) time.Duration {
	t, err := cmd.Root().PersistentFlags().GetDuration(timeoutFlag)
	if err != nil {
		return 0.0
	}
	return t
}
