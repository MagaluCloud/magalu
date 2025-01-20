package schema

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/jsonpointer"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
	"go.uber.org/zap"
	"golang.org/x/exp/constraints"
)

var convertLogger = utils.NewLazyLoader(func() *zap.SugaredLogger {
	return logger().Named("convert")
})

// TODO: if not required and struct, make nullable

func ToCoreSchema(s *jsonschema.Schema) (schema *Schema, err error) {
	if s == nil {
		return nil, fmt.Errorf("invalid jsonschema.Schema passed to 'toCoreSchema' function")
	}

	convertLogger().Debugw("ToCoreSchema called: will convert schema", "schema", s)

	refResolver := newRefResolver(s)
	oapiSchemaRef := convertJsonSchemaToOpenAPISchemaRef(s, refResolver)
	err = refResolver.resolvePending()
	if err != nil {
		return
	}

	schemaRefCow := NewCOWSchemaRef(oapiSchemaRef)
	err = SimplifySchemaRefCOW(schemaRefCow)
	if err != nil {
		return
	}
	schema = schemaRefCow.Value()
	return
}

func isJsonSchemaSchemaNullable(input *jsonschema.Schema) bool {
	for _, sub := range input.OneOf {
		if sub != nil && sub.Type == "null" {
			return true
		}
	}
	return false
}

func convertJsonSchemaNumberToOpenAPIPointer[T constraints.Integer | constraints.Float](v int) (r *T) {
	if v == 0 {
		return nil
	}
	r = new(T)
	*r = T(v)
	return
}

func addExtensions(output *openapi3.Schema, name string, value any) {
	if output.Extensions == nil {
		output.Extensions = map[string]interface{}{}
	}
	output.Extensions[name] = value
}

