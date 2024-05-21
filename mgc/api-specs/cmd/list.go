package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var listSpecsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available specs",
	Run: func(cmd *cobra.Command, args []string) {
		var fromViveiro bool
		cmd.Flags().BoolVarP(&fromViveiro, "viveiro", "v", false, "Função utilizando viveiro")
		_ = verificarEAtualizarDiretorio(SPEC_DIR)

		currentConfig, err := loadListFromViper(toWriteViveiro(fromViveiro))

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
