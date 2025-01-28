package schema

import (
	"errors"
	"fmt"
	"reflect"

	"slices"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/jsonpointer"
)

var (
	errorUnsupported     = errors.New("unsupported")
	errorSingleNullChild = errors.New("it makes no sense to have single null child")
)

type simplifyContextItem struct {
	name                string
	schemaRefCow        *COWSchemaRef
	schemaCow           *COWSchema
	schema              *Schema
	pendingSchemaRefCow []*COWSchemaRef
}

type simplifyContext struct {
	trail []*simplifyContextItem
}

func (c *simplifyContext) find(schema *Schema) (i int, item *simplifyContextItem) {
	for i, item = range c.trail {
		if item.schema == schema {
			return
		}
	}
	return -1, nil
}

func (c *simplifyContext) startSchema(spec *simplifyContextItem) bool {
	i, item := c.find(spec.schema)
	if i >= 0 {
		if spec.schemaRefCow == nil {
			panic("programming error: expected schemaRefCow for replicated schema references")
		}
		item.pendingSchemaRefCow = append(item.pendingSchemaRefCow, spec.schemaRefCow)
		return false
	}
	c.trail = append(c.trail, spec)
	return true
}

func (c *simplifyContext) finishSchema(spec *simplifyContextItem) {
	last := len(c.trail) - 1
	if len(c.trail) == 0 || c.trail[last] != spec {
		panic("programming error: trail mismatch")
	}
	path := ""
	for _, item := range c.trail {
		name := item.name
		if name != "" {
			path += "/" + name
		}
	}
	if path == "" {
		path = "/"
	}
	schema := spec.schemaCow.Peek()
	schemaRef := &SchemaRef{Ref: path, Value: (*openapi3.Schema)(schema)}
	for _, schemaRefCow := range spec.pendingSchemaRefCow {
		// Replace so schema is referenced and not copied by its COW
		schemaRefCow.Replace(schemaRef)
	}
	c.trail = c.trail[:last]
}

func hasSchemaRefValue(r *SchemaRef) bool {
	return r != nil && r.Value != nil
}

func getCommonType(children []*SchemaRef) (t string, err error) {
	for i, schemaRef := range children {
		if !hasSchemaRefValue(schemaRef) {
			continue
		}

		schema := schemaRef.Value

		if schema.Type == "" {
			continue
		}

		if t == "" {
			t = schema.Type
			continue
		}

		if t != schema.Type {
			return "", &utils.ChainedError{
				Name: fmt.Sprint(i), Err: &utils.ChainedError{
					Name: "type", Err: &utils.CompareError{A: t, B: schema.Type},
				},
			}
		}
	}

	return t, nil
}

func findNullSchema(schemaRefs []*SchemaRef) int {
	return slices.IndexFunc(schemaRefs, func(schemaRef *SchemaRef) bool {
		return hasSchemaRefValue(schemaRef) && schemaRef.Value.Type == "null"
	})
}

func createAnyOfIfNeeded(a, b *SchemaRef) *SchemaRef {
	if !hasSchemaRefValue(a) {
		return b
	}
	if !hasSchemaRefValue(b) {
		return a
	}
	return &openapi3.SchemaRef{Value: openapi3.NewAnyOfSchema(a.Value, b.Value)}
}

// If all types are the same, set the input to that
//
// NOTE: no child is simplified, that should be done beforehand
func simplifyTypeIfAllMatches(input *COWSchema, children []*SchemaRef) (err error) {
	t := input.Type()
	if t != "" {
		return
	}

	if t, err = getCommonType(children); t == "" {
		return
	}

	input.SetType(t)
	return nil
}

func (c *simplifyContext) simplifySchemaRefs(cowSchemaRefs *utils.COWSliceOfCOW[*SchemaRef, *COWSchemaRef], name string) (err error) {
	_ = cowSchemaRefs.ForEachCOW(func(index int, schemaRef *COWSchemaRef) (run bool) {
		err = c.simplifySchemaRefCOW(schemaRef, fmt.Sprintf("%s/%d", name, index))
		return err == nil
	})
	return
}