func convertJsonSchemaToOpenAPISchema(input *jsonschema.Schema, refResolver *refResolver) (output *openapi3.Schema) {
	convertLogger().Debugw("starting conversion from 'jsonschema.Schema' to 'kin-openapi.Schema'", "jsonschema", input)

	if input == nil {
		convertLogger().Debugw("returning nil, since input was nil")
		return nil
	}

	// We used to MarshalJSON() from jsonschema.Schema and UnmarshalJSON() into openapi3.Schema, but
	// jsonschema's MarshalJSON() will return "true" for empty schema and this is not handled by openapi3.Schema's UnmarshalJSON()
	// Then do it manually.

	if input == jsonschema.TrueSchema {
		convertLogger().Debugw("returning 'any' schema")
		return (*openapi3.Schema)(NewAnySchema())
	}

	additionalProperties := openapi3.AdditionalProperties{}
	if input.AdditionalProperties != nil && input.AdditionalProperties != jsonschema.FalseSchema {
		convertLogger().Debugw("will convert and add additional properties")
		has := true
		additionalProperties.Has = &has
		additionalProperties.Schema = convertJsonSchemaToOpenAPISchemaRef(input.AdditionalProperties, refResolver)
	}
	if len(input.PatternProperties) > 0 && additionalProperties.Has == nil {
		convertLogger().Debugw("will convert and add additional properties")
		has := true
		additionalProperties.Has = &has
		additionalProperties.Schema = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				AnyOf: convertJsonSchemaToOpenAPISchemaMapToSlice(input.PatternProperties, refResolver),
			},
		}
	}

	output = &openapi3.Schema{
		OneOf:        convertJsonSchemaToOpenAPISchemaSlice(input.OneOf, refResolver),
		AnyOf:        convertJsonSchemaToOpenAPISchemaSlice(input.AnyOf, refResolver),
		AllOf:        convertJsonSchemaToOpenAPISchemaSlice(input.AllOf, refResolver),
		Not:          convertJsonSchemaToOpenAPISchemaRef(input.Not, refResolver),
		Type:         input.Type,
		Title:        input.Title,
		Format:       input.Format,
		Description:  input.Description,
		Enum:         input.Enum,
		Default:      input.Default,
		UniqueItems:  input.UniqueItems,
		ExclusiveMin: input.ExclusiveMinimum,
		ExclusiveMax: input.ExclusiveMaximum,
		Nullable:     isJsonSchemaSchemaNullable(input),
		ReadOnly:     input.ReadOnly,
		WriteOnly:    input.WriteOnly,
		// Does not exist: AllowEmptyValue:      input.AllowEmptyValue,
		Deprecated:           input.Deprecated,
		Min:                  convertJsonSchemaNumberToOpenAPIPointer[float64](input.Minimum),
		Max:                  convertJsonSchemaNumberToOpenAPIPointer[float64](input.Maximum),
		MultipleOf:           convertJsonSchemaNumberToOpenAPIPointer[float64](input.MultipleOf),
		MinLength:            uint64(input.MinLength),
		MaxLength:            convertJsonSchemaNumberToOpenAPIPointer[uint64](input.MaxLength),
		Pattern:              input.Pattern,
		MinItems:             uint64(input.MinItems),
		MaxItems:             convertJsonSchemaNumberToOpenAPIPointer[uint64](input.MaxItems),
		Items:                convertJsonSchemaToOpenAPISchemaRef(input.Items, refResolver),
		Required:             input.Required,
		Properties:           convertJsonSchemaToOpenAPISchemaMap(input.Properties, refResolver),
		MinProps:             uint64(input.MinProperties),
		MaxProps:             convertJsonSchemaNumberToOpenAPIPointer[uint64](input.MaxProperties),
		AdditionalProperties: additionalProperties,
		// Does not exist: Discriminator:        input.Discriminator,
	}

	if len(input.Examples) > 0 {
		convertLogger().Debugw("will add add examples")
		output.Example = input.Examples[0]
	}

	if input.ContentMediaType != "" {
		convertLogger().Debugw("will add x-contentMediaType extension")
		addExtensions(output, "x-contentMediaType", input.ContentMediaType)
	}
	if input.ContentEncoding != "" {
		convertLogger().Debugw("will add x-contentEncoding extension")
		addExtensions(output, "x-contentEncoding", input.ContentEncoding)
	}
	if input.ContentSchema != nil {
		convertLogger().Debugw("will add x-contentSchema extension")
		addExtensions(output, "x-contentSchema", input.ContentSchema)
	}

	convertLogger().Debugw("finished converting 'jsonschema.Schema' to 'kin-openapi.Schema'", "jsonschema", input, "kin-openapi", output)

	return
}

func convertJsonSchemaToOpenAPISchemaRef(input *jsonschema.Schema, refResolver *refResolver) (output *openapi3.SchemaRef) {
	if input == nil {
		return nil
	}

	convertLogger().Debugw("starting conversion from 'jsonschema.Schema' to 'kin-openapi.SchemaRef'", "ref", input.Ref)

	s := convertJsonSchemaToOpenAPISchema(input, refResolver)
	if s == nil {
		return nil
	}

	ref := &openapi3.SchemaRef{}
	if input.Ref != "" {
		convertLogger().Debugw("adding ref to be resolved later", "ref", input.Ref)
		ref = refResolver.add(input.Ref)
	}

	if isSchemaEmpty(s) {
		if input.Ref != "" {
			convertLogger().Debugw("returning ref without resolution, as it was empty", "ref", input.Ref)
			return ref
		} else {
			convertLogger().Debugw("changing empty output to 'any' schema")
			s = (*openapi3.Schema)(NewAnySchema())
		}
	}

	ref.Value = s
	convertLogger().Debugw("finished converting 'jsonschema.Schema' to 'kin-openapi.SchemaRef'", "input", input, "outputRef", ref.Ref, "outputValue", ref.Value)
	return ref
}

func convertJsonSchemaToOpenAPISchemaSlice(input []*jsonschema.Schema, refResolver *refResolver) (output []*openapi3.SchemaRef) {
	if len(input) == 0 {
		return nil
	}

	convertLogger().Debugw("starting conversion from '[]jsonschema.Schema' to '[]kin-openapi.SchemaRef'", "jsonschemas", input)

	output = make([]*openapi3.SchemaRef, len(input))
	for i, value := range input {
		convertLogger().Debugw("will convert schema in array", "index", i)
		output[i] = convertJsonSchemaToOpenAPISchemaRef(value, refResolver)
	}
	return
}

