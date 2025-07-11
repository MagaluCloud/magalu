package spec

import (
	"github.com/spf13/cobra"
)

func SpecCmd() *cobra.Command {
	specMenu := &cobra.Command{
		Use:   "specs",
		Short: "Menu com opções para manipulação de specs",
	}

	specMenu.AddCommand(downloadSpecsCmd()) // download all
	specMenu.AddCommand(specAddNewCmd())    // add spec
	specMenu.AddCommand(deleteSpecCmd())    // delete spec
	specMenu.AddCommand(listSpecsCmd())     // list specs
	specMenu.AddCommand(prepareToGoCmd())   // convert spec to golang
	specMenu.AddCommand(downgradeSpec())    // downgrade spec
	specMenu.AddCommand(mergeSpecsCmd())    // spc merge
	specMenu.AddCommand(validateSpec())     // validate spec
	specMenu.AddCommand(diffCheckerCmd())   // diff checker

	return specMenu
}
