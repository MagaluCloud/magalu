package spec

import (
	"fmt"
	"os"

	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/spec/downgrader"
	"github.com/spf13/cobra"
)

func downgradeFile(options *DowngradeUniqueSpec) {
	downgrader := downgrader.NewOpenAPIDowngrader()
	err := downgrader.DowngradeFile(options.specIN, options.specOUT)
	if err != nil {
		fmt.Printf("Erro ao converter: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Arquivo convertido com sucesso: %s\n", options.specOUT)
}

type DowngradeUniqueSpec struct {
	specIN  string
	specOUT string
}

func DowngradeUniqueSpecCmd() *cobra.Command {
	options := &DowngradeUniqueSpec{}

	cmd := &cobra.Command{
		Use:    "downgrade",
		Short:  "Downgrade specs from 3.1.x to 3.0.x",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			downgradeFile(options)
		},
	}

	cmd.Flags().StringVarP(&options.specIN, "spec-in", "i", "", "Spec to downgrade")
	cmd.Flags().StringVarP(&options.specOUT, "spec-out", "o", "", "Spec to downgrade")

	return cmd
}
