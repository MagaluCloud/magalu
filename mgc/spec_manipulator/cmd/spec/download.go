package spec

import (
	"fmt"
	"os"

	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/spec/downloader"
	"github.com/spf13/cobra"
)

func downloadSpecFile(options *DownloadSpec) {
	down := downloader.NewSpecDownloader()

	down.ValidateSpec = !options.skipValidation

	err := down.DownloadSpec(options.source, options.destination)
	if err != nil {
		fmt.Printf("Erro ao baixar/carregar especificação: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Especificação salva com sucesso em: %s\n", options.destination)
}

type DownloadSpec struct {
	source         string
	destination    string
	skipValidation bool
}

func DownloadSpecCmd() *cobra.Command {
	options := &DownloadSpec{}

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Baixa ou carrega uma especificação OpenAPI de uma URL ou arquivo local",
		Example: `
  # Baixar de uma URL:
  spec download -s https://example.com/api/openapi.json -d ./specs/

  # Carregar de um arquivo local (com prefixo @):
  spec download -s @./caminho/para/spec.yaml -d ./specs/destino.yaml
  
  # Baixar sem validar a especificação:
  spec download -s https://example.com/api/openapi.json -d ./specs/ --skip-validation`,
		Run: func(cmd *cobra.Command, args []string) {
			downloadSpecFile(options)
		},
	}

	cmd.Flags().StringVarP(&options.source, "source", "s", "", "URL ou caminho do arquivo de origem (com prefixo @ para arquivos locais)")
	cmd.Flags().StringVarP(&options.destination, "destination", "d", "", "Caminho do arquivo ou diretório de destino")
	cmd.Flags().BoolVar(&options.skipValidation, "skip-validation", false, "Pular a validação da especificação OpenAPI")

	_ = cmd.MarkFlagRequired("source")
	_ = cmd.MarkFlagRequired("destination")

	return cmd
}
