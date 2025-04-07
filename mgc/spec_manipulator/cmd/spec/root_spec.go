package spec

import (
	"github.com/spf13/cobra"
)

func SpecCmd() *cobra.Command {
	specMenu := &cobra.Command{
		Use:   "spec",
		Short: "Menu com opções para manipulação de specs",
	}

	specMenu.AddCommand(downloadSpecsCmd()) // download all
	specMenu.AddCommand(specAddNewCmd())    // add spec
	specMenu.AddCommand(deleteSpecsCmd)     // delete spec
	specMenu.AddCommand(listSpecsCmd())     // list specs
	specMenu.AddCommand(prepareToGoCmd())   // convert spec to golang
	specMenu.AddCommand(downgradeSpec())    // downgrade spec
	specMenu.AddCommand(mergeSpecsCmd())    // spc merge

	return specMenu
}