func convertJsonSchemaToOpenAPISchemaMap(input *orderedmap.OrderedMap, refResolver *refResolver) (output map[string]*openapi3.SchemaRef) {
	if input == nil {
		return nil
	}
	values := input.Values()
	if len(values) == 0 {
		return nil
	}
	convertLogger().Debugw("starting conversion from 'map[string]jsonschema.Schema' to 'map[string]kin-openapi.SchemaRef'", "jsonschemas", input)
	output = make(map[string]*openapi3.SchemaRef, len(values))
	for key, value := range values {
		output[key] = convertJsonSchemaToOpenAPISchemaRef(value.(*jsonschema.Schema), refResolver)
	}
	return
}

func convertJsonSchemaToOpenAPISchemaMapToSlice(input map[string]*jsonschema.Schema, refResolver *refResolver) (output []*openapi3.SchemaRef) {
	if len(input) == 0 {
		return nil
	}
	output = make([]*openapi3.SchemaRef, 0, len(input))
	for _, value := range input {
		output = append(output, convertJsonSchemaToOpenAPISchemaRef(value, refResolver))
	}
	return
}

type refResolver struct {
	doc             *jsonschema.Schema
	jsonSchemaCache map[string]*jsonschema.Schema
	oapiSchemaCache map[string]*openapi3.Schema
	pending         []*openapi3.SchemaRef
}

