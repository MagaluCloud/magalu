package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var lisCommandsCmd = &cobra.Command{
	Use:    "list",
	Short:  "List all available commands",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		currentCommands, err := loadList()

		if err != nil {
			fmt.Println(err)
			return
		}

		out, err := yaml.Marshal(currentCommands)
		if err == nil {
			fmt.Println(string(out))
		}

	},
}
