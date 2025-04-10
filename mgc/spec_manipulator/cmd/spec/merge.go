package spec

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func mergeSpecsMain(options MergeSpecs) {
	// Carregar a primeira especificação
	file := filepath.Join(options.specA)
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	docA, err := libopenapi.NewDocument(fileBytes)
	if err != nil {
		fmt.Printf("Erro ao carregar a especificação A: %v\n", err)
		os.Exit(1)
	}

	// Carregar a segunda especificação
	file = filepath.Join(options.specB)
	fileBytes, err = os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	docB, err := libopenapi.NewDocument(fileBytes)
	if err != nil {
		fmt.Printf("Erro ao carregar a especificação B: %v\n", err)
		os.Exit(1)
	}

	// Realizar o merge
	mergedSpec := mergeSpecs(docA, docB, options)

	// Salvar a especificação mesclada
	err = saveSpec(mergedSpec, options.output)
	if err != nil {
		fmt.Printf("Erro ao salvar a especificação mesclada: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Especificações mescladas com sucesso. Resultado salvo em '" + options.output + "'")
}

func mergeSpecs(specA, specB libopenapi.Document, options MergeSpecs) libopenapi.Document {
	mergedSpec := specA

	//Add Servers
	createExtension := func(value interface{}) *yaml.Node {
		node := &yaml.Node{}
		err := node.Encode(value)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return node
	}

	mergedSpecA, _ := mergedSpec.BuildV3Model()
	specModelB, _ := specB.BuildV3Model()

	url := "https://{env}/{region}/" + options.productPathURL
	if options.globalDb {
		url = "https://{env}/" + options.productPathURL
	}

	server := &v3.Server{
		URL:       url,
		Variables: orderedmap.New[string, *v3.ServerVariable](),
	}

	regionVar := &v3.ServerVariable{
		Enum:       []string{"br-ne-1", "br-se1", "br-mgl1"},
		Extensions: orderedmap.New[string, *yaml.Node](),
	}

	regionVar.Description = "Region to reach the service"
	regionVar.Default = "br-se1"

	regionTransforms := []map[string]interface{}{
		{
			"type":         "translate",
			"allowMissing": true,
			"translations": []map[string]string{
				{"from": "br-ne1", "to": "br-ne-1"},
				{"from": "br-mgl1", "to": "br-se-1"},
			},
		},
	}
	regionVar.Extensions.Set("x-mgc-transforms", createExtension(regionTransforms))
	server.Variables.Set("region", regionVar)

	envVar := &v3.ServerVariable{
		Enum:       []string{"api.magalu.cloud", "api.pre-prod.jaxyendy.com"},
		Extensions: orderedmap.New[string, *yaml.Node](),
	}
	envVar.Description = "Environment to use"
	envVar.Default = "api.magalu.cloud"
	envTransforms := []map[string]interface{}{
		{
			"type": "translate",
			"translations": []map[string]string{
				{"from": "prod", "to": "api.magalu.cloud"},
				{"from": "pre-prod", "to": "api.pre-prod.jaxyendy.com"},
			},
		},
	}
	envVar.Extensions.Set("x-mgc-transforms", createExtension(envTransforms))
	server.Variables.Set("env", envVar)

	mergedSpecA.Model.Servers = append(mergedSpecA.Model.Servers, server)

	// Merge paths
	for path := specModelB.Model.Paths.PathItems.Oldest(); path != nil; path = path.Next() {
		if pathItemA, isPathPresent := mergedSpecA.Model.Paths.PathItems.Get(path.Key); isPathPresent {
			// Path existe no B e também no A
			operationsA := pathItemA.GetOperations()
			operationsB := path.Value.GetOperations()

			for opB := operationsB.Oldest(); opB != nil; opB = opB.Next() {
				if operationA, isOpPresent := operationsA.Get(opB.Key); isOpPresent {

					if opB.Value.Tags != nil {
						operationA.Tags = opB.Value.Tags
					}

					if opB.Value.Summary != "" {
						operationA.Summary = opB.Value.Summary
					}

					if opB.Value.Description != "" {
						operationA.Description = opB.Value.Description
					}

					if opB.Value.ExternalDocs != nil {
						operationA.ExternalDocs = opB.Value.ExternalDocs
					}

					if opB.Value.OperationId != "" {
						operationA.OperationId = opB.Value.OperationId
					}

					if opB.Value.RequestBody != nil {
						operationA.RequestBody = opB.Value.RequestBody
					}

					if opB.Value.Responses != nil {
						// operationA.Responses = opB.Value.Responses

						for cB := orderedmap.First(opB.Value.Responses.Codes); cB != nil; cB = cB.Next() {

							if cA, isCAPresent := operationA.Responses.Codes.Get(cB.Key()); isCAPresent {
								ccB := cB.Value()

								if ccB.Content != nil {
									cA.Content = ccB.Content
								}

								if ccB.Description != "" {
									cA.Description = ccB.Description
								}

								if ccB.Extensions != nil {
									cA.Extensions = ccB.Extensions
								}

								if ccB.Headers != nil {
									cA.Headers = ccB.Headers
								}

								if ccB.Links != nil {
									cA.Links = ccB.Links
								}

							} else {
								operationA.Responses.Codes.Set(cB.Key(), cB.Value())
							}
						}
					}

					// ARRAY ITEMS - ITERATE IT
					if opB.Value.Parameters != nil { //array
						// operationA.Parameters = opB.Value.Parameters

						for _, pB := range opB.Value.Parameters {
							for _, pA := range operationA.Parameters {
								if pA.Name == pB.Name {
									// SAME
									if pB.In != "" {
										pA.In = pB.In
									}
									if pB.Description != "" {
										pA.Description = pB.Description
									}
									if pB.Style != "" {
										pA.Style = pB.Style
									}

									if pB.Required != nil {
										pA.Required = pB.Required
									}

									if pB.Explode != nil {
										pA.Explode = pB.Explode
									}

									if pB.Schema != nil {
										pA.Schema = pB.Schema
									}
									if pB.Example != nil {
										pA.Example = pB.Example
									}
									if pB.Examples != nil {
										pA.Examples = pB.Examples
									}
									if pB.Content != nil {
										pA.Content = pB.Content
									}

									pA.Deprecated = pB.Deprecated
									pA.AllowEmptyValue = pB.AllowEmptyValue
									pA.AllowReserved = pB.AllowReserved

								}
							}
						}
					}

					if opB.Value.Deprecated != nil {
						operationA.Deprecated = opB.Value.Deprecated
					}
					if opB.Value.Callbacks != nil { //array
						operationA.Callbacks = opB.Value.Callbacks
					}
					if opB.Value.Security != nil { //array
						operationA.Security = opB.Value.Security
					}

					if opB.Value.Servers != nil { //array
						operationA.Servers = opB.Value.Servers
					}

					if opB.Value.Extensions != nil { //array
						operationA.Extensions = opB.Value.Extensions
					}

				}
			}
		} else {
			// Path existe no B e não no A
			mergedSpecA.Model.Paths.PathItems.Set(path.Key, path.Value)
		}

		if patchedItem, isOk := mergedSpecA.Model.Paths.PathItems.Get(path.Key); isOk {
			if patchedItem.Get != nil {
				for _, tag := range patchedItem.Get.Tags {
					alreadyExists := false
					for _, tt := range mergedSpecA.Model.Tags {
						if tt.Name == tag {
							alreadyExists = true
							break
						}

					}
					if !alreadyExists {
						mergedSpecA.Model.Tags = append(mergedSpecA.Model.Tags, &base.Tag{
							Name:        tag,
							Description: tag,
						})
					}
				}
			}

			if patchedItem.Put != nil {
				for _, tag := range patchedItem.Put.Tags {
					alreadyExists := false
					for _, tt := range mergedSpecA.Model.Tags {
						if tt.Name == tag {
							alreadyExists = true
							break
						}

					}
					if !alreadyExists {
						mergedSpecA.Model.Tags = append(mergedSpecA.Model.Tags, &base.Tag{
							Name:        tag,
							Description: tag,
						})
					}
				}
			}

			if patchedItem.Delete != nil {
				for _, tag := range patchedItem.Delete.Tags {
					alreadyExists := false
					for _, tt := range mergedSpecA.Model.Tags {
						if tt.Name == tag {
							alreadyExists = true
							break
						}

					}
					if !alreadyExists {
						mergedSpecA.Model.Tags = append(mergedSpecA.Model.Tags, &base.Tag{
							Name:        tag,
							Description: tag,
						})
					}
				}
			}

			if patchedItem.Post != nil {
				for _, tag := range patchedItem.Post.Tags {
					alreadyExists := false
					for _, tt := range mergedSpecA.Model.Tags {
						if tt.Name == tag {
							alreadyExists = true
							break
						}

					}
					if !alreadyExists {
						mergedSpecA.Model.Tags = append(mergedSpecA.Model.Tags, &base.Tag{
							Name:        tag,
							Description: tag,
						})
					}
				}
			}

			// add more http methods if necessary

		}
	}

	// Merge schemas
	for schema := specModelB.Model.Components.Schemas.Oldest(); schema != nil; schema = schema.Next() {
		if schemaItemA, isSchemaPresent := mergedSpecA.Model.Components.Schemas.Get(schema.Key); isSchemaPresent {
			// Schema existe no B e também no A
			*schemaItemA = *schema.Value
		} else {
			// Schema existe no B e não no A
			mergedSpecA.Model.Components.Schemas.Set(schema.Key, schema.Value)
		}
	}

	// Merge or create Tags
	for _, tags := range specModelB.Model.Tags {
		tagOnlyAtB := true
		for _, tagMerge := range mergedSpecA.Model.Tags {
			if tags.Name == tagMerge.Name {
				tagMerge = tags
				tagOnlyAtB = false
				continue
			}
		}
		if tagOnlyAtB {
			mergedSpecA.Model.Tags = append(mergedSpecA.Model.Tags, tags)
		}
	}

	return mergedSpec
}

func saveSpec(spec libopenapi.Document, filename string) error {
	model, errr := spec.BuildV3Model()
	if errr != nil {
		fmt.Printf("Erro ao gerar modelo do spec mesclado: %v\n", errr)
		os.Exit(1)
	}

	byteFile, err := model.Model.RenderJSON("  ")
	if err != nil {
		fmt.Printf("Erro ao salvar spec mesclado: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(filename, byteFile, 0644)
	if err != nil {
		fmt.Printf("Erro ao salvar specfile: %v\n", err)
		os.Exit(1)
	}
	return nil
}

type MergeSpecs struct {
	specA          string
	specB          string
	output         string
	globalDb       bool
	productPathURL string
}

func mergeSpecsCmd() *cobra.Command {

	options := &MergeSpecs{}

	cmd := &cobra.Command{
		Use:     "merge",
		Short:   "Merge OpenAPI specifications",
		Example: "merge -s1 spec1.yaml -s2 spec2.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			mergeSpecsMain(*options)
		},
	}

	cmd.Flags().StringVarP(&options.specA, "target", "t", "", "Path to the first OpenAPI specification file")
	cmd.Flags().StringVarP(&options.specB, "from", "f", "", "Path to the second OpenAPI specification file")
	cmd.Flags().StringVarP(&options.output, "output", "o", "", "Output filename OpenAPI specification file")
	cmd.Flags().BoolVarP(&options.globalDb, "globaldb", "g", false, "Is globalDB?")
	cmd.Flags().StringVarP(&options.productPathURL, "product", "p", "", "productURLPath")

	_ = cmd.MarkFlagRequired("target")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("output")
	_ = cmd.MarkFlagRequired("globaldb")
	_ = cmd.MarkFlagRequired("product")
	return cmd
}
