package schema

import (
	"fmt"
	"reflect"

	"slices"

	"maps"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/jsonschema"
)

// NOTE: TODO: should we duplicate this, or find a more generic package?
type Schema openapi3.Schema
type SchemaRef = openapi3.SchemaRef
type SchemaRefs = openapi3.SchemaRefs

func (s *Schema) VisitJSON(value any, opts ...openapi3.SchemaValidationOption) error {
	opts = append(opts, openapi3.MultiErrors())
	return (*openapi3.Schema)(s).VisitJSON(value, opts...)
}

func (s *Schema) Equals(other *Schema) bool {
	return reflect.DeepEqual(s, other)
}

func (s *Schema) IsEmpty() bool {
	return (*openapi3.Schema)(s).IsEmpty()
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

func NewAnySchema() *Schema {
	s := &openapi3.Schema{
		Nullable: true,
		AnyOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{Value: &openapi3.Schema{Type: "null", Nullable: true}},
			&openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
			&openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
			&openapi3.SchemaRef{Value: openapi3.NewFloat64Schema()},
			&openapi3.SchemaRef{Value: openapi3.NewIntegerSchema()},
			&openapi3.SchemaRef{Value: openapi3.NewArraySchema().WithItems(&openapi3.Schema{})},
			&openapi3.SchemaRef{Value: openapi3.NewObjectSchema().WithAnyAdditionalProperties()},
		},
	}

	return (*Schema)(s)
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

func getJsonEnumType(v *Schema) (string, error) {
	types := []string{}
	for _, v := range v.Enum {
		var t string
		switch v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			t = "integer"
		case float32, float64:
			t = "number"
		case string:
			t = "string"
		case bool:
			t = "boolean"
		default:
			return "", fmt.Errorf("unsupported enum value: %+v", v)
		}
		if !slices.Contains(types, t) {
			types = append(types, t)
		}
	}
	if len(types) != 1 {
		return "", fmt.Errorf("must provide values of a single type in a enum, got %+v", types)
	}

	return types[0], nil
}

// Similar schemas are those with the same type and, depending on the type,
// similar properties or restrictions.
func CheckSimilarJsonSchemas(a, b *Schema) bool {
	return CompareJsonSchemas(a, b) == nil
}

func CheckSimilarJsonSchemasRefs(a, b *SchemaRef) bool {
	return CompareJsonSchemaRefs(a, b) == nil
}

func CompareJsonSchemaRefs(a, b *SchemaRef) (err error) {
	if a == b {
		return nil
	} else if a == nil || b == nil {
		return &utils.CompareError{A: a, B: b}
	}
	return CompareJsonSchemas((*Schema)(a.Value), (*Schema)(b.Value))
}

func compareIfNonZero(vA, vB reflect.Value) (err error) {
	if vA.IsZero() || vB.IsZero() {
		return nil
	}
	return utils.ReflectValueDeepEqualCompare(vA, vB)
}

func compareEnum(vA, vB reflect.Value) (err error) {
	if vA.IsZero() || vB.IsZero() {
		return nil
	}
	return utils.ReflectValueUnorderedSliceCompareDeepEqual(vA, vB)
}

func compareSchemaRef(vA, vB reflect.Value) (err error) {
	a, ok := vA.Interface().(*SchemaRef)
	if !ok {
		return &utils.ChainedError{Name: utils.CompareTypeErrorKey, Err: fmt.Errorf("expected *SchemaRef, got %s", vA)}
	}

	b, ok := vB.Interface().(*SchemaRef)
	if !ok {
		return &utils.ChainedError{Name: utils.CompareTypeErrorKey, Err: fmt.Errorf("expected *SchemaRef, got %s", vA)}
	}

	return CompareJsonSchemaRefs(a, b)
}

func compareSchemaRefSlice(vA, vB reflect.Value) (err error) {
	return utils.ReflectValueUnorderedSliceCompare(vA, vB, compareSchemaRef)
}

func compareSchemaRefMap(vA, vB reflect.Value) (err error) {
	a, ok := vA.Interface().(openapi3.Schemas)
	if !ok {
		return &utils.ChainedError{Name: utils.CompareTypeErrorKey, Err: fmt.Errorf("expected *SchemaRef, got %s", vA)}
	}

	b, ok := vB.Interface().(openapi3.Schemas)
	if !ok {
		return &utils.ChainedError{Name: utils.CompareTypeErrorKey, Err: fmt.Errorf("expected *SchemaRef, got %s", vA)}
	}

	if !maps.EqualFunc(a, b, CheckSimilarJsonSchemasRefs) {
		return &utils.CompareError{A: a, B: b}
	}

	return nil
}

