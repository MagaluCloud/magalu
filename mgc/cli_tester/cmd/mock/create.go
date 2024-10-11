package mock

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mockCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("xpto")

	},
}
