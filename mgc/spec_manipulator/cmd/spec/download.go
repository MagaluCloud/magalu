package spec

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

type DownloadOptions struct {
	All      bool
	Menu     string
	SpecsDir string
}

func DownloadSpecsCmd() *cobra.Command {
	var opts DownloadOptions
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download all available specs",
		Run: func(cmd *cobra.Command, args []string) {
			_ = verificarEAtualizarDiretorio(opts.SpecsDir)

			currentConfig, err := loadList()

			if err != nil {
				fmt.Println(err)
				return
			}

			for _, v := range currentConfig {
				if v.Enabled && ((opts.Menu == "" && opts.All) || v.Menu == opts.Menu) {
					_ = getAndSaveFile(v.Url, filepath.Join(opts.SpecsDir, v.File))
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Download all specs")
	cmd.Flags().StringVarP(&opts.Menu, "menu", "m", "", "Menu to download")
	cmd.Flags().StringVarP(&opts.SpecsDir, "specs-dir", "s", "", "Directory containing the specs")

	return cmd
}
