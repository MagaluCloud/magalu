package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.0.0")
	},
}
