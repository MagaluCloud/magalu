package spec

import (
	"fmt"
	"os"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/spec/customizer"
	"github.com/spf13/cobra"
)

func customizeSpec(options *CustomizeSpec) {
	magalizer := customizer.NewMagaluCustomizer()

	var paramsToRemove []string
	if options.removeParams != "" {
		paramsToRemove = strings.Split(options.removeParams, ",")
		for i, param := range paramsToRemove {
			paramsToRemove[i] = strings.TrimSpace(param)
		}
	}

	customOptions := &customizer.CustomizeOptions{
		IncludeRegion:       options.includeRegion,
		IncludeGlobalRegion: options.includeGlobalRegion,
		ProductPathURL:      options.productPathURL,
		DowngradeToVersion:  options.downgradeToVersion,
		ParamsToRemove:      paramsToRemove,
		ConfigureSecurity:   options.configureSecurity,
	}

	err := magalizer.CustomizeSpec(options.specIn, options.output, customOptions)
	if err != nil {
		fmt.Printf("Erro ao personalizar especificação: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Especificação personalizada com sucesso. Resultado salvo em '" + options.output + "'")

	if len(paramsToRemove) > 0 {
		fmt.Println("Parâmetros removidos: " + strings.Join(paramsToRemove, ", "))
	}
}

type CustomizeSpec struct {
	specIn              string
	output              string
	includeRegion       bool
	includeGlobalRegion bool
	productPathURL      string
	downgradeToVersion  string
	removeParams        string
	configureSecurity   bool
}

func CustomizeSpecCmd() *cobra.Command {
	options := &CustomizeSpec{}

	cmd := &cobra.Command{
		Use:     "customize",
		Short:   "Personaliza uma especificação OpenAPI para o padrão Magalu",
		Example: "customize -i spec.yaml -o output.yaml -p product-path -r -g -v 3.0.3 --remove-params 'param1,param2' --configure-security",
		Run: func(cmd *cobra.Command, args []string) {
			customizeSpec(options)
		},
	}

	cmd.Flags().StringVarP(&options.specIn, "spec-in", "i", "", "Arquivo de especificação OpenAPI de entrada")
	cmd.Flags().StringVarP(&options.output, "output", "o", "", "Arquivo de saída para a especificação personalizada")
	cmd.Flags().BoolVarP(&options.includeRegion, "region", "r", false, "Incluir variável de região na URL")
	cmd.Flags().BoolVarP(&options.includeGlobalRegion, "global-region", "g", false, "Incluir região global na lista de regiões")
	cmd.Flags().StringVarP(&options.productPathURL, "product", "p", "", "Caminho do produto na URL")
	cmd.Flags().StringVarP(&options.downgradeToVersion, "version", "v", "3.0.3", "Versão OpenAPI para downgrade (padrão: 3.0.3)")
	cmd.Flags().StringVar(&options.removeParams, "remove-params", "", "Lista de parâmetros a serem removidos, separados por vírgula (ex: 'param1,param2')")
	cmd.Flags().BoolVar(&options.configureSecurity, "configure-security", false, "Configurar informações de segurança OAuth2 nas rotas que não possuem")

	_ = cmd.MarkFlagRequired("spec-in")
	_ = cmd.MarkFlagRequired("output")
	_ = cmd.MarkFlagRequired("product")

	return cmd
}
