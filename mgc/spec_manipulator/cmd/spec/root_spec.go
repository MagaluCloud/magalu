package spec

import (
	"github.com/spf13/cobra"
)

func SpecCmd() *cobra.Command {
	specMenu := &cobra.Command{
		Use:   "spec",
		Short: "Menu com opções para manipulação de specs",
	}

	specMenu.AddCommand(MergeSpecsCmd()) // spc merge
	specMenu.AddCommand(DowngradeUniqueSpecCmd())
	specMenu.AddCommand(CustomizeSpecCmd()) // personalizar para padrão Magalu
	specMenu.AddCommand(DownloadSpecCmd())  // downgrade spec
	return specMenu
}
