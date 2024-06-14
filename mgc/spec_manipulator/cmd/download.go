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
		_ = verificarEAtualizarDiretorio(SPEC_DIR)

		currentConfig, err := loadList()

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, v := range currentConfig {
			// _ = removerArquivosOld(filepath.Join(SPEC_DIR))
			// _ = verificarERenomearArquivo(filepath.Join(SPEC_DIR, v.File))
			_ = getAndSaveFile(v.Url, filepath.Join(SPEC_DIR, v.File))
		}
		runPrepare(cmd, args)

	},
}
