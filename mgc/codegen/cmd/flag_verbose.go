package cmd

import (
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

const (
	flagVerbose = "verbose"
)

func addVerboseFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&verbose, flagVerbose, "v", false, "be verbose while generating")
}

func getVerboseFlag(cmd *cobra.Command) bool {
	return verbose
}
