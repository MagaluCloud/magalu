package spec

import (
	"fmt"
	"os"

	mergers "github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/spec/merger"
	"github.com/spf13/cobra"
)

func mergeSpecs(options *MergeSpecs) {
	merger := mergers.NewSpecMerger()

	mOptions := &mergers.MergeOptions{
		DowngradeToVersion: "3.0.3",
	}

	err := merger.MergeSpecs(options.specA, options.specB, options.output, mOptions)
	if err != nil {
		fmt.Printf("Erro ao mesclar especificações: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Especificações mescladas com sucesso. Resultado salvo em '" + options.output + "'")
}

type MergeSpecs struct {
	specA  string
	specB  string
	output string
}

func MergeSpecsCmd() *cobra.Command {
	options := &MergeSpecs{}

	cmd := &cobra.Command{
		Use:     "merge",
		Short:   "Mescla duas especificações OpenAPI",
		Example: "merge -a speca.yaml -b specb.yaml -o output.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			mergeSpecs(options)
		},
	}

	cmd.Flags().StringVarP(&options.specA, "speca", "a", "", "Caminho para o primeiro arquivo de especificação OpenAPI")
	cmd.Flags().StringVarP(&options.specB, "specb", "b", "", "Caminho para o segundo arquivo de especificação OpenAPI")
	cmd.Flags().StringVarP(&options.output, "output", "o", "", "Nome do arquivo de saída para a especificação mesclada")

	_ = cmd.MarkFlagRequired("speca")
	_ = cmd.MarkFlagRequired("output")
	return cmd
}
