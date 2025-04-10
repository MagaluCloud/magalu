package spec

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func listSpecsCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available specs",
		Run: func(cmd *cobra.Command, args []string) {
			_ = verificarEAtualizarDiretorio(dir)

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
	cmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to save the converted specs")
	return cmd
}