func (c *simplifyContext) simplifyNot(input *COWSchema) (err error) {
	if !hasSchemaRefValue(input.Not()) {
		return
	}

	// TODO: is this needed? This complicates a lot the CLI/TF handling.
	// One option would be to create a anyOf() schema that would not match the given schema,
	// but it's quite complex as "not" can be complex on its on (not->anyOf->...)
	return errorUnsupported
}

func (c *simplifyContext) simplifyOneOf(input *COWSchema) (err error) {
	if len(input.OneOf()) == 0 {
		return
	}

	if err = c.simplifySchemaRefs(input.OneOfCOW(), "oneOf"); err != nil {
		return
	}

	children := input.OneOf()
	if len(children) == 1 {
		childSchemaRef := children[0]
		if hasSchemaRefValue(childSchemaRef) {
			childSchema := (*Schema)(childSchemaRef.Value)
			if childSchema.Type == "null" {
				return errorSingleNullChild
			}
			err = mergeIntoParent(input, childSchema)
			if err != nil {
				return &utils.ChainedError{
					Name: "0",
					Err:  fmt.Errorf("could not promote: %w", err),
				}
			}
		}
		input.SetOneOf(nil)
		return c.simplifySchemaCOWInternal(input) // again, so merged lists and items can be simplified with the final values
	}

	if len(children) != 2 {
		_ = simplifyTypeIfAllMatches(input, children)
		return
	}

	// if oneOf is [otherSchema, nullSchema], then promote otherSchema, mark as nullable and drop oneOf
	nullIndex := findNullSchema(children)
	if nullIndex < 0 {
		_ = simplifyTypeIfAllMatches(input, children)
		return
	}

	otherIndex := (nullIndex + 1) % 2
	otherSchemaRef := children[otherIndex]
	if hasSchemaRefValue(otherSchemaRef) {
		otherSchema := (*Schema)(otherSchemaRef.Value)
		err = mergeIntoParent(input, otherSchema)
		if err != nil {
			return &utils.ChainedError{
				Name: fmt.Sprint(otherIndex),
				Err:  fmt.Errorf("could not promote: %w", err),
			}
		} else {
			input.SetNullable(true)
			input.SetOneOf(nil)
		}
		return
	}

	return c.simplifySchemaCOWInternal(input) // again, so merged lists and items can be simplified with the final values
}

func (c *simplifyContext) simplifyAllOf(input *COWSchema) (err error) {
	if len(input.AllOf()) == 0 {
		return
	}

	if err = c.simplifySchemaRefs(input.AllOfCOW(), "allOf"); err != nil {
		return
	}

	children := input.AllOf()

	if input.Type() == "" {
		if err = simplifyTypeIfAllMatches(input, children); err != nil {
			return err
		}
	} else {
		t, err := getCommonType(children)
		if err != nil {
			return err
		}
		if input.Type() != t {
			return &utils.ChainedError{Name: "type", Err: &utils.CompareError{A: t, B: input.Type()}}
		}
	}

	// allOf can be eliminated if we promote all children to parent
	for i, schemaRef := range children {
		if hasSchemaRefValue(schemaRef) {
			schema := (*Schema)(schemaRef.Value)
			if err = mergeIntoParent(input, schema); err != nil {
				return &utils.ChainedError{
					Name: fmt.Sprint(i),
					Err:  fmt.Errorf("could not promote: %w", err),
				}
			}
		}
	}

	input.SetAllOf(nil)

	return c.simplifySchemaCOWInternal(input) // again, so merged lists and items can be simplified with the final values
}

