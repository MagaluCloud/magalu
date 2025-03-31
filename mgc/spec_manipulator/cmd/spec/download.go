package spec

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var DownloadSpecsCmd = &cobra.Command{
	Use:   "download",
	Short: "Download all available specs",
	Run: func(cmd *cobra.Command, args []string) {
		_ = verificarEAtualizarDiretorio(CurrentDir())

		currentConfig, err := loadList()

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, v := range currentConfig {
			_ = getAndSaveFile(v.Url, filepath.Join(CurrentDir(), v.File))
		}
		fmt.Println("Now, run '" + cmd.Root().Name() + " prepare'")

	},
}
