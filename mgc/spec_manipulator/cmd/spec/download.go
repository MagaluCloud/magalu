package spec

import (
	"path/filepath"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/tui"
	"github.com/spf13/cobra"
)

func downloadSpecsCmd() *cobra.Command {
	var dir string
	var menu string

	cmd := &cobra.Command{
		Use:   "download [dir] [menu]",
		Short: "Download available spec",
		Run: func(cmd *cobra.Command, args []string) {

			_ = verificarEAtualizarDiretorio(dir)

			var currentConfig []specList
			var err error

			if menu != "" {
				currentConfig, err = loadList(menu)
			} else {
				currentConfig, err = getConfigToRun()
			}
			if err != nil {
				return
			}
			spinner := tui.NewSpinner()
			spinner.Start("Downloading ...")
			for _, v := range currentConfig {
				spinner.UpdateText("Downloading " + v.File)
				if !strings.Contains(v.Url, "gitlab.luizalabs.com") {
					err = getAndSaveFile(v.Url, filepath.Join(dir, v.File), v.Menu)
					if err != nil {
						spinner.Fail(err)
						return
					}
				}

				if strings.Contains(v.Url, "gitlab.luizalabs.com") {
					err = downloadGitlab(v.Url, filepath.Join(dir, v.File))
					if err != nil {
						spinner.Fail(err)
						return
					}
				}

				justRunValidate(dir, v)
			}
			spinner.Success("Specs downloaded successfully")
		},
	}
	cmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to save the converted specs")
	cmd.Flags().StringVarP(&menu, "menu", "m", "", "Menu to download the specs")
	return cmd
}
