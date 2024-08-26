package spec

import (
	"github.com/spf13/cobra"
)

func SpecCmd() *cobra.Command {
	specMenu := &cobra.Command{
		Use:   "spec",
		Short: "Exibe o submenu",
	}

	specMenu.AddCommand(DownloadSpecsCmd) // download all
	specMenu.AddCommand(SpecAddNewCmd())  // add spec
	specMenu.AddCommand(DeleteSpecsCmd)   // delete spec
	specMenu.AddCommand(ListSpecsCmd)     // list specs
	specMenu.AddCommand(PrepareToGoCmd)   // convert spec to golang
	specMenu.AddCommand(DowngradeSpecCmd) // downgrade spec

	return specMenu
}
