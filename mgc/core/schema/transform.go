package schema

import (
	"errors"
)

// This error notifies the iteration should stop, but it's not returned to the end user
var TransformStop = errors.New("stop walking transformations")

type ConstraintKind string

const (
	ConstraintAllOf ConstraintKind = "allOf"
	ConstraintAnyOf ConstraintKind = "anyOf"
	ConstraintOneOf ConstraintKind = "oneOf"
	ConstraintNot   ConstraintKind = "not"
)

type Transformer[T any] interface {
	// Transform the current value given the schema.
	//
	// If the transformation should stop recursion, then return err == TransformStop.
	// This error is not propagated, that is, callers will receive err == nil.
	//
	// If the returned err == nil, it will proceed to transform Scalar, Array, Object and Constraints.
	//
	// Other errors will abort the transformation with error.
	Transform(schema *Schema, value T) (T, error)

	// Transform a regular scalar
	Scalar(schema *Schema, value T) (T, error)

	// Transform an array.
	//
	// itemSchema may be null if none was provided!
	Array(schema *Schema, itemSchema *Schema, value T) (T, error)

	// Transform an object.
	//
	// Callers usually will call proceed with TransformObjectProperties()
	Object(schema *Schema, value T) (T, error)

	// Transform a series of constraints
	//
	// kind = allOf, anyOf, oneOf and not.
	// note that "not" is a slice with a single element
	//
	// This is only called if the len(schemaRefs) != 0
	//
	// Callers will usually call TransformSchemasArray()
	Constraints(kind ConstraintKind, schemaRefs SchemaRefs, value T) (T, error)
}

// Apply the transformer to the given schema and value
//
// If will first try t.Transform(), if that doesn't return any errors (err == nil),
// then it will proceed to specific types: Scalar(), Object(), Array() or Constraints()
// (allOf, anyOf, oneOf).
//
// This is **NOT** recursive, Transformer implementations are expected to do it
// for Array(), Object() and Constraints(), as needed.
// Object() can use TransformObjectProperties() while
// Constraints() can use TransformSchemasArray().
func Transform[T any](
	t Transformer[T],
	schema *Schema,
	value T,
) (T, error) {
	var err error
	value, err = t.Transform(schema, value)
	if err != nil {
		if err == TransformStop {
			err = nil
		}
		return value, err
	}

	if schema.Type != nil {
		switch {
		case schema.Type.Includes("string"), schema.Type.Includes("number"), schema.Type.Includes("integer"), schema.Type.Includes("boolean"), schema.Type.Includes("null"):
			return t.Scalar(schema, value)

		case schema.Type.Includes("object"):
			return t.Object(schema, value)

		case schema.Type.Includes("array"):
			var itemSchema *Schema
			if schema.Items != nil && schema.Items.Value != nil {
				itemSchema = (*Schema)(schema.Items.Value)
			}
			return t.Array(schema, itemSchema, value)

		}
	}
	notSchemaRefs := SchemaRefs{}
	if schema.Not != nil {
		notSchemaRefs = append(notSchemaRefs, schema.Not)
	}
	sub := map[ConstraintKind]SchemaRefs{
		ConstraintAllOf: schema.AllOf,
		ConstraintAnyOf: schema.AnyOf,
		ConstraintOneOf: schema.OneOf,
		ConstraintNot:   notSchemaRefs,
	}

	for kind, refs := range sub {
		if len(refs) == 0 {
			continue
		}
		value, err = t.Constraints(kind, refs, value)
		if err != nil {
			if err == TransformStop {
				err = nil
			}
			break
		}
	}

	return value, err
}

// Applies the Transformer to every element in the list, in order.
//
// The value returned by one transform is given to the next one,
// the last one is returned
//
// If the transformation should stop looping, then return err == TransformStop.
// This error is not propagated, that is, callers will receive err == nil.
//
// If the returned err == nil, it will proceed to call Transformer using the next schema.
//
// Other errors will abort the transformation with error.
func TransformSchemasArray[T any](
	t Transformer[T],
	schemaRefs SchemaRefs,
	value T,
) (T, error) {
	var err error
	for _, ref := range schemaRefs {
		if ref.Value != nil {
			subSchema := (*Schema)(ref.Value)
			value, err = Transform[T](t, subSchema, value)
			if err != nil {
				if err == TransformStop {
					err = nil
				}
				break
			}
		}
	}

	return value, err
}

// Transform every property of an object given its schema
//
// Only named properties are iterated, additional and pattern properties are not handled.
//
// The value returned by one cbProperty is given to the next one,
// the last one is returned.
//
// If the transformation should stop looping, then return err == TransformStop.
// This error is not propagated, that is, callers will receive err == nil.
//
// If the returned err == nil, it will proceed to call cbProperty using the next property name and schema.
//
// Other errors will abort the transformation with error.
func TransformObjectProperties[T any](
	schema *Schema,
	value T,
	cbProperty func(propName string, propSchema *Schema, value T) (T, error),
) (T, error) {
	var err error

	for k, ref := range schema.Properties {
		propSchema := (*Schema)(ref.Value)
		if propSchema != nil {
			value, err = cbProperty(k, propSchema, value)
			if err != nil {
				if err == TransformStop {
					err = nil
				}
				break
			}
		}
	}

	return value, err
}
