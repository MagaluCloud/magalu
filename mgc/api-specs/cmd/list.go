package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var listSpecsCmd = &cobra.Command{
	Use:    "list",
	Short:  "List all available specs",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = verificarEAtualizarDiretorio(SPEC_DIR)

		currentConfig, err := loadList()

		if err != nil {
			fmt.Println(err)
			return
		}

		out, err := yaml.Marshal(currentConfig)
		if err == nil {
			fmt.Println(string(out))
		}

	},
}