func (c *simplifyContext) simplifyAnyOf(input *COWSchema) (err error) {
	if len(input.AnyOf()) == 0 {
		return
	}

	if err = c.simplifySchemaRefs(input.AnyOfCOW(), "anyOf"); err != nil {
		return
	}

	children := input.AnyOf()

	if len(children) == 1 {
		childSchemaRef := children[0]
		if hasSchemaRefValue(childSchemaRef) {
			childSchema := (*Schema)(childSchemaRef.Value)
			if childSchema.Type == "null" {
				return errorSingleNullChild
			}
			if err = mergeIntoParent(input, childSchema); err != nil {
				return &utils.ChainedError{
					Name: "0",
					Err:  fmt.Errorf("could not promote: %w", err),
				}
			}
		}
		input.SetAnyOf(nil)
		return c.simplifySchemaCOWInternal(input) // again, so merged lists and items can be simplified with the final values
	}

	if input.Type() == "" {
		_ = simplifyTypeIfAllMatches(input, children)
	}

	nullIndex := findNullSchema(children)
	if nullIndex >= 0 {
		input.SetNullable(true)
		switch len(children) {
		case 1:
			input.SetAnyOf(nil)

		case 2:
			// just like oneOf, mark as nullable and promote the other element
			otherIndex := (nullIndex + 1) % 2
			otherSchemaRef := children[otherIndex]
			if hasSchemaRefValue(otherSchemaRef) {
				otherSchema := (*Schema)(otherSchemaRef.Value)
				if err = mergeIntoParent(input, otherSchema); err != nil {
					return &utils.ChainedError{
						Name: fmt.Sprint(otherIndex),
						Err:  fmt.Errorf("could not promote: %w", err),
					}
				}
			}
			input.SetAnyOf(nil)

		default:
			remaining := slices.Delete(children, nullIndex, nullIndex+1)
			input.SetAnyOf(remaining)
		}

		return c.simplifyAnyOf(input) // try again, maybe type can be simplified
	}

	return
}

func (c *simplifyContext) simplifyItems(input *COWSchema) (err error) {
	if !hasSchemaRefValue(input.Items()) {
		return
	}
	return c.simplifySchemaRefCOW(input.ItemsCOW(), "items")
}

func (c *simplifyContext) simplifyProperties(input *COWSchema) (err error) {
	if len(input.Properties()) == 0 {
		return
	}

	_ = input.PropertiesCOW().ForEachCOW(func(k string, propRefCow *COWSchemaRef) (run bool) {
		if err = c.simplifySchemaRefCOW(propRefCow, "properties/"+jsonpointer.Escape(k)); err != nil {
			err = &utils.ChainedError{Name: k, Err: err}
			return false
		}
		return true
	})

	return
}

func (c *simplifyContext) simplifyAdditionalProperties(input *COWSchema) (err error) {
	v := input.AdditionalProperties()
	if !hasSchemaRefValue(v.Schema) {
		return
	}

	schemaRef := NewCOWSchemaRef(v.Schema)
	if err = c.simplifySchemaRefCOW(schemaRef, "additionalProperties"); err != nil {
		return
	}

	v.Schema = schemaRef.Peek()
	input.SetAdditionalProperties(v)
	return
}

func mergeComparable[T comparable](get func() T, set func(T) bool, target T) (err error) {
	var empty T
	if target == empty {
		return
	}

	existing := get()
	if existing == empty {
		_ = set(target) // we use the return bool signature just to be compliant with COW setters
		return
	} else if existing == target {
		return
	}
	return &utils.CompareError{A: existing, B: target}
}

func mergeSchemaRefsSlices(get func() SchemaRefs, set func(SchemaRefs) bool, target SchemaRefs) {
	if len(target) == 0 {
		return
	}

	existing := get()
	if len(existing) == 0 {
		_ = set(target) // we use the return bool signature just to be compliant with COW setters
		return
	}

	merged := make([]*openapi3.SchemaRef, 0, len(existing)+len(target))
	merged = append(merged, existing...)
	merged = append(merged, target...)
	_ = set(merged) // we use the return bool signature just to be compliant with COW setters
}

func mergeSchemaRefs(get func() *SchemaRef, set func(*SchemaRef) bool, target *SchemaRef) {
	if !hasSchemaRefValue(target) {
		return
	}

	existing := get()
	if !hasSchemaRefValue(existing) {
		_ = set(target) // we use the return bool signature just to be compliant with COW setters
	}

	_ = set(createAnyOfIfNeeded(existing, target)) // we use the return bool signature just to be compliant with COW setters
}