func newRefResolver(doc *jsonschema.Schema) *refResolver {
	// Def overrides for kin-openapi Schema of Schema
	// schemaSchema, schemasSchema, schemaRefSchema, schemaRefsSchema := newjsonschemaSchemaSchema()
	// schemaSchema, schemasSchema, schemaRefSchema, schemaRefsSchema := newkinopenapiSchemaSchema()
	resolver := &refResolver{
		doc:             doc,
		jsonSchemaCache: map[string]*jsonschema.Schema{},
		oapiSchemaCache: map[string]*openapi3.Schema{},
	}

	hasSchemaOfSchema := false

	for defs := range doc.Definitions {
		if defs == "Schema" || defs == "Schemas" || defs == "SchemaRef" || defs == "SchemaRefs" {
			hasSchemaOfSchema = true
			break
		}
	}

	// If we:
	// - Generate a reflected Schema via 'jsonschema' of a struct that has a field of type '*core.Schema'
	// - Convert the reflect Schema to OpenAPI Schema (core.Schema)
	// - Try to validate a variable of type '*core.Schema'
	// We will get errors because the converted version of the Schema doesn't allow for 'null' to be passed
	// to the 'AdditionalProperties' field. So we manually create a Schema to represent the '*core.Schema'
	// if necessary, and then add it to the cache. If not necessary, don't even bother
	if !hasSchemaOfSchema {
		return resolver
	}

	schemaOfSchemaRef := openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
		"ref": openapi3.NewStringSchema(),
		// "value": openapi3.NewSchemaSchema(), // Will be inserted later
	}).WithNullable()

	schemaOfSchemas := openapi3.NewObjectSchema().WithAdditionalProperties((*openapi3.Schema)(schemaOfSchemaRef))
	schemaOfSchemaRefs := openapi3.NewArraySchema().WithItems((*openapi3.Schema)(schemaOfSchemaRef))

	refOfSchemaOfSchemaRefs := resolver.add("#/$defs/SchemaRefs")
	refOfSchemaOfSchemaRef := resolver.add("#/$defs/SchemaRef")
	refOfSchemaOfSchemas := resolver.add("#/$defs/Schemas")

	schemaOfSchema := &openapi3.Schema{
		Type:     "object",
		Nullable: true,
		Properties: openapi3.Schemas{
			"extensions":  openapi3.NewSchemaRef("", openapi3.NewObjectSchema().WithAnyAdditionalProperties()),
			"oneOf":       refOfSchemaOfSchemaRefs,
			"anyOf":       refOfSchemaOfSchemaRefs,
			"allOf":       refOfSchemaOfSchemaRefs,
			"not":         refOfSchemaOfSchemaRef,
			"type":        openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
			"title":       openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
			"format":      openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
			"description": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
			"enum":        openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems((*openapi3.Schema)(NewAnySchema()))),
			"default":     openapi3.NewSchemaRef("", (*openapi3.Schema)(NewAnySchema())),
			"example":     openapi3.NewSchemaRef("", (*openapi3.Schema)(NewAnySchema())),
			"externalDocs": openapi3.NewSchemaRef("", openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
				"extensions":  openapi3.NewObjectSchema().WithAnyAdditionalProperties(),
				"description": openapi3.NewStringSchema(),
				"url":         openapi3.NewStringSchema(),
			})),
			"uniqueItems":     openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
			"exclusiveMin":    openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
			"exclusiveMax":    openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
			"nullable":        openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
			"readOnly":        openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
			"writeOnly":       openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
			"allowEmptyValue": openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
			"deprecated":      openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
			"xml": openapi3.NewSchemaRef("", openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
				"extensions": openapi3.NewObjectSchema().WithAnyAdditionalProperties(),
				"name":       openapi3.NewStringSchema(),
				"namespace":  openapi3.NewStringSchema(),
				"prefix":     openapi3.NewStringSchema(),
				"attribute":  openapi3.NewBoolSchema(),
				"wrapped":    openapi3.NewBoolSchema(),
			})),
			"min":        openapi3.NewSchemaRef("", openapi3.NewFloat64Schema()),
			"max":        openapi3.NewSchemaRef("", openapi3.NewFloat64Schema()),
			"multipleOf": openapi3.NewSchemaRef("", openapi3.NewFloat64Schema()),
			"minLength":  openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
			"maxLength":  openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
			"pattern":    openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
			// "compiledPattern": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
			"minItems":   openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
			"maxItems":   openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
			"items":      refOfSchemaOfSchemaRef,
			"required":   openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
			"properties": refOfSchemaOfSchemas,
			"minProps":   openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
			"maxProps":   openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
			"additionalProperties": openapi3.NewSchemaRef("", &openapi3.Schema{
				Type: "object",
				Properties: openapi3.Schemas{
					"has":    openapi3.NewSchemaRef("", openapi3.NewBoolSchema().WithNullable()),
					"schema": refOfSchemaOfSchemaRef,
				},
			}),
			"discriminator": openapi3.NewSchemaRef("", openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
				"extensions":   openapi3.NewObjectSchema().WithAnyAdditionalProperties(),
				"propertyName": openapi3.NewStringSchema(),
				"mapping":      openapi3.NewObjectSchema().WithAdditionalProperties(openapi3.NewStringSchema()),
			})),
		},
	}

	schemaOfSchemaRef.Properties["value"] = openapi3.NewSchemaRef("", (*openapi3.Schema)(schemaOfSchema).WithNullable())

	resolver.oapiSchemaCache["/$defs/Schema"] = schemaOfSchema
	resolver.oapiSchemaCache["/$defs/Schemas"] = schemaOfSchemas
	resolver.oapiSchemaCache["/$defs/SchemaRef"] = schemaOfSchemaRef
	resolver.oapiSchemaCache["/$defs/SchemaRefs"] = schemaOfSchemaRefs

	return resolver
}

func (r *refResolver) add(ref string) (schemaRef *openapi3.SchemaRef) {
	schemaRef = &openapi3.SchemaRef{Ref: ref}
	r.pending = append(r.pending, schemaRef)
	return
}

const maxResolvePendingIterations = 10

