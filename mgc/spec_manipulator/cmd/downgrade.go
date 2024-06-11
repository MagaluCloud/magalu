package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var downgradeSpecCmd = &cobra.Command{
	Use:    "downgrade",
	Short:  "Downgrade specs from 3.1.x to 3.0.x",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// runPrepare(cmd, args)
		_ = verificarEAtualizarDiretorio(SPEC_DIR)

		currentConfig, err := loadList()

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, v := range currentConfig {
			file := filepath.Join(SPEC_DIR, v.File)
			fileBytes, err := os.ReadFile(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			document, err := libopenapi.NewDocument(fileBytes)
			if err != nil {
				panic(fmt.Sprintf("cannot read document: %e", err))
			}
			docModel, errors := document.BuildV3Model()
			if len(errors) > 0 {
				for i := range errors {
					fmt.Printf("error: %e\n", errors[i])
				}
				panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported", len(errors)))
			}

			if spl := strings.Split(docModel.Model.Version, "."); spl[0] == "3" && spl[1] == "0" {
				fmt.Printf("Document %s is in 3.0.x format\n", v.File)
				continue
			}

			// downgrade to 3.0.x
			docModel.Model.Version = "3.0.3"
			_, document, docModel, errors = document.RenderAndReload()
			if len(errors) > 0 {
				for i := range errors {
					fmt.Printf("error: %e\n", errors[i])
				}
				panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported", len(errors)))
			}

			// Schemas
			for pair := docModel.Model.Components.Schemas.Oldest(); pair != nil; pair = pair.Next() {
				if pair.Value.Schema().Properties == nil {
					continue
				}
				for schema := pair.Value.Schema().Properties.Oldest(); schema != nil; schema = schema.Next() {
					if len(schema.Value.Schema().AnyOf) > 1 {

						hasNull := new(bool)
						*hasNull = false

						example := &yaml.Node{}
						if schema.Value.Schema().Examples != nil {
							example = schema.Value.Schema().Examples[0]
						}

						typeAny := ""
						format := ""
						var valMin *float64
						var valMax *float64

						newExclusiveMin := &base.DynamicValue[bool, float64]{}
						newExclusiveMax := &base.DynamicValue[bool, float64]{}

						var varDefault *yaml.Node

						var oneOf []*base.SchemaProxy
						var items *base.DynamicValue[*base.SchemaProxy, bool]

						for _, anyOf := range schema.Value.Schema().AnyOf {

							if anyOf.Schema().Type == nil {
								continue
							}

							if anyOf.Schema().Type[0] == "null" {
								*hasNull = true
							} else if anyOf.Schema().Type[0] == "integer" || anyOf.Schema().Type[0] == "number" {
								typeAny = anyOf.Schema().Type[0]
								format = anyOf.Schema().Format
								newExclusiveMin = anyOf.Schema().ExclusiveMinimum
								newExclusiveMax = anyOf.Schema().ExclusiveMaximum
								valMin = anyOf.Schema().Minimum
								valMax = anyOf.Schema().Maximum

								if newExclusiveMin != nil && newExclusiveMin.IsA() {
									//Time está com spec invalida, tratar como 3.0.x
									newExclusiveMin = &base.DynamicValue[bool, float64]{
										A: newExclusiveMin.IsA(),
									}
									if valMin == nil {
										valMin = &newExclusiveMin.B
									}
								}

								if newExclusiveMin != nil && newExclusiveMin.IsB() {
									//Time está com spec invalida, tratar como 3.0.x
									newExclusiveMin = &base.DynamicValue[bool, float64]{
										A: newExclusiveMin.IsA(),
									}
									valMin = &newExclusiveMin.B
								}

								if newExclusiveMax != nil && newExclusiveMax.IsA() {
									//Time está com spec invalida, tratar como 3.0.x
									newExclusiveMax = &base.DynamicValue[bool, float64]{
										A: newExclusiveMax.IsA(),
									}
									if valMax == nil {
										valMax = &newExclusiveMax.B
									}
								}

								if newExclusiveMax != nil && newExclusiveMax.IsB() {
									//Time está com spec invalida, tratar como 3.0.x
									newExclusiveMax = &base.DynamicValue[bool, float64]{
										A: newExclusiveMax.IsA(),
									}
									valMax = &newExclusiveMax.B
								}

								varDefault = anyOf.Schema().Default
							} else if anyOf.Schema().Type[0] == "object" {
								typeAny = anyOf.Schema().Type[0]
								format = anyOf.Schema().Format
								oneOf = append(oneOf, anyOf)
							} else if anyOf.Schema().Type[0] == "array" {
								// DISSO
								// "subnets": {
								// 	"anyOf": [
								// 		{
								// 			"items": {
								// 				"type": "string"
								// 			},
								// 			"type": "array"
								// 		},
								// 		{
								// 			"type": "null"
								// 		}
								// 	],
								// 	"title": "Subnets",
								// 	"default": []
								// },

								// FICA ASSIM
								// "tags": {
								// 	"type": "array",
								// 	"description": "List of tags applied to the Kubernetes cluster.",
								// 	"items": {
								// 		"type": "string",
								// 		"nullable": true,
								// 		"description": "Items from the list of tags applied to the Kubernetes cluster.",
								// 		"example": "tag-example"
								// 	}
								// },
								typeAny = anyOf.Schema().Type[0]

								xptos := anyOf.Schema().Items
								if xptos.A.Schema().Type == nil {

									items = &base.DynamicValue[*base.SchemaProxy, bool]{
										A: base.CreateSchemaProxy(
											&base.Schema{
												Type:        []string{"string"},
												Description: "Array",
												Example:     example,
												Nullable:    hasNull,
											},
										),
									}
								} else {
									xpp := xptos.A.Schema()
									fmt.Println(xpp.Type)

									items = &base.DynamicValue[*base.SchemaProxy, bool]{
										A: base.CreateSchemaProxy(
											&base.Schema{
												Type:        []string{xpp.Type[0]},
												Description: xpp.Description,
												Example:     example,
												Nullable:    hasNull,
											},
										),
									}
								}

								fmt.Println(anyOf.Schema().Items)

							} else {
								typeAny = anyOf.Schema().Type[0]
							}
						}
						// Sorry for it!
						// I'm not proud of this code

						propMap := base.CreateSchemaProxy(
							&base.Schema{
								Type:             []string{typeAny},
								Example:          example,
								Nullable:         hasNull,
								Description:      schema.Value.Schema().Description,
								Title:            schema.Value.Schema().Title,
								Format:           format,
								ExclusiveMaximum: newExclusiveMax,
								Maximum:          valMax,
								ExclusiveMinimum: newExclusiveMin,
								Minimum:          valMin,
								Default:          varDefault,
								OneOf:            oneOf,
								Items:            items,
							},
						)

						propMap.BuildSchema()
						schemaToChange, ok := docModel.Model.Components.Schemas.Get(pair.Key)
						if ok {
							pprops := schemaToChange.Schema().Properties
							pprops.Set(schema.Key, propMap)
						}
					}

				}

			}

			//Paths
			for path := docModel.Model.Paths.PathItems.Oldest(); path != nil; path = path.Next() {
				operations := path.Value.GetOperations()
				if operations == nil {
					continue
				}
				for op := operations.Oldest(); op != nil; op = op.Next() {
					if op.Value.Parameters != nil {
						for _, param := range op.Value.Parameters {
							if param.Schema != nil {
								if param.Schema.Schema().AnyOf != nil {
									hasNull := new(bool)
									*hasNull = false

									example := &yaml.Node{}
									if param.Schema.Schema().Examples != nil {
										example = param.Schema.Schema().Examples[0]
									}
									typeAny := ""
									format := ""

									var valMin *float64
									var valMax *float64

									newExclusiveMin := &base.DynamicValue[bool, float64]{}
									newExclusiveMax := &base.DynamicValue[bool, float64]{}

									var varDefault *yaml.Node
									var oneOf []*base.SchemaProxy
									for _, anyOf := range param.Schema.Schema().AnyOf {
										if anyOf.Schema().Type[0] == "null" {
											*hasNull = true
										} else if anyOf.Schema().Type[0] == "integer" || anyOf.Schema().Type[0] == "number" {
											typeAny = anyOf.Schema().Type[0]
											format = anyOf.Schema().Format
											newExclusiveMin = anyOf.Schema().ExclusiveMinimum
											newExclusiveMax = anyOf.Schema().ExclusiveMaximum
											valMin = anyOf.Schema().Minimum
											valMax = anyOf.Schema().Maximum

											if newExclusiveMin != nil && newExclusiveMin.IsA() {
												//Time está com spec invalida, tratar como 3.0.x
												newExclusiveMin = &base.DynamicValue[bool, float64]{
													A: newExclusiveMin.IsA(),
												}
												if valMin == nil {
													valMin = &newExclusiveMin.B
												}
											}

											if newExclusiveMin != nil && newExclusiveMin.IsB() {
												//Time está com spec invalida, tratar como 3.0.x
												newExclusiveMin = &base.DynamicValue[bool, float64]{
													A: newExclusiveMin.IsA(),
												}
												valMin = &newExclusiveMin.B
											}

											if newExclusiveMax != nil && newExclusiveMax.IsA() {
												//Time está com spec invalida, tratar como 3.0.x
												newExclusiveMax = &base.DynamicValue[bool, float64]{
													A: newExclusiveMax.IsA(),
												}
												if valMax == nil {
													valMax = &newExclusiveMax.B
												}
											}

											if newExclusiveMax != nil && newExclusiveMax.IsB() {
												//Time está com spec invalida, tratar como 3.0.x
												newExclusiveMax = &base.DynamicValue[bool, float64]{
													A: newExclusiveMax.IsA(),
												}
												valMax = &newExclusiveMax.B
											}

											varDefault = anyOf.Schema().Default
										} else if anyOf.Schema().Type[0] == "object" {
											typeAny = anyOf.Schema().Type[0]
											format = anyOf.Schema().Format
											oneOf = append(oneOf, anyOf)
										} else if anyOf.Schema().Type[0] == "array" {
											typeAny = anyOf.Schema().Type[0]
											for _, it := range anyOf.Schema().Items.A.GetReferenceNode().Value {
												fmt.Println(it)

											}
										} else {
											typeAny = anyOf.Schema().Type[0]
										}
									}

									// Sorry for it!
									// I'm not proud of this code
									// var schemaToRelace *base.SchemaProxy
									schemaToRelace := base.CreateSchemaProxy(&base.Schema{
										Type:             []string{typeAny},
										Example:          example,
										Nullable:         hasNull,
										Description:      param.Schema.Schema().Description,
										Title:            param.Schema.Schema().Title,
										Format:           format,
										ExclusiveMaximum: newExclusiveMax,
										Maximum:          valMax,
										ExclusiveMinimum: newExclusiveMin,
										Minimum:          valMin,
										Default:          varDefault,
										OneOf:            oneOf,
									})

									_, _ = schemaToRelace.BuildSchema()

									pathKey, ok := docModel.Model.Paths.PathItems.Get(path.Key)
									if ok {
										operation, ok := pathKey.GetOperations().Get(op.Key)
										if ok {
											for opar := range operation.Parameters {
												if operation.Parameters[opar].Schema == param.Schema {
													operation.Parameters[opar].Schema = schemaToRelace
												}
											}
										}
									}

								}

							}
						}
					}
				}
			}

			fileBytes, _, _, errs := document.RenderAndReload()
			if len(errors) > 0 {
				panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
			}

			os.WriteFile(filepath.Join(SPEC_DIR, "conv."+v.File), fileBytes, 0644)
		}
	},
}