func mergeEnum(input *COWSchema, target []any) (err error) {
	if len(target) == 0 {
		return
	}

	if len(input.Enum()) == 0 {
		input.SetEnum(target)
	} else {
		currentEnum := input.Enum()
		for i, targetValue := range target {
			foundIndex := slices.IndexFunc(currentEnum, func(currentValue any) bool {
				return reflect.DeepEqual(targetValue, currentValue)
			})
			if foundIndex < 0 {
				return &utils.CompareError{A: currentEnum, B: targetValue, Message: fmt.Sprintf("missing element %d: %#v", i, targetValue)}
			}
		}
	}

	return nil
}

func mergeRequired(input *COWSchema, target []string) {
	if len(target) == 0 {
		return
	}

	if len(input.Required()) == 0 {
		input.SetRequired(target)
	} else {
		requiredCow := input.RequiredCOW()
		for _, name := range target {
			requiredCow.Add(name)
		}
	}
}

func mergeProperties(input *COWSchema, target openapi3.Schemas) (err error) {
	if len(target) == 0 {
		return
	}

	if len(input.Properties()) == 0 {
		input.SetProperties(target)
		return
	}

	propertiesCow := input.PropertiesCOW()
	for k, propRef := range target {
		if existing, ok := propertiesCow.Get(k); !ok || !hasSchemaRefValue(existing) {
			propertiesCow.Set(k, propRef)
		} else if !equalSchemaRef(propRef, existing) {
			return &utils.ChainedError{
				Name: k,
				Err:  &utils.CompareError{A: existing, B: propRef},
			}
		}
	}

	return
}

func mergeExtensions(parent *COWSchema, target map[string]any) (err error) {
	if len(target) == 0 {
		return
	}

	if len(parent.Extensions()) == 0 {
		parent.SetExtensions(target)
		return
	}

	extensionsCow := parent.ExtensionsCOW()
	for k, v := range target {
		if existing, ok := extensionsCow.Get(k); !ok {
			extensionsCow.Set(k, v)
		} else if !reflect.DeepEqual(v, existing) {
			return &utils.ChainedError{
				Name: k,
				Err:  &utils.CompareError{A: existing, B: v},
			}
		}
	}

	return
}

func mergeAdditionalProperties(input *COWSchema, target openapi3.AdditionalProperties) {
	additionalProperties := input.AdditionalProperties()
	if additionalProperties.Has == nil || !*additionalProperties.Has {
		additionalProperties.Has = target.Has
	}
	if hasSchemaRefValue(target.Schema) {
		additionalProperties.Schema = createAnyOfIfNeeded(additionalProperties.Schema, target.Schema)
	}

	input.SetAdditionalProperties(additionalProperties)
}