func (r *refResolver) resolvePending() error {
	for i := 0; i < maxResolvePendingIterations && len(r.pending) > 0; i++ {
		convertLogger().Debugw("will start resolving all unresolved refs", "iteration", i, "count", len(r.pending))
		pending := r.pending
		r.pending = nil
		for _, schemaRef := range pending {
			convertLogger().Debugw("will resolve ref", "ref", schemaRef.Ref)
			schema, err := r.resolveOpenAPI(schemaRef.Ref)
			if err != nil {
				return err
			}

			if schemaRef.Value == nil || reflect.DeepEqual(schemaRef.Value, schema) {
				convertLogger().Debugw("resolved ref, setting value", "ref", schemaRef.Ref, "value", schema)
				schemaRef.Value = schema
			} else {
				if isSchemaEmpty(schema) {
					convertLogger().Debugw("resolved ref is empty, ignoring", "ref", schemaRef.Ref)
					continue
				} else if isSchemaEmpty(schemaRef.Value) {
					convertLogger().Debugw("resolved ref, setting value", "ref", schemaRef.Ref, "value", schema)
					schemaRef.Value = schema
				} else {
					// We need to merge the existing value with the ref because `jsonschema` will insert the Struct Tags customizations
					// only in the non-ref value, so that the ref isn't changed everywhere, and only in the struct that used the tags.
					// Use AllOf and then simplify, should never fail
					convertLogger().Debugw("will merge current value with resolved as allOf", "ref", schemaRef.Ref, "current", schemaRef.Value, "resolved", schema)
					allOf := NewAllOfSchema((*Schema)(schema), (*Schema)(schemaRef.Value))
					allOf, err = SimplifySchema(allOf)
					if err != nil {
						return err
					}

					convertLogger().Debugw("resolved ref, setting value", "ref", schemaRef.Ref, "value", allOf)
					schemaRef.Value = (*openapi3.Schema)(allOf)
				}
			}
		}
	}

	convertLogger().Debugw("finished resolving refs")

	if len(r.pending) > 0 {
		return fmt.Errorf(
			"could not finish resolving pending %d references in %d iterations, likely a unsolvable cycle happened",
			len(r.pending),
			maxResolvePendingIterations,
		)
	}

	return nil
}

func (r *refResolver) resolveJsonSchema(ref string) (jsonSchema *jsonschema.Schema, err error) {
	path, _ := strings.CutPrefix(ref, "#")
	if jsonSchema = r.jsonSchemaCache[path]; jsonSchema != nil {
		convertLogger().Debugw("returning schema from 'jsonschema' cache, it was already resolved", "ref", ref)
		return
	}

	jp, err := jsonpointer.New(path)
	if err != nil {
		return
	}

	v, _, err := jp.Get(r.doc)
	if err != nil {
		return
	}

	jsonSchema, ok := v.(*jsonschema.Schema)
	if !ok {
		err = fmt.Errorf("invalid type: expected *jsonschema.Schema, got %T", v)
		return
	}

	if jsonSchema.Ref != "" {
		jsonSchema, err = r.resolveJsonSchema(jsonSchema.Ref)
		if err != nil {
			return
		}
	}

	convertLogger().Debugw("adding ref to jsonschema cache", "ref", ref)
	r.jsonSchemaCache[path] = jsonSchema

	return
}

func (r *refResolver) resolveOpenAPI(ref string) (oapiSchema *openapi3.Schema, err error) {
	path, _ := strings.CutPrefix(ref, "#")
	if oapiSchema = r.oapiSchemaCache[path]; oapiSchema != nil {
		convertLogger().Debugw("returning schema from 'kin-openapi' cache, it was already resolved", "ref", ref)
		return
	}

	jsonSchema, err := r.resolveJsonSchema(path)
	if err != nil {
		return
	}

	oapiSchema = convertJsonSchemaToOpenAPISchema(jsonSchema, r)
	convertLogger().Debugw("adding ref to 'kin-openapi' cache", "ref", ref)
	r.oapiSchemaCache[path] = oapiSchema

	return
}

func isSchemaEmpty(s *openapi3.Schema) bool {
	return s.IsEmpty() && s.Description == ""
}
