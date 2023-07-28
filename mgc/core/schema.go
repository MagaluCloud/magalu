package core

import "github.com/getkin/kin-openapi/openapi3"

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

func NewBooleanSchema() *Schema {
	return (*Schema)(openapi3.NewBoolSchema())
}

func NewNullSchema() *Schema {
	return &Schema{
		Type: "null",
	}
}

func NewArraySchema(item *Schema) *Schema {
	return &Schema{
		Type:  "array",
		Items: &openapi3.SchemaRef{Value: (*openapi3.Schema)(item)},
	}
}

func SetDefault(schema *Schema, value any) *Schema {
	schema.Default = value
	return schema
}

func SetDescription(schema *Schema, description string) *Schema {
	schema.Description = description
	return schema
}

// UnmarshalJSON sets Schema to a copy of data.
func (schema *Schema) UnmarshalJSON(data []byte) error {
	return (*openapi3.Schema)(schema).UnmarshalJSON(data)
}

// MarshalJSON returns the JSON encoding of Schema.
func (schema Schema) MarshalJSON() ([]byte, error) {
	return openapi3.Schema(schema).MarshalJSON()
}