// NOTE: this does not simplify parent after it's merged. Do it explicitly in the caller
func mergeIntoParent(parent *COWSchema, child *Schema) (err error) {
	if err = mergeComparable(parent.Type, parent.SetType, child.Type); err != nil {
		return &utils.ChainedError{Name: "type", Err: err}
	}

	if parent.Description() == "" { // this is okay to have diverging, no need to merge
		parent.SetDescription(child.Description)
	}

	if err = mergeEnum(parent, child.Enum); err != nil {
		return &utils.ChainedError{Name: "enum", Err: err}
	}

	if err = mergeComparable(parent.Format, parent.SetFormat, child.Format); err != nil {
		return &utils.ChainedError{Name: "format", Err: err}
	}

	if parent.Default() == nil {
		parent.SetDefault(child.Default)
	}

	if parent.Example() == nil {
		parent.SetExample(child.Example)
	}

	if err = mergeComparable(parent.UniqueItems, parent.SetUniqueItems, child.UniqueItems); err != nil {
		return &utils.ChainedError{Name: "uniqueItems", Err: err}
	}

	if err = mergeComparable(parent.ExclusiveMin, parent.SetExclusiveMin, child.ExclusiveMin); err != nil {
		return &utils.ChainedError{Name: "exclusiveMin", Err: err}
	}

	if err = mergeComparable(parent.ExclusiveMax, parent.SetExclusiveMax, child.ExclusiveMax); err != nil {
		return &utils.ChainedError{Name: "exclusiveMax", Err: err}
	}

	if err = mergeComparable(parent.Nullable, parent.SetNullable, child.Nullable); err != nil {
		return &utils.ChainedError{Name: "nullable", Err: err}
	}

	if err = mergeComparable(parent.ReadOnly, parent.SetReadOnly, child.ReadOnly); err != nil {
		return &utils.ChainedError{Name: "readOnly", Err: err}
	}

	if err = mergeComparable(parent.WriteOnly, parent.SetWriteOnly, child.WriteOnly); err != nil {
		return &utils.ChainedError{Name: "writeOnly", Err: err}
	}

	if err = mergeComparable(parent.AllowEmptyValue, parent.SetAllowEmptyValue, child.AllowEmptyValue); err != nil {
		return &utils.ChainedError{Name: "allowEmptyValue", Err: err}
	}

	if err = mergeComparable(parent.Deprecated, parent.SetDeprecated, child.Deprecated); err != nil {
		return &utils.ChainedError{Name: "deprecated", Err: err}
	}

	if err = mergeComparable(parent.Min, parent.SetMin, child.Min); err != nil {
		return &utils.ChainedError{Name: "minimum", Err: err}
	}

	if err = mergeComparable(parent.Max, parent.SetMax, child.Max); err != nil {
		return &utils.ChainedError{Name: "maximum", Err: err}
	}

	if err = mergeComparable(parent.MultipleOf, parent.SetMultipleOf, child.MultipleOf); err != nil {
		return &utils.ChainedError{Name: "multipleOf", Err: err}
	}

	if err = mergeComparable(parent.MinLength, parent.SetMinLength, child.MinLength); err != nil {
		return &utils.ChainedError{Name: "minLength", Err: err}
	}

	if err = mergeComparable(parent.MaxLength, parent.SetMaxLength, child.MaxLength); err != nil {
		return &utils.ChainedError{Name: "maxLength", Err: err}
	}

	if err = mergeComparable(parent.Pattern, parent.SetPattern, child.Pattern); err != nil {
		return &utils.ChainedError{Name: "pattern", Err: err}
	}

	if err = mergeComparable(parent.MinItems, parent.SetMinItems, child.MinItems); err != nil {
		return &utils.ChainedError{Name: "minItems", Err: err}
	}

	if err = mergeComparable(parent.MaxItems, parent.SetMaxItems, child.MaxItems); err != nil {
		return &utils.ChainedError{Name: "maxItems", Err: err}
	}

	mergeRequired(parent, child.Required)

	if err = mergeComparable(parent.MinProps, parent.SetMinProps, child.MinProps); err != nil {
		return &utils.ChainedError{Name: "minProperties", Err: err}
	}

	if err = mergeComparable(parent.MaxProps, parent.SetMaxProps, child.MaxProps); err != nil {
		return &utils.ChainedError{Name: "maxProperties", Err: err}
	}

	mergeSchemaRefsSlices(parent.OneOf, parent.SetOneOf, child.OneOf)
	mergeSchemaRefsSlices(parent.AnyOf, parent.SetAnyOf, child.AnyOf)
	mergeSchemaRefsSlices(parent.AllOf, parent.SetAllOf, child.AllOf)

	mergeSchemaRefs(parent.Not, parent.SetNot, child.Not)
	mergeSchemaRefs(parent.Items, parent.SetItems, child.Items)

	if err = mergeProperties(parent, child.Properties); err != nil {
		return &utils.ChainedError{Name: "properties", Err: err}
	}

	mergeAdditionalProperties(parent, child.AdditionalProperties)

	if err = mergeExtensions(parent, child.Extensions); err != nil {
		return &utils.ChainedError{Name: "extensions", Err: err}
	}

	return nil
}

