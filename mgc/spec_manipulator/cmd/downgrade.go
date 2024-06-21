package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"

	validator "github.com/pb33f/libopenapi-validator"

	"github.com/spf13/cobra"
)

func prepareSchema(xchema *base.Schema) *base.Schema {
	newChema := &base.Schema{}
	// Compatible with all versions
	newChema.Not = xchema.Not
	newChema.Title = xchema.Title
	newChema.MultipleOf = xchema.MultipleOf
	newChema.Maximum = xchema.Maximum
	newChema.Minimum = xchema.Minimum
	newChema.MaxLength = xchema.MaxLength
	newChema.MinLength = xchema.MinLength
	newChema.Pattern = xchema.Pattern
	newChema.Format = xchema.Format
	newChema.MaxItems = xchema.MaxItems
	newChema.MinItems = xchema.MinItems
	newChema.UniqueItems = xchema.UniqueItems
	newChema.MaxProperties = xchema.MaxProperties
	newChema.MinProperties = xchema.MinProperties
	newChema.Required = xchema.Required
	newChema.Enum = xchema.Enum
	newChema.Description = xchema.Description

	newChema.Default = xchema.Default
	newChema.Nullable = xchema.Nullable
	newChema.ReadOnly = xchema.ReadOnly
	newChema.WriteOnly = xchema.WriteOnly
	newChema.XML = xchema.XML
	newChema.ExternalDocs = xchema.ExternalDocs
	newChema.Example = xchema.Example
	newChema.Deprecated = xchema.Deprecated
	newChema.Extensions = xchema.Extensions

	newChema.Const = nil

	// In versions 2 and 3.0, this Type is a single value, so array will only ever have one value
	// in version 3.1, Type can be multiple values
	//Type []string
	forceType := true
	if xchema.Type != nil {
		forceType = false
		for _, tp := range xchema.Type {
			if tp == "null" {
				newChema.Nullable = new(bool)
				*newChema.Nullable = true
				continue
			}
			newChema.Type = []string{tp}
		}
	}

	// newChema.AdditionalProperties = xchema.AdditionalProperties
	if xchema.AdditionalProperties != nil {
		if xchema.AdditionalProperties.A != nil {
			trated := xchema.AdditionalProperties.A.Schema()
			var anyof []*base.SchemaProxy
			for _, ao := range trated.AnyOf {
				if ao.Schema().Type[0] == "null" {
					trated.Nullable = new(bool)
					*trated.Nullable = true
					continue
				}
				anyof = append(anyof, ao)
				trated.AnyOf = anyof
			}
			newChema.AdditionalProperties = &base.DynamicValue[*base.SchemaProxy, bool]{
				A: base.CreateSchemaProxy(trated),
				N: 0,
			}
		}
	}

	// 3.1 properties need to be converted to 3.0.x
	if xchema.Properties != nil {
		if newChema.Properties == nil {
			newChema.Properties = orderedmap.New[string, *base.SchemaProxy]()
		}
		propMap := orderedmap.New[string, *base.SchemaProxy]()

		for prop := xchema.Properties.Oldest(); prop != nil; prop = prop.Next() {
			propMap.Set(prop.Key, base.CreateSchemaProxy(prepareSchema(prop.Value.Schema())))
		}

		newChema.Properties = propMap

	}

	// 3.1 only, used to define a dialect for this schema, label is '$schema'.
	//SchemaTypeRef string

	// Schemas are resolved on demand using a SchemaProxy
	//AllOf []*SchemaProxy
	newAllOf := []*base.SchemaProxy{}
	if xchema.AllOf != nil {
		for _, xA := range xchema.AllOf {
			if xA.IsReference() {
				newAllOf = append(newAllOf, xA)
				continue
			}
			for _, xT := range xA.Schema().Type {
				if xT == "null" {
					newChema.Nullable = new(bool)
					*newChema.Nullable = true
					continue
				}
				newAllOf = append(newAllOf, base.CreateSchemaProxy(prepareSchema(xA.Schema())))
			}
		}
	}
	if len(newAllOf) > 0 {
		newChema.AllOf = newAllOf
	}

	// Polymorphic Schemas are only available in version 3+
	//OneOf         []*SchemaProxy
	newOneOf := []*base.SchemaProxy{}
	if xchema.OneOf != nil {
		for _, xO := range xchema.OneOf {
			if xO.IsReference() {
				newOneOf = append(newOneOf, xO)
				continue
			}
			for _, xT := range xO.Schema().Type {
				if xT == "null" {
					newChema.Nullable = new(bool)
					*newChema.Nullable = true
					continue
				}
				newOneOf = append(newOneOf, base.CreateSchemaProxy(prepareSchema(xO.Schema())))
			}
		}
	}
	if len(newOneOf) > 0 {
		newChema.OneOf = newOneOf
	}

	//AnyOf         []*SchemaProxy
	newAnyOf := []*base.SchemaProxy{}
	if xchema.AnyOf != nil {
		for _, xA := range xchema.AnyOf {
			if xA.IsReference() {
				newAnyOf = append(newAnyOf, xA)
				continue
			}

			for _, xT := range xA.Schema().Type {
				if xT == "null" {
					newChema.Nullable = new(bool)
					*newChema.Nullable = true
					continue
				}
				newAnyOf = append(newAnyOf, base.CreateSchemaProxy(prepareSchema(xA.Schema())))
			}

		}
	}
	if len(newAnyOf) > 0 {
		newChema.AnyOf = newAnyOf
	}

	//Discriminator *Discriminator

	// in 3.1 examples can be an array (which is recommended)
	//Examples []*yaml.Node
	if xchema.Examples != nil {
		if len(xchema.Examples) > 0 {
			newChema.Example = xchema.Examples[0]
		}
	}

	// in 3.1 prefixItems provides tuple validation support.
	//PrefixItems []*SchemaProxy
	if xchema.PrefixItems != nil {
		xchema.PrefixItems = nil
	}

	// 3.1 Specific properties
	//Contains          *SchemaProxy
	//MinContains       *int64
	//MaxContains       *int64
	//If                *SchemaProxy
	//Else              *SchemaProxy
	//Then              *SchemaProxy
	//DependentSchemas  *orderedmap.Map[string, *SchemaProxy]
	//PatternProperties *orderedmap.Map[string, *SchemaProxy]
	//PropertyNames     *SchemaProxy
	//UnevaluatedItems  *SchemaProxy

	// in 3.1 UnevaluatedProperties can be a Schema or a boolean
	// https://github.com/pb33f/libopenapi/issues/118
	//UnevaluatedProperties *DynamicValue[*SchemaProxy, bool]

	// in 3.1 Items can be a Schema or a boolean
	//Items *DynamicValue[*SchemaProxy, bool]
	if xchema.Items != nil {
		newChema.Items = &base.DynamicValue[*base.SchemaProxy, bool]{
			A: base.CreateSchemaProxy(prepareSchema(xchema.Items.A.Schema())),
		}

	}

	// 3.1 only, part of the JSON Schema spec provides a way to identify a sub-schema
	//Anchor string

	// In versions 2 and 3.0, this ExclusiveMaximum can only be a boolean.
	// In version 3.1, ExclusiveMaximum is a number.
	//ExclusiveMaximum *DynamicValue[bool, float64]
	if xchema.ExclusiveMaximum != nil && xchema.ExclusiveMaximum.IsB() {
		//assume que é um valor numérico
		newChema.ExclusiveMaximum = &base.DynamicValue[bool, float64]{
			A: true,
			B: xchema.ExclusiveMaximum.B,
			N: 0,
		}

		if newChema.Minimum == nil {
			newChema.Minimum = new(float64)
			*newChema.Minimum = xchema.ExclusiveMaximum.B
		}
	} else if xchema.ExclusiveMaximum != nil && xchema.ExclusiveMaximum.IsA() {
		//assume que é um boolean
		newChema.ExclusiveMaximum = &base.DynamicValue[bool, float64]{
			A: xchema.ExclusiveMaximum.A,
			B: 0,
			N: 0,
		}
	}

	// In versions 2 and 3.0, this ExclusiveMinimum can only be a boolean.
	// In version 3.1, ExclusiveMinimum is a number.
	//ExclusiveMinimum *DynamicValue[bool, float64]
	if xchema.ExclusiveMinimum != nil && xchema.ExclusiveMinimum.IsB() {
		//assume que é um valor numérico
		newChema.ExclusiveMinimum = &base.DynamicValue[bool, float64]{
			A: true,
			N: 0,
		}

		if newChema.Minimum == nil {
			newChema.Minimum = new(float64)
			*newChema.Minimum = xchema.ExclusiveMinimum.B
		}
	} else if xchema.ExclusiveMinimum != nil && xchema.ExclusiveMinimum.IsA() {
		//assume que é um boolean
		newChema.ExclusiveMinimum = &base.DynamicValue[bool, float64]{
			A: xchema.ExclusiveMinimum.A,
			N: 0,
		}
	}
	if newChema.Type != nil && newChema.Type[0] == "array" && newChema.Items == nil {
		newChema.Items = &base.DynamicValue[*base.SchemaProxy, bool]{
			A: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"string"},
			}),
		}
	}

	if forceType && newChema.Type == nil && newChema.AnyOf == nil && newChema.OneOf == nil && newChema.AllOf == nil {
		newChema.Type = []string{"string"}
	}

	return newChema

}

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

			docModel.Model.Security = nil

			_, document, docModel, errors = document.RenderAndReload()
			if len(errors) > 0 {
				for i := range errors {
					fmt.Printf("error: %e\n", errors[i])
				}
				panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported", len(errors)))
			}
			fmt.Println("\nBEGIN FILE:" + file + "\n")
			// Schemas
			for pair := docModel.Model.Components.Schemas.Oldest(); pair != nil; pair = pair.Next() {
				xchema := pair.Value.Schema()
				*xchema = *prepareSchema(xchema)
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
							xchema := param.Schema.Schema()
							*xchema = *prepareSchema(xchema)
						}
					}
				}
			}

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
				}
			}

			fileBytes, _, _, errs = document.RenderAndReload()
			if len(errors) > 0 {
				panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
			}

			_ = os.WriteFile(filepath.Join(SPEC_DIR, "conv."+v.File), fileBytes, 0644)
		}
	},
}
