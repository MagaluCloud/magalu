package generator

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
)

type generatorTemplateTypeField struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type generatorTemplateTypeKind string

var (
	generatorTemplateTypeKindAlias  = generatorTemplateTypeKind("alias")
	generatorTemplateTypeKindStruct = generatorTemplateTypeKind("struct")
)

type generatorTemplateTypeDefinition struct {
	Name string
	Doc  string
	Kind generatorTemplateTypeKind // alias or struct

	Target string // if alias

	Fields []generatorTemplateTypeField // if struct
}

type generatorTemplateTypes struct {
	Definitions []*generatorTemplateTypeDefinition
	ByName      map[string]*generatorTemplateTypeDefinition
	BySchema    map[*mgcSchemaPkg.Schema]*generatorTemplateTypeDefinition
}

func (t *generatorTemplateTypes) add(def *generatorTemplateTypeDefinition, schema *mgcSchemaPkg.Schema) {
	if t.ByName == nil {
		t.ByName = map[string]*generatorTemplateTypeDefinition{}
	}
	if t.BySchema == nil {
		t.BySchema = map[*mgcSchemaPkg.Schema]*generatorTemplateTypeDefinition{}
	}
	t.ByName[def.Name] = def
	t.BySchema[schema] = def
	t.Definitions = append(t.Definitions, def)
}

func (t *generatorTemplateTypes) addSchemaRef(name string, schemaRef *mgcSchemaPkg.SchemaRef, required bool) (signature string, err error) {
	if schemaRef == nil {
		err = errMissingSchema
		return
	}
	return t.addSchema(name, (*mgcSchemaPkg.Schema)(schemaRef.Value), required)
}

func (t *generatorTemplateTypes) addSchema(name string, schema *mgcSchemaPkg.Schema, required bool) (signature string, err error) {
	signature, err = t.addSchemaInternal(name, schema)
	if err != nil {
		return
	}

	//check if OAD schema is nullable and if the property is not required
	if schema.Nullable || !required {
		signature = "*" + signature
	}
	return
}

// force the name to be added as a type, useful for entry points
func (t *generatorTemplateTypes) addSchemaOrAlias(name string, schema *mgcSchemaPkg.Schema) (signature string, err error) {
	signature, err = t.addSchemaInternal(name, schema)
	if err != nil {
		return
	}
	if signature != name {
		def := &generatorTemplateTypeDefinition{
			Name:   name,
			Kind:   generatorTemplateTypeKindAlias,
			Target: signature,
		}
		t.add(def, schema)
		signature = name
	}

	if schema.Nullable {
		signature = "*" + signature
	}
	return
}

func (t *generatorTemplateTypes) addSchemaInternal(name string, schema *mgcSchemaPkg.Schema) (signature string, err error) {
	if schema == nil {
		err = errMissingSchema
		return
	}

	if mgcSchemaPkg.CheckSimilarJsonSchemas(schema, anySchema) {
		return "any", nil
	}

	switch schema.Type {
	case "string":
		return t.addScalar(name, "string", schema)

	case "integer":
		return t.addScalar(name, "int", schema)

	case "number":
		return t.addScalar(name, "float64", schema)

	case "boolean":
		return t.addScalar(name, "bool", schema)

	case "array":
		return t.addArray(name, schema)

	case "object":
		return t.addObject(name, schema)

	case "":
		if len(schema.OneOf) > 0 {
			return t.addAlternatives(name, "one of", schema, schema.OneOf)
		}
		if len(schema.AnyOf) > 0 {
			return t.addAlternatives(name, "any of", schema, schema.AnyOf)
		}

		return "any", nil // let's behave as any

	default:
		return "", fmt.Errorf("unsupported JSON schema type: %s (%#v)", schema.Type, schema)
	}
}

func (t *generatorTemplateTypes) addScalar(name string, goType string, schema *mgcSchemaPkg.Schema) (signature string, err error) {
	// TODO: enum: alias and emit variables

	signature = goType
	return
}

func (t *generatorTemplateTypes) addArray(name string, schema *mgcSchemaPkg.Schema) (signature string, err error) {
	var def *generatorTemplateTypeDefinition
	if def = t.ByName[name]; def != nil {
		return def.Name, nil
	}
	if def = t.BySchema[schema]; def != nil {
		return def.Name, nil
	}

	var itemSignature string

	// don't create pointer types for array of primitives
	itemSignature, err = t.addSchemaRef(name+"Item", schema.Items, true)
	if err != nil {
		if errors.Is(err, errMissingSchema) {
			// TODO: this should happen in the schema, but it does, just fallback to "any"
			itemSignature = "any"
		} else {
			err = &utils.ChainedError{Name: "items", Err: err}
			return
		}
	}
	def = &generatorTemplateTypeDefinition{
		Name:   name,
		Kind:   generatorTemplateTypeKindAlias,
		Target: "[]" + itemSignature,
	}

	t.add(def, schema)

	return def.Name, nil
}