func compareAdditionalProperties(vA, vB reflect.Value) (err error) {
	if vA.IsZero() || vB.IsZero() {
		return nil
	}
	if vA.Type() != vB.Type() {
		return &utils.ChainedError{Name: utils.CompareTypeErrorKey, Err: &utils.CompareError{A: vA.Interface(), B: vB.Interface()}}
	}

	a, ok := vA.Interface().(openapi3.AdditionalProperties)
	if !ok {
		return &utils.ChainedError{Name: utils.CompareTypeErrorKey, Err: fmt.Errorf("expected openapi3.AdditionalProperties, got %s", vA)}
	}

	b, ok := vB.Interface().(openapi3.AdditionalProperties)
	if !ok {
		return &utils.ChainedError{Name: utils.CompareTypeErrorKey, Err: fmt.Errorf("expected openapi3.AdditionalProperties, got %s", vA)}
	}

	if a.Has != nil && b.Has != nil && *a.Has != *b.Has {
		return &utils.ChainedError{Name: "Has", Err: &utils.CompareError{A: a.Has, B: b.Has}}
	}

	err = CompareJsonSchemaRefs(a.Schema, b.Schema)
	if err != nil {
		return &utils.ChainedError{Name: "Schema", Err: &utils.CompareError{A: a.Schema, B: b.Schema}}
	}
	return
}

var schemaFieldComparator utils.StructFieldsComparator

func init() {
	// break initialization cycle
	schemaFieldComparator = utils.NewMapStructFieldsComparator[Schema](map[string]utils.ReflectValueCompareFn{
		"Type":                 utils.ReflectValueDeepEqualCompare,
		"OneOf":                compareSchemaRefSlice,
		"AnyOf":                compareSchemaRefSlice,
		"AllOf":                compareSchemaRefSlice,
		"Not":                  compareSchemaRef,
		"Format":               compareIfNonZero,
		"Enum":                 compareEnum,
		"Default":              compareIfNonZero,
		"Nullable":             utils.ReflectValueDeepEqualCompare,
		"ReadOnly":             utils.ReflectValueDeepEqualCompare,
		"WriteOnly":            utils.ReflectValueDeepEqualCompare,
		"AllowEmptyValue":      utils.ReflectValueDeepEqualCompare,
		"Min":                  compareIfNonZero,
		"Max":                  compareIfNonZero,
		"MultipleOf":           compareIfNonZero,
		"MinLength":            compareIfNonZero,
		"MaxLength":            compareIfNonZero,
		"Pattern":              compareIfNonZero,
		"MinItems":             compareIfNonZero,
		"MaxItems":             compareIfNonZero,
		"Items":                compareSchemaRef,
		"Required":             utils.ReflectValueUnorderedSliceCompareDeepEqual,
		"Properties":           compareSchemaRefMap,
		"MinProps":             compareIfNonZero,
		"MaxProps":             compareIfNonZero,
		"AdditionalProperties": compareAdditionalProperties,
	})
}

func CompareJsonSchemas(a, b *Schema) (err error) {
	if a == b {
		return nil
	} else if a == nil || b == nil {
		return &utils.CompareError{A: a, B: b}
	}

	return utils.StructFieldsCompare(a, b, schemaFieldComparator)
}

func SchemaFromType[T any]() (*Schema, error) {
	t := new(T)
	tp := reflect.TypeOf(t).Elem()
	kind := tp.Kind()
	if tp.Name() == "" && kind == reflect.Interface {
		return NewAnySchema(), nil
	}

	s, err := ToCoreSchema(schemaReflector.Reflect(t))
	if err != nil {
		return nil, fmt.Errorf("unable to create JSON Schema for type '%T': %w", t, err)
	}

	isArray := kind == reflect.Array || kind == reflect.Slice

	// schemaReflector seems to lose the fact that it's an array, so we bring that back
	if isArray && s.Type == "object" {
		arrSchema := NewArraySchema(s)
		s = arrSchema
	}

	return s, nil
}

var schemaReflector *jsonschema.Reflector

func init() {
	schemaReflector = &jsonschema.Reflector{
		DoNotReference: false,
	}
}
