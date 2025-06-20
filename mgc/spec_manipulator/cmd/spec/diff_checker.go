package spec

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/tui"
	// oasChanges "github.com/pb33f/openapi-changes"
	"github.com/spf13/cobra"
)

func diffCheckerCmd() *cobra.Command {
	var dir string
	var menu string

	cmd := &cobra.Command{
		Use:   "diff [dir] [menu]",
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
				dir = filepath.Join(dir, "tmp")
				os.MkdirAll(dir, 0755)

				tmpFile := filepath.Join(dir, v.File)

				if !strings.Contains(v.Url, "gitlab.luizalabs.com") {
					err = getAndSaveFile(v.Url, tmpFile, v.Menu)
					if err != nil {
						return
					}
				}

				if strings.Contains(v.Url, "gitlab.luizalabs.com") {
					err = downloadGitlab(v.Url, tmpFile)
					if err != nil {
						return
					}
				}

				justRunValidate(dir, v)

				//

			}
			spinner.Success("Specs downloaded successfully")
		},
	}
	cmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to save the converted specs")
	cmd.Flags().StringVarP(&menu, "menu", "m", "", "Menu to download the specs")
	return cmd
}