func (t *generatorTemplateTypes) addAlternatives(name string, doc string, schema *mgcSchemaPkg.Schema, schemaRefs mgcSchemaPkg.SchemaRefs) (signature string, err error) {
	var def *generatorTemplateTypeDefinition
	if def = t.ByName[name]; def != nil {
		return def.Name, nil
	}
	if def = t.BySchema[schema]; def != nil {
		return def.Name, nil
	}

	def = &generatorTemplateTypeDefinition{
		Name:   name,
		Doc:    doc + ": ",
		Kind:   generatorTemplateTypeKindAlias,
		Target: "any",
	}
	t.add(def, schema)

	for i, childRef := range schemaRefs {
		if childRef == nil || childRef.Value == nil {
			continue
		}
		var childSignature string
		k := fmt.Sprint(i)
		childSignature, err = t.addSchemaRef(name+k, childRef, slices.Contains(schema.Required, childRef.Ref))
		if err != nil {
			err = &utils.ChainedError{Name: k, Err: err}
			return
		}

		if i != 0 {
			def.Doc += ", "
		}
		def.Doc += childSignature
	}

	return def.Name, nil
}

func (t *generatorTemplateTypes) addMapOf(name string, schema *mgcSchemaPkg.Schema, propSchemaRef *mgcSchemaPkg.SchemaRef) (signature string, err error) {
	var def *generatorTemplateTypeDefinition
	if def = t.ByName[name]; def != nil {
		return def.Name, nil
	}
	if def = t.BySchema[schema]; def != nil {
		return def.Name, nil
	}

	def = &generatorTemplateTypeDefinition{
		Name: name,
		Doc:  schema.Description,
		Kind: generatorTemplateTypeKindAlias,
	}
	t.add(def, schema)

	childSignature := "any"
	if propSchemaRef != nil && propSchemaRef.Value != nil {
		childSignature, err = t.addSchema(name+"Property", (*mgcSchemaPkg.Schema)(propSchemaRef.Value), slices.Contains(schema.Required, propSchemaRef.Ref))
		if err != nil {
			return
		}
	}
	def.Target = "map[string]" + childSignature

	return def.Name, nil
}

func (t *generatorTemplateTypes) addObject(name string, schema *mgcSchemaPkg.Schema) (signature string, err error) {
	if schema.AdditionalProperties.Has != nil && *schema.AdditionalProperties.Has && schema.AdditionalProperties.Schema != nil {
		return t.addMapOf(name, schema, schema.AdditionalProperties.Schema)
	}

	// TODO: try to reduce number of generated structs by identifying the same fields
	var def *generatorTemplateTypeDefinition
	def = t.ByName[name]

	if def == nil {
		def = t.BySchema[schema]
	}

	fieldCount := len(schema.Properties) + len(schema.AnyOf) + len(schema.OneOf)
	if def == nil {
		def = &generatorTemplateTypeDefinition{
			Name:   name,
			Doc:    schema.Description,
			Kind:   generatorTemplateTypeKindStruct,
			Fields: make([]generatorTemplateTypeField, 0, fieldCount),
		}

		t.add(def, schema)
	}

	keys := make([]string, 0, len(schema.Properties))
	for k := range schema.Properties {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		propRef := schema.Properties[k]
		fieldName := removeChars(strcase.UpperCamelCase(k), ".")

		for _, f := range def.Fields {
			if f.Name == fieldName {
				return def.Name, nil
			}
		}

		var fieldType string
		fieldType, err = t.addSchemaRef(name+fieldName, propRef, slices.Contains(schema.Required, k))
		if err != nil {
			err = &utils.ChainedError{Name: k, Err: err}
			return
		}

		for _, f := range def.Fields {
			if f.Name == fieldName {
				return def.Name, nil
			}
		}

		def.Fields = append(def.Fields, generatorTemplateTypeField{
			Name: fieldName,
			Type: fieldType,
			Tag:  buildFieldTag(k, slices.Contains(schema.Required, k)),
		})

	}

	if len(schema.OneOf) > 0 {
		err = t.addObjectAlternatives(def, name, "one of", schema.OneOf)
		if err != nil {
			err = &utils.ChainedError{Name: "one of", Err: err}
			return
		}
	}

	if len(schema.AnyOf) > 0 {
		err = t.addObjectAlternatives(def, name, "any of", schema.AnyOf)
		if err != nil {
			err = &utils.ChainedError{Name: "any of", Err: err}
			return
		}
	}

	return def.Name, nil
}

func removeChars(s string, chars string) string {
	for _, c := range chars {
		s = strings.ReplaceAll(s, string(c), "")
	}
	return s
}

func (t *generatorTemplateTypes) addObjectAlternatives(objDef *generatorTemplateTypeDefinition, name string, doc string, schemaRefs mgcSchemaPkg.SchemaRefs) (err error) {
	if objDef.Doc != "" {
		objDef.Doc += "\n"
	}
	objDef.Doc += doc + ": "

	schema := &openapi3.Schema{}
	schema.Properties = make(map[string]*mgcSchemaPkg.SchemaRef)

	for i, childRef := range schemaRefs {
		if childRef == nil || childRef.Value == nil {
			continue
		}
		if i == 0 {
			schema = childRef.Value
		}

		for k, v := range childRef.Value.Properties {
			schema.Properties[k] = v
		}

	}

	var childSignature string

	// always create pointer types for alternatives
	childSignature, err = t.addSchema(name, (*mgcSchemaPkg.Schema)(schema), true)
	if err != nil {
		err = &utils.ChainedError{Name: name, Err: err}
		return
	}

	objDef.Doc += childSignature

	return
}

func buildFieldTag(fieldKey string, required bool) (tag string) {
	tag = `json:"` + fieldKey
	if !required {
		tag += ",omitempty"
	}
	tag += `"`
	return
}
