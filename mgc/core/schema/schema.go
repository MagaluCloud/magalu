package schema

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// NOTE: TODO: should we duplicate this, or find a more generic package?
type Schema openapi3.Schema

func (s *Schema) VisitJSON(value any, opts ...openapi3.SchemaValidationOption) error {
	opts = append(opts, openapi3.MultiErrors())
	return (*openapi3.Schema)(s).VisitJSON(value, opts...)
}

// UnmarshalJSON sets Schema to a copy of data.
func (schema *Schema) UnmarshalJSON(data []byte) error {
	return (*openapi3.Schema)(schema).UnmarshalJSON(data)
}

// MarshalJSON returns the JSON encoding of Schema.
func (schema Schema) MarshalJSON() ([]byte, error) {
	return openapi3.Schema(schema).MarshalJSON()
}

func NewSchemaRef(ref string, schema *Schema) *openapi3.SchemaRef {
	return openapi3.NewSchemaRef(ref, (*openapi3.Schema)(schema))
}

func NewObjectSchema(properties map[string]*Schema, required []string) *Schema {
	hasAdditionalProperties := false

	p := openapi3.Schemas{}
	for k, v := range properties {
		p[k] = &openapi3.SchemaRef{Value: (*openapi3.Schema)(v)}
	}

	return &Schema{
		Type:                 "object",
		AdditionalProperties: openapi3.AdditionalProperties{Has: &hasAdditionalProperties},
		Properties:           p,
		Required:             required,
	}
}

func NewStringSchema() *Schema {
	return (*Schema)(openapi3.NewStringSchema())
}

func NewNumberSchema() *Schema {
	return (*Schema)(openapi3.NewFloat64Schema())
}

func NewIntegerSchema() *Schema {
	return (*Schema)(openapi3.NewInt64Schema())
}

func NewBooleanSchema() *Schema {
	return (*Schema)(openapi3.NewBoolSchema())
}

func NewNullSchema() *Schema {
	return &Schema{
		Type:     "null",
		Nullable: true,
	}
}

func NewArraySchema(item *Schema) *Schema {
	return &Schema{
		Type:  "array",
		Items: &openapi3.SchemaRef{Value: (*openapi3.Schema)(item)},
	}
}

func NewAnyOfSchema(anyOfs ...*Schema) *Schema {
	anyOfsCast := make([]*openapi3.Schema, 0, len(anyOfs))
	for _, v := range anyOfs {
		anyOfsCast = append(anyOfsCast, (*openapi3.Schema)(v))
	}
	return (*Schema)(openapi3.NewAnyOfSchema(anyOfsCast...))
}
func NewOneOfSchema(oneOfs ...*Schema) *Schema {
	anyOfsCast := make([]*openapi3.Schema, 0, len(oneOfs))
	for _, v := range oneOfs {
		anyOfsCast = append(anyOfsCast, (*openapi3.Schema)(v))
	}
	return (*Schema)(openapi3.NewOneOfSchema(anyOfsCast...))
}
func NewAllOfSchema(allOfs ...*Schema) *Schema {
	anyOfsCast := make([]*openapi3.Schema, 0, len(allOfs))
	for _, v := range allOfs {
		anyOfsCast = append(anyOfsCast, (*openapi3.Schema)(v))
	}
	return (*Schema)(openapi3.NewAllOfSchema(anyOfsCast...))
}

func SetDefault(schema *Schema, value any) *Schema {
	schema.Default = value
	return schema
}

func SetDescription(schema *Schema, description string) *Schema {
	schema.Description = description
	return schema
}

// *Recursively* checks if schema can be nullable. Function supports `nullable`
// fields and `type: 'null'` fields, including if included in anyOf, allOf and oneOf
// properties. It also checks for type refs to check if they are nullable, alas why recursive.
func IsSchemaNullable(schema *Schema) bool {
	// Object is nullable
	if schema.Nullable || schema.Type == "null" {
		return true
	}
	// Object has nullable type in type list
	possibleRefs := []openapi3.SchemaRefs{schema.AnyOf, schema.OneOf, schema.AllOf}
	for _, refs := range possibleRefs {
		for _, typeRef := range refs {
			// ! Recursive call
			if IsSchemaNullable((*Schema)(typeRef.Value)) {
				return true
			}
		}
	}
	return false
}