// Simplifies the Schema, if needed.
//
// The following simplifications and adjustments are done:
//   - `type` is set if possible (from enum or children schemas)
//   - if `oneOf`/`anyOf` is [otherSchema, nullSchema], then it is removed, the input schema is made nullable
//   - `anyOf` contains nullSchema, it's removed from the list and the input schema is made nullable
//   - `allOf` schemas are merged into the input schema, allOf is then removed
//   - `not` schemas are not supported and will cause an error
func SimplifySchemaCOW(input *COWSchema) (err error) {
	c := &simplifyContext{}
	spec := &simplifyContextItem{
		schemaCow: input,
		schema:    input.Peek(),
	}
	return c.simplifySchemaCOW(spec)
}

func (c *simplifyContext) simplifySchemaCOW(spec *simplifyContextItem) (err error) {
	schema := spec.schema
	if schema == nil {
		return
	}

	if !c.startSchema(spec) {
		return
	}
	err = c.simplifySchemaCOWInternal(spec.schemaCow)
	c.finishSchema(spec)
	return
}

func (c *simplifyContext) simplifySchemaCOWInternal(input *COWSchema) (err error) {
	if input.Type() == "" {
		if len(input.Enum()) > 0 {
			if t, _ := getJsonEnumType(input.Peek()); t != "" { // ignore errors, just don't set the type
				input.SetType(t)
			}
		}
	}

	if err = c.simplifyNot(input); err != nil {
		return &utils.ChainedError{Name: "not", Err: err}
	}

	if err = c.simplifyOneOf(input); err != nil {
		return &utils.ChainedError{Name: "oneOf", Err: err}
	}

	if err = c.simplifyAllOf(input); err != nil {
		return &utils.ChainedError{Name: "allOf", Err: err}
	}

	if err = c.simplifyAnyOf(input); err != nil {
		return &utils.ChainedError{Name: "anyOf", Err: err}
	}

	if err = c.simplifyItems(input); err != nil {
		return &utils.ChainedError{Name: "items", Err: err}
	}

	if err = c.simplifyProperties(input); err != nil {
		return &utils.ChainedError{Name: "properties", Err: err}
	}

	if err = c.simplifyAdditionalProperties(input); err != nil {
		return &utils.ChainedError{Name: "additionalProperties", Err: err}
	}

	return
}

// Helper on top of SimplifySchemaCOW()
//
// The input pointer is NOT modified, a new copy is returned if it was changed
func SimplifySchema(input *Schema) (output *Schema, err error) {
	if input == nil {
		return
	}

	cow := NewCOWSchema(input)
	if err = SimplifySchemaCOW(cow); err != nil {
		return
	}
	return cow.Peek(), nil
}

// Simplifies the Value and make sure the Ref string is unset.
func SimplifySchemaRefCOW(input *COWSchemaRef) (err error) {
	c := &simplifyContext{}
	return c.simplifySchemaRefCOW(input, "")
}

// Special case: when the schema is recursive, the ref should NOT be unset. Otherwise, infinite
// recursion would happen when trying to print it, marshal it, etc.
func shouldUnsetRefWhenSimplifying(input *COWSchemaRef) bool {
	ref := input.Ref()
	return ref != "#/$defs/Schema" &&
		ref != "#/$defs/Schemas" &&
		ref != "#/$defs/SchemaRef" &&
		ref != "#/$defs/SchemaRefs"
}

func (c *simplifyContext) simplifySchemaRefCOW(input *COWSchemaRef, name string) (err error) {
	if input == nil {
		return
	}

	if shouldUnsetRefWhenSimplifying(input) {
		input.UnsetRef() // make sure there is no Ref, even if the value wasn't changed
	}
	spec := &simplifyContextItem{
		name:         name,
		schemaRefCow: input,
		schemaCow:    input.ValueCOW(),
		schema:       input.ValueCOW().Peek(),
	}
	if err = c.simplifySchemaCOW(spec); err != nil {
		return
	}
	return
}
