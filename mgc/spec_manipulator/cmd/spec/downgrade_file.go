package spec

import (
	"fmt"
	"os"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"

	validator "github.com/pb33f/libopenapi-validator"

	"github.com/spf13/cobra"
)

func downgradeFile(file string) {
	// extract the file name from the path, example: /home/gfz/go/src/mgc/spec_manipulator/cmd/spec/virtual-machine.openapi.yaml > virtual-machine.openapi
	fileName := strings.Split(file, "/")
	fileName = strings.Split(fileName[len(fileName)-1], ".")

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

	// if spl := strings.Split(docModel.Model.Version, "."); spl[0] == "3" && spl[1] == "0" {
	// 	fmt.Printf("Skipping %s. Already in 3.0.x format\n", file)
	// 	return
	// }

	// downgrade to 3.0.x
	docModel.Model.Version = "3.0.3"

	docModel.Model.Security = nil

	_, document, docModel, errors = document.RenderAndReload()
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}
		panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported", len(errors)))
	}
	fmt.Println("Downgrading " + file)
	// Schemas
	for pair := docModel.Model.Components.Schemas.Oldest(); pair != nil; pair = pair.Next() {
		xchema := pair.Value.Schema()
		*xchema = *PrepareSchema(xchema)
	}

	//Paths
	for path := docModel.Model.Paths.PathItems.Oldest(); path != nil; path = path.Next() {
		operations := path.Value.GetOperations()
		if operations == nil {
			continue
		}

		for op := operations.Oldest(); op != nil; op = op.Next() {
			var newParams []*v3.Parameter
			if op.Value.Parameters != nil {
				for _, param := range op.Value.Parameters {
					if strings.ToLower(param.Name) != "x-tenant-id" {
						xchema := param.Schema.Schema()
						*xchema = *PrepareSchema(xchema)
						newParams = append(newParams, param)
					}
				}
				if len(newParams) > 0 {
					op.Value.Parameters = newParams
				} else {
					op.Value.Parameters = nil
				}
			}

			if len(op.Value.Security) == 0 {
				sec := orderedmap.New[string, []string]()
				if op.Key == "get" {
					sec.Set("OAuth2", []string{fmt.Sprintf("%s:read", fileName[0])})
				} else {
					sec.Set("OAuth2", []string{fmt.Sprintf("%s:write", fileName[0])})
				}

				op.Value.Security = []*base.SecurityRequirement{
					{
						Requirements: sec,
					},
				}
			}
		}

	}

	// if docModel.Model.Components.SecuritySchemes == nil || docModel.Model.Components.SecuritySchemes.Len() == 0 {
	// 	orderedMap := orderedmap.New[string, *v3.SecurityScheme]()
	// 	orderedMap.Set("OAuth2", &v3.SecurityScheme{
	// 		Type: "oauth2",
	// 		Name: "OAuth2",
	// 	})
	// 	docModel.Model.Components.SecuritySchemes = orderedMap
	// }

	_, document, _, errs := document.RenderAndReload()
	if len(errors) > 0 {
		panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
	}
	docValidator, validatorErrs := validator.NewValidator(document)
	if len(validatorErrs) > 0 {
		panic(fmt.Sprintf("cannot create validator: %d errors reported", len(validatorErrs)))
	}

	valid, validationErrs := docValidator.ValidateDocument()

	if !valid {
		for _, e := range validationErrs {
			// 5. Handle the error
			fmt.Printf("Type: %s, Failure: %s\n", e.ValidationType, e.Message)
			fmt.Printf("Fix: %s\n\n", e.HowToFix)
			for _, se := range e.SchemaValidationErrors {
				// fmt.Printf("Schema: %s, Failure: %s\n", se., se.Message)
				// fmt.Printf("Fix: %s\n\n", se.HowToFix)
				fmt.Printf("Reason: %s\n", se.Reason)
				fmt.Printf("Location: %s\n", se.Location)
				fmt.Printf("DeepLocation: %s\n", se.DeepLocation)
				fmt.Printf("AbsoluteLocation: %s\n", se.AbsoluteLocation)
				fmt.Printf("Line: %d\n", se.Line)
				fmt.Printf("Column: %d\n", se.Column)

			}

		}
	}

	fileBytes, _, _, errs = document.RenderAndReload()
	if len(errors) > 0 {
		panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
	}

	_ = os.WriteFile(file, fileBytes, 0644)

}

type DowngradeUniqueSpec struct {
	spec string
}

func DowngradeUniqueSpecCmd() *cobra.Command {
	options := &DowngradeUniqueSpec{}

	cmd := &cobra.Command{
		Use:    "downgrade-unique",
		Short:  "Downgrade specs from 3.1.x to 3.0.x",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			downgradeFile(options.spec)
		},
	}

	cmd.Flags().StringVarP(&options.spec, "spec", "s", "", "Spec to downgrade")

	return cmd
}
