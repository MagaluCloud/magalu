package schema

import (
	"errors"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
)

func Test_CompareJsonSchemas(t *testing.T) {
	required1 := []string{"a", "b"}
	required2 := []string{"x", "y"}

	enum1 := []any{"a", 1, true}
	enum2 := []any{2, "x", false}

	schemaRef1 := &SchemaRef{Value: openapi3.NewStringSchema()}
	schemaRef2 := &SchemaRef{Value: openapi3.NewIntegerSchema()}
	schemaRef3 := &SchemaRef{Value: openapi3.NewBoolSchema()}

	oneOf1 := SchemaRefs{schemaRef1, schemaRef2}
	oneOf2 := SchemaRefs{schemaRef1, schemaRef3}

	properties1 := openapi3.Schemas{"a": schemaRef1, "b": schemaRef2}
	properties2 := openapi3.Schemas{"a": schemaRef2, "c": schemaRef2}

	floatPtr1, floatPtr2, floatPtr3 := new(float64), new(float64), new(float64)
	*floatPtr1 = 123
	*floatPtr2 = 345
	*floatPtr3 = *floatPtr1

	has1, has2, has3 := new(bool), new(bool), new(bool)
	*has1 = true
	*has2 = false
	*has3 = *has1

	tests := []struct {
		name      string
		a, b      *Schema
		wantError error
	}{
		{
			name: "equal/type",
			a:    NewStringSchema(),
			b:    NewStringSchema(),
		},
		{
			name: "equal/required",
			a:    &Schema{Required: required1},
			b:    &Schema{Required: []string{"b", "a"}},
		},
		{
			name: "equal/enum",
			a:    &Schema{Enum: enum1},
			b:    &Schema{Enum: []any{1, "a", true}},
		},
		{
			name: "equal/enum/a-unset",
			a:    &Schema{},
			b:    &Schema{Enum: []any{1, "a", true}},
		},
		{
			name: "equal/enum/b-unset",
			a:    &Schema{Enum: enum1},
			b:    &Schema{},
		},
		{
			name: "equal/format/equal",
			a:    &Schema{Format: "f"},
			b:    &Schema{Format: "f"},
		},
		{
			name: "equal/format/a-unset",
			a:    &Schema{},
			b:    &Schema{Format: "f"},
		},
		{
			name: "equal/format/b-unset",
			a:    &Schema{Format: "f"},
			b:    &Schema{},
		},
		{
			name: "equal/min/same",
			a:    &Schema{Min: floatPtr1},
			b:    &Schema{Min: floatPtr1},
		},
		{
			name: "equal/min/equal",
			a:    &Schema{Min: floatPtr1},
			b:    &Schema{Min: floatPtr3},
		},
		{
			name: "equal/min/a-unset",
			a:    &Schema{},
			b:    &Schema{Min: floatPtr1},
		},
		{
			name: "equal/min/b-unset",
			a:    &Schema{Min: floatPtr1},
			b:    &Schema{},
		},
		{
			name: "equal/items",
			a:    &Schema{Items: schemaRef1},
			b:    &Schema{Items: &SchemaRef{Value: openapi3.NewStringSchema()}},
		},
		{
			name: "equal/oneOf",
			a:    &Schema{OneOf: oneOf1},
			b:    &Schema{OneOf: SchemaRefs{schemaRef2, schemaRef1}},
		},
		{
			name: "equal/properties",
			a:    &Schema{Properties: properties1},
			b:    &Schema{Properties: openapi3.Schemas{"a": schemaRef1, "b": schemaRef2}},
		},
		{
			name: "equal/additionalProperties/has/same",
			a:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Has: has1}},
			b:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Has: has1}},
		},
		{
			name: "equal/additionalProperties/has/equal",
			a:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Has: has1}},
			b:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Has: has3}},
		},
		{
			name: "equal/additionalProperties/has/a-unset",
			a:    &Schema{},
			b:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Has: has1}},
		},
		{
			name: "equal/additionalProperties/has/b-unset",
			a:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Has: has1}},
			b:    &Schema{},
		},
		{
			name: "equal/additionalProperties/schema/same",
			a:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Schema: schemaRef1}},
			b:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Schema: schemaRef1}},
		},
		{
			name: "equal/additionalProperties/schema/equal",
			a:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Schema: schemaRef1}},
			b:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Schema: &SchemaRef{Value: openapi3.NewStringSchema()}}},
		},
		{
			name: "equal/additionalProperties/schema/a-unset",
			a:    &Schema{},
			b:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Schema: schemaRef1}},
		},
		{
			name: "equal/additionalProperties/schema/b-unset",
			a:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Schema: schemaRef1}},
			b:    &Schema{},
		},
		// {
		// 	name:      "error/type",
		// 	a:         NewStringSchema(),
		// 	b:         NewIntegerSchema(),
		// 	wantError: &utils.ChainedError{Name: "Type", Err: &utils.CompareError{A: "string", B: "integer"}},
		// },
		{
			name:      "error/required",
			a:         &Schema{Required: required1},
			b:         &Schema{Required: required2},
			wantError: &utils.ChainedError{Name: "Required", Err: &utils.CompareError{A: required1, B: required2}},
		},
		{
			name:      "error/enum",
			a:         &Schema{Enum: enum1},
			b:         &Schema{Enum: enum2},
			wantError: &utils.ChainedError{Name: "Enum", Err: &utils.CompareError{A: enum1, B: enum2}},
		},
		{
			name:      "error/format",
			a:         &Schema{Format: "fa"},
			b:         &Schema{Format: "fb"},
			wantError: &utils.ChainedError{Name: "Format", Err: &utils.CompareError{A: "fa", B: "fb"}},
		},
		{
			name:      "error/min",
			a:         &Schema{Min: floatPtr1},
			b:         &Schema{Min: floatPtr2},
			wantError: &utils.ChainedError{Name: "Min", Err: &utils.CompareError{A: floatPtr1, B: floatPtr2}},
		},
		{
			name:      "error/items/a-unset",
			a:         &Schema{},
			b:         &Schema{Items: schemaRef1},
			wantError: &utils.ChainedError{Name: "Items", Err: &utils.CompareError{A: (*SchemaRef)(nil), B: schemaRef1}},
		},
		{
			name:      "error/items/b-unset",
			a:         &Schema{Items: schemaRef1},
			b:         &Schema{},
			wantError: &utils.ChainedError{Name: "Items", Err: &utils.CompareError{A: schemaRef1, B: (*SchemaRef)(nil)}},
		},
		// {
		// 	name: "error/items/different",
		// 	a:    &Schema{Items: schemaRef1},
		// 	b:    &Schema{Items: schemaRef2},
		// 	wantError: &utils.ChainedError{
		// 		Name: "Items",
		// 		Err: &utils.ChainedError{
		// 			Name: "Type", Err: &utils.CompareError{A: schemaRef1.Value.Type.Slice()[0], B: schemaRef2.Value.Type.Slice()[0]},
		// 		},
		// 	},
		// },
		{
			name: "error/oneOf",
			a:    &Schema{OneOf: oneOf1},
			b:    &Schema{OneOf: oneOf2},
			wantError: &utils.ChainedError{
				Name: "OneOf",
				Err:  &utils.CompareError{A: oneOf1, B: oneOf2},
			},
		},
		{
			name: "error/properties",
			a:    &Schema{Properties: properties1},
			b:    &Schema{Properties: properties2},
			wantError: &utils.ChainedError{
				Name: "Properties",
				Err:  &utils.CompareError{A: properties1, B: properties2},
			},
		},
		{
			name: "error/additionalProperties/has",
			a:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Has: has1}},
			b:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Has: has2}},
			wantError: &utils.ChainedError{
				Name: "AdditionalProperties",
				Err: &utils.ChainedError{
					Name: "Has",
					Err:  &utils.CompareError{A: has1, B: has2},
				},
			},
		},
		{
			name: "error/additionalProperties/schema",
			a:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Schema: schemaRef1}},
			b:    &Schema{AdditionalProperties: openapi3.AdditionalProperties{Schema: schemaRef2}},
			wantError: &utils.ChainedError{
				Name: "AdditionalProperties",
				Err: &utils.ChainedError{
					Name: "Schema",
					Err:  &utils.CompareError{A: schemaRef1, B: schemaRef2},
				},
			},
		},
	}

	for _, tc := range tests {
		name := "CompareJsonSchemas/" + tc.name
		t.Run(name, func(t *testing.T) {
			err := CompareJsonSchemas(tc.a, tc.b)
			if tc.wantError == nil && err != nil {
				t.Errorf("%s: unexpected error: %#v\n%s", name, err, err.Error())
			} else if !errors.Is(err, tc.wantError) {
				t.Errorf("%s: expected\n%#v\n%s\ngot:\n%#v\n%s", name, tc.wantError, tc.wantError.Error(), err, err.Error())
			}
		})
	}
}
