package spec

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func downloadSpecsCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download all available specs",
		Run: func(cmd *cobra.Command, args []string) {
			_ = verificarEAtualizarDiretorio(dir)

			currentConfig, err := loadList()

			if err != nil {
				fmt.Println(err)
				return
			}

			for _, v := range currentConfig {
				_ = getAndSaveFile(v.Url, filepath.Join(dir, v.File))
			}
			fmt.Println("Now, run '" + cmd.Root().Name() + " prepare'")
		},
	}
	cmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to save the converted specs")
	return cmd
}
