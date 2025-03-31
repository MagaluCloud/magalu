package spec

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ListSpecsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available specs",
	Run: func(cmd *cobra.Command, args []string) {
		_ = verificarEAtualizarDiretorio(CurrentDir())

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
