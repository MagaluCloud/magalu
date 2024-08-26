package spec

import (
	"fmt"
	"os"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func main() {
	// Carregue as duas especificações
	spec1, err := loadSpec("spec1.yaml")
	if err != nil {
		fmt.Printf("Erro ao carregar spec1: %v\n", err)
		return
	}

	spec2, err := loadSpec("spec2.yaml")
	if err != nil {
		fmt.Printf("Erro ao carregar spec2: %v\n", err)
		return
	}

	// Mescle as especificações
	mergedSpec, err := mergeSpecs(spec1, spec2)
	if err != nil {
		fmt.Printf("Erro ao mesclar specs: %v\n", err)
		return
	}

	// Valide a especificação mesclada
	if err := validateSpec(mergedSpec); err != nil {
		fmt.Printf("Erro na validação da spec mesclada: %v\n", err)
		return
	}

	// Salve a especificação mesclada
	if err := saveSpec(mergedSpec, "merged_spec.yaml"); err != nil {
		fmt.Printf("Erro ao salvar a spec mesclada: %v\n", err)
		return
	}

	fmt.Println("Especificações mescladas com sucesso!")
}

func loadSpec(filename string) (*v3.Document, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	doc, err := libopenapi.NewDocument(data)
	if err != nil {
		return nil, err
	}

	return doc.V3(), nil
}

func mergeSpecs(spec1, spec2 *v3.Document) (*v3.Document, error) {
	// mergedSpec := &v3.Document{
	// 	OpenAPI: spec1.OpenAPI,
	// 	Info:    spec1.Info, // Você pode escolher qual info manter ou combinar
	// 	Paths:   make(v3.Paths),
	// 	Components: &v3.Components{
	// 		Schemas: make(v3.SchemaMap),
	// 		// Adicione outros componentes conforme necessário
	// 	},
	// }

	// // Mesclar caminhos
	// for path, item := range spec1.Paths {
	// 	mergedSpec.Paths[path] = item
	// }
	// for path, item := range spec2.Paths {
	// 	if _, exists := mergedSpec.Paths[path]; exists {
	// 		// Lidar com conflitos de caminho, se necessário
	// 		fmt.Printf("Aviso: Caminho duplicado encontrado: %s\n", path)
	// 	}
	// 	mergedSpec.Paths[path] = item
	// }

	// // Mesclar schemas
	// for name, schema := range spec1.Components.Schemas {
	// 	mergedSpec.Components.Schemas[name] = schema
	// }
	// for name, schema := range spec2.Components.Schemas {
	// 	if _, exists := mergedSpec.Components.Schemas[name]; exists {
	// 		// Lidar com conflitos de schema, se necessário
	// 		fmt.Printf("Aviso: Schema duplicado encontrado: %s\n", name)
	// 	}
	// 	mergedSpec.Components.Schemas[name] = schema
	// }

	// Mesclar outros componentes conforme necessário

	return mergedSpec, nil
}

func validateSpec(spec *v3.Document) error {
	// Implemente a validação da especificação aqui
	// A libopenapi pode não fornecer validação direta, então você pode precisar
	// implementar sua própria lógica de validação ou usar outra biblioteca
	return nil
}

func saveSpec(spec *v3.Document, filename string) error {
	// Implemente a lógica para salvar a especificação em um arquivo
	// Você pode precisar converter o Document de volta para YAML ou JSON
	return nil
}

var MergeSpecsCmd = &cobra.Command{
	Use:   "merge",
	Short: "List all available specs",
	Run: func(cmd *cobra.Command, args []string) {
		_ = verificarEAtualizarDiretorio(CurrentDir())

		currentConfig, err := loadList()

		if err != nil {
			fmt.Println(err)
			return
		}

		out, err := yaml.Marshal(currentConfig)
		if err == nil {
			fmt.Println(string(out))
		}

	},
}
