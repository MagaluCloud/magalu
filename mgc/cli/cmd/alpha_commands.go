package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// isAlphaCommand tells whether cmd comes from an installed alpha module.
func isAlphaCommand(cmd *cobra.Command) bool {
	return strings.HasPrefix(cmd.CommandPath(), cmd.Root().Name()+" alpha ")
}

// isAlphaModule tells whether args name a module right under 'mgc alpha'.
func isAlphaModule(root *cobra.Command, args []string) bool {
	found, _, err := root.Find(args)
	return err == nil && found.CommandPath() == root.Name()+" alpha"
}
