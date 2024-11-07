package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var downloadSpecsCmd = &cobra.Command{
	Use:    "download",
	Short:  "Download all available specs",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = verificarEAtualizarDiretorio(currentDir())

		currentConfig, err := loadList()

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, v := range currentConfig {
			_ = getAndSaveFile(v.Url, filepath.Join(currentDir(), v.File))
		}
		fmt.Println("Now, run '" + cmd.Root().Name() + " prepare'")

	},
}
