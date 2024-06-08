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

						newMinValue := &base.DynamicValue[bool, float64]{}
						newMaxValue := &base.DynamicValue[bool, float64]{}

						var varDefault *yaml.Node

						for _, anyOf := range schema.Value.Schema().AnyOf {
							if anyOf.Schema().Type == nil {
								continue
							}
							if anyOf.Schema().Type[0] == "null" {
								*hasNull = true
							} else {
								typeAny = anyOf.Schema().Type[0]
								format = anyOf.Schema().Format

								if anyOf.Schema().ExclusiveMinimum != nil {
									newMinValue = &base.DynamicValue[bool, float64]{
										A: anyOf.Schema().ExclusiveMinimum != nil,
									}
									valMin = &anyOf.Schema().ExclusiveMinimum.B
								}
								if anyOf.Schema().ExclusiveMaximum != nil {
									newMaxValue = &base.DynamicValue[bool, float64]{
										A: anyOf.Schema().ExclusiveMaximum != nil,
									}
									valMax = &anyOf.Schema().ExclusiveMaximum.B
								}
								varDefault = anyOf.Schema().Default
							}
						}

						propMap := base.CreateSchemaProxy(
							&base.Schema{
								Type:             []string{typeAny},
								Example:          example,
								Nullable:         hasNull,
								Description:      schema.Value.Schema().Description,
								Title:            schema.Value.Schema().Title,
								Format:           format,
								ExclusiveMaximum: newMaxValue,
								Maximum:          valMax,
								ExclusiveMinimum: newMinValue,
								Minimum:          valMin,
								Default:          varDefault,
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

									newMinValue := &base.DynamicValue[bool, float64]{}
									newMaxValue := &base.DynamicValue[bool, float64]{}

									var varDefault *yaml.Node

									for _, anyOf := range param.Schema.Schema().AnyOf {
										if anyOf.Schema().Type[0] == "null" {
											*hasNull = true
										} else {
											typeAny = anyOf.Schema().Type[0]
											format = anyOf.Schema().Format
											if anyOf.Schema().ExclusiveMinimum != nil {
												newMinValue = &base.DynamicValue[bool, float64]{
													A: anyOf.Schema().ExclusiveMinimum != nil,
												}
												valMin = &anyOf.Schema().ExclusiveMinimum.B
											}
											if anyOf.Schema().ExclusiveMaximum != nil {
												newMaxValue = &base.DynamicValue[bool, float64]{
													A: anyOf.Schema().ExclusiveMaximum != nil,
												}
												valMax = &anyOf.Schema().ExclusiveMaximum.B
											}
											varDefault = anyOf.Schema().Default
										}
									}

									schemaToRelace := base.CreateSchemaProxy(&base.Schema{
										Type:             []string{typeAny},
										Example:          example,
										Nullable:         hasNull,
										Description:      param.Schema.Schema().Description,
										Title:            param.Schema.Schema().Title,
										Format:           format,
										ExclusiveMaximum: newMaxValue,
										Maximum:          valMax,
										ExclusiveMinimum: newMinValue,
										Minimum:          valMin,
										Default:          varDefault,
									})

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
