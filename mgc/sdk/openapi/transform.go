package openapi

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
	"go.uber.org/zap"
	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
)

type transformSpec struct {
	Type string `json:"type" yaml:"type"`
	// See more about the 'remain' directive here: https://pkg.go.dev/github.com/mitchellh/mapstructure#hdr-Remainder_Values
	Parameters map[string]any `json:",remain"` // nolint
	Schema     *openapi3.Schema
}

type transformRegExpSpec struct {
	// Regular Expression as per https://pkg.go.dev/regexp#Compile
	Pattern string `json:"pattern" yaml:"pattern"`
	// Replacement Template as per https://pkg.go.dev/regexp#Regexp.Expand
	Replacement string `json:"replacement" yaml:"replacement"`
}

func transformRegExp(params map[string]any, s string) (result string, err error) {
	spec, err := utils.DecodeNewValue[transformRegExpSpec](params)
	if err != nil {
		return s, fmt.Errorf("invalid regexp parameters: %w", err)
	}
	if len(spec.Pattern) == 0 {
		return s, fmt.Errorf("invalid regexp parameters: missing pattern")
	}
	re, err := regexp.Compile(spec.Pattern)
	if err != nil {
		return s, fmt.Errorf("invalid regexp pattern %q: %w", spec.Pattern, err)
	}

	b := []byte{}
	for _, submatches := range re.FindAllStringSubmatchIndex(s, -1) {
		b = re.ExpandString(b, spec.Replacement, s, submatches)
	}
	return string(b), nil
}

type transformTranslateSpecItem struct {
	From any `json:"from" yaml:"from"`
	To   any `json:"to" yaml:"to"`
}

type transformTranslateSpec struct {
	Translations []transformTranslateSpecItem `json:"translations" yaml:"translations"`
	AllowMissing bool                         `json:"allowMissing,omitempty" yaml:"allowMissing,omitempty"`
}

func transformTranslate(params map[string]any, value any) (result any, err error) {
	spec, err := utils.DecodeNewValue[transformTranslateSpec](params)
	if err != nil {
		return value, fmt.Errorf("invalid translation parameters: %w", err)
	}
	if len(spec.Translations) == 0 {
		return value, fmt.Errorf("invalid translation parameters: missing translations")
	}
	for _, item := range spec.Translations {
		if reflect.DeepEqual(item.From, value) {
			return item.To, nil
		}
	}
	if spec.AllowMissing {
		return value, nil
	}
	return value, fmt.Errorf("translation not found: %+v", value)
}

func doTransformValue(spec *transformSpec, value any) (any, error) {
	switch spec.Type {
	case "uppercase", "upper-case", "upper":
		if s, ok := value.(string); ok {
			return strings.ToUpper(s), nil
		}
	case "lowercase", "lower-case", "lower":
		if s, ok := value.(string); ok {
			return strings.ToLower(s), nil
		}
	case "kebabcase", "kebab-case", "kebab":
		if s, ok := value.(string); ok {
			return strcase.KebabCase(s), nil
		}
	case "snakecase", "snake-case", "snake":
		if s, ok := value.(string); ok {
			return strcase.SnakeCase(s), nil
		}
	case "pascal", "pascalcase", "pascal-case", "upper-camel":
		if s, ok := value.(string); ok {
			return strcase.UpperCamelCase(s), nil
		}
	case "camel", "camelcase", "camel-case", "lower-camel":
		if s, ok := value.(string); ok {
			return strcase.LowerCamelCase(s), nil
		}
	case "regexp", "regexp-replace":
		if s, ok := value.(string); ok {
			return transformRegExp(spec.Parameters, s)
		}
	case "translate":
		return transformTranslate(spec.Parameters, value)
	}

	return value, nil
}

func doTransformSchema(spec *transformSpec, value *mgcSchemaPkg.COWSchema) (*mgcSchemaPkg.COWSchema, error) {
	if spec.Schema != nil {
		value.Replace((*mgcSchemaPkg.Schema)(spec.Schema))
		return value, nil
	}
	switch spec.Type {
	default:
		return value, nil
	case "translate":
		return transformTranslateSchema(spec.Parameters, value)
	}
}

func doTransformsToValue(logger *zap.SugaredLogger, specs []*transformSpec, value any) (result any, err error) {
	result = value
	for _, spec := range specs {
		result, err = doTransformValue(spec, result)
		if err != nil {
			logger.Debugw("transformation attempt failed", "value", value, "type", spec.Type)
			return
		}
	}
	logger.Debugw("transformed value", "input", value, "output", result)
	return
}

func doTransformsToSchema(logger *zap.SugaredLogger, specs []*transformSpec, value *mgcSchemaPkg.COWSchema) (result *mgcSchemaPkg.COWSchema, err error) {
	result = value
	for _, spec := range specs {
		result, err = doTransformSchema(spec, result)
		if err != nil {
			logger.Debugw("transformation attempt failed", "value", value, "type", spec.Type)
			return
		}
	}
	logger.Debugw("transformed schema", "input", value.Peek(), "output", result.Peek())
	return
}

func getTransformKey(extensionPrefix *string) string {
	if extensionPrefix == nil || *extensionPrefix == "" {
		return ""
	}
	return *extensionPrefix + "-transforms"
}

func newTransformSpecFromString(s string) *transformSpec {
	if len(s) == 0 {
		return nil
	}
	return &transformSpec{Type: s}
}

func newTransformSpecFromMap(m map[string]any) *transformSpec {
	if len(m) == 0 {
		return nil
	}
	spec, err := utils.DecodeNewValue[transformSpec](m)
	if err != nil || len(spec.Type) == 0 {
		return nil
	}
	return spec
}

// spec must be string or map
func newTransformSpec(spec any) *transformSpec {
	if s, ok := spec.(string); ok {
		return newTransformSpecFromString(s)
	} else if m, ok := spec.(map[string]any); ok {
		return newTransformSpecFromMap(m)
	}
	return nil
}

func newTransformSpecSlice(sl []any) []*transformSpec {
	ret := make([]*transformSpec, 0, len(sl))
	for _, spec := range sl {
		if ts := newTransformSpec(spec); ts != nil {
			ret = append(ret, ts)
		}
	}
	if len(ret) == 0 {
		return nil
	}
	return ret
}

func getTransformationSpecs(extensions map[string]any, transformationKey string) []*transformSpec {
	if spec, ok := extensions[transformationKey]; !ok {
		return nil
	} else if sl, ok := spec.([]any); ok {
		return newTransformSpecSlice(sl)
	} else if ts := newTransformSpec(spec); ts != nil {
		return []*transformSpec{ts}
	} else {
		return nil
	}
}

// The returned function does NOT and should NOT alter the value that was passed by it
// (maps, for example, when passed as input, won't be altered, a new copy will be made)
func createTransform[T any](logger *zap.SugaredLogger, schema *core.Schema, extensionPrefix *string) (func(value T) (T, error), *core.Schema, error) {
	transformationKey := getTransformKey(extensionPrefix)
	if transformationKey == "" {
		return nil, schema, nil
	}

	needs, err := needsTransformation(schema, transformationKey)
	if err != nil {
		return nil, schema, err
	}
	if !needs {
		return nil, schema, nil
	}

	transformedSchema, err := transformSchema(logger, schema, transformationKey, schema)
	if err != nil {
		return nil, schema, err
	}

	return func(value T) (converted T, err error) {
		r, err := transformValue(logger, schema, transformationKey, value)
		if err != nil {
			return
		}
		converted, ok := r.(T)
		if !ok {
			err = fmt.Errorf("invalid conversion result, expected %T, got %+v", converted, r)
			return
		}
		return
	}, transformedSchema, nil
}

// Common pattern that checks existing specs, if they exist then call transformSpecs(),
// otherwise process Arrays and Objects.
//
// Scalars are passed thru while Constraints() are recursively processed.
type commonSchemaTransformer[T any] struct {
	tKey                 string
	transformSpecs       func(specs []*transformSpec, value T) (T, error)
	transformArray       func(t mgcSchemaPkg.Transformer[T], schema *core.Schema, itemSchema *core.Schema, value T) (T, error)
	transformObject      func(t mgcSchemaPkg.Transformer[T], schema *core.Schema, value T) (T, error)
	transformConstraints func(t mgcSchemaPkg.Transformer[T], kind mgcSchemaPkg.ConstraintKind, schemaRefs mgcSchemaPkg.SchemaRefs, value T) (T, error)
}

func (t *commonSchemaTransformer[T]) Transform(schema *core.Schema, value T) (T, error) {
	specs := getTransformationSpecs(schema.Extensions, t.tKey)
	var err error
	if len(specs) > 0 {
		value, err = t.transformSpecs(specs, value)
		if err == nil {
			err = mgcSchemaPkg.TransformStop
		}
	}
	return value, err
}

func (t *commonSchemaTransformer[T]) Scalar(schema *core.Schema, value T) (T, error) {
	return value, nil
}

func (t *commonSchemaTransformer[T]) Array(schema *core.Schema, itemSchema *core.Schema, value T) (T, error) {
	if itemSchema == nil {
		return value, nil
	}
	return t.transformArray(t, schema, itemSchema, value)
}

func (t *commonSchemaTransformer[T]) Constraints(kind mgcSchemaPkg.ConstraintKind, schemaRefs mgcSchemaPkg.SchemaRefs, value T) (T, error) {
	return t.transformConstraints(t, kind, schemaRefs, value)
}

func (t *commonSchemaTransformer[T]) Object(schema *core.Schema, value T) (T, error) {
	return t.transformObject(t, schema, value)
}

var _ mgcSchemaPkg.Transformer[any] = (*commonSchemaTransformer[any])(nil)

// Recursively checks whenever the given schema needs transformation
func needsTransformation(schema *core.Schema, transformationKey string) (bool, error) {
	t := &commonSchemaTransformer[bool]{
		tKey:                 transformationKey,
		transformSpecs:       func(specs []*transformSpec, value bool) (bool, error) { return true, nil },
		transformArray:       transformArrayNeedsTransformation,
		transformObject:      transformObjectNeedsTransformation,
		transformConstraints: transformConstraintsNeedsTransformation,
	}
	return mgcSchemaPkg.Transform[bool](t, schema, false)
}

func transformArrayNeedsTransformation(t mgcSchemaPkg.Transformer[bool], schema *core.Schema, itemSchema *core.Schema, value bool) (bool, error) {
	if itemSchema == nil {
		return value, nil
	}
	return mgcSchemaPkg.Transform(t, itemSchema, value)
}

func transformObjectNeedsTransformation(t mgcSchemaPkg.Transformer[bool], schema *core.Schema, value bool) (bool, error) {
	return mgcSchemaPkg.TransformObjectProperties(schema, value, func(propName string, propSchema *core.Schema, value bool) (bool, error) {
		value, err := mgcSchemaPkg.Transform(t, propSchema, value)
		if err != nil {
			return value, err
		}
		if value {
			return true, mgcSchemaPkg.TransformStop
		}
		return false, nil
	})
}

func transformConstraintsNeedsTransformation(t mgcSchemaPkg.Transformer[bool], kind mgcSchemaPkg.ConstraintKind, schemaRefs mgcSchemaPkg.SchemaRefs, value bool) (bool, error) {
	value, err := mgcSchemaPkg.TransformSchemasArray(t, schemaRefs, value)
	if err != nil {
		return value, err
	}
	if value {
		return true, mgcSchemaPkg.TransformStop
	}
	return false, nil

}

// Recursively transforms the value based on the schema that may contain transformations
// If the schema doesn't contain any transformation, then the value is unchanged
func transformValue(logger *zap.SugaredLogger, schema *core.Schema, transformationKey string, value any) (any, error) {
	t := &commonSchemaTransformer[any]{
		tKey:                 transformationKey,
		transformSpecs:       func(specs []*transformSpec, value any) (any, error) { return doTransformsToValue(logger, specs, value) },
		transformArray:       transformArrayValue,
		transformObject:      transformObjectValue,
		transformConstraints: transformConstraintsValue,
	}
	return mgcSchemaPkg.Transform[any](t, schema, value)
}

func transformArrayValue(t mgcSchemaPkg.Transformer[any], schema *core.Schema, itemSchema *core.Schema, value any) (any, error) {
	valueSlice, ok := value.([]any)
	if !ok {
		return value, fmt.Errorf("expected []any, got %T %#v", value, value)
	}

	cs := utils.NewCOWSliceFunc(valueSlice, utils.IsSameValueOrPointer)
	for i, itemValue := range valueSlice {
		convertedValue, err := mgcSchemaPkg.Transform(t, itemSchema, itemValue)
		if err != nil {
			return value, err
		}
		cs.Set(i, convertedValue)
	}

	valueSlice, _ = cs.Release()
	return valueSlice, nil
}

func transformObjectValue(t mgcSchemaPkg.Transformer[any], schema *core.Schema, value any) (any, error) {
	valueMap, ok := value.(map[string]any)
	if !ok {
		return value, fmt.Errorf("expected map[string]any, got %T %#v", value, value)
	}
	cm, err := mgcSchemaPkg.TransformObjectProperties(
		schema,
		utils.NewCOWMapFunc(valueMap, utils.IsSameValueOrPointer),
		func(propName string, propSchema *core.Schema, cm *utils.COWMap[string, any],
		) (*utils.COWMap[string, any], error) {
			propValue, ok := valueMap[propName]
			if !ok {
				return cm, nil
			}

			convertedFieldValue, err := mgcSchemaPkg.Transform(t, propSchema, propValue)
			if err != nil {
				return cm, err
			}
			cm.Set(propName, convertedFieldValue)
			return cm, nil
		},
	)
	if err != nil {
		return value, err
	}

	valueMap, _ = cm.Release()
	return valueMap, nil
}

func transformConstraintsValue(t mgcSchemaPkg.Transformer[any], kind mgcSchemaPkg.ConstraintKind, schemaRefs mgcSchemaPkg.SchemaRefs, value any) (any, error) {
	// TODO: handle kind properly, see https://swagger.io/docs/specification/data-models/oneof-anyof-allof-not/
	return mgcSchemaPkg.TransformSchemasArray(t, schemaRefs, value)
}

func transformSchema(logger *zap.SugaredLogger, schema *core.Schema, transformationKey string, value *core.Schema) (*core.Schema, error) {
	t := &commonSchemaTransformer[*mgcSchemaPkg.COWSchema]{
		tKey: transformationKey,
		transformSpecs: func(specs []*transformSpec, value *mgcSchemaPkg.COWSchema) (*mgcSchemaPkg.COWSchema, error) {
			return doTransformsToSchema(logger, specs, value)
		},
		transformArray:       transformArraySchema,
		transformObject:      transformObjectSchema,
		transformConstraints: transformConstraintsSchema,
	}
	cowSchema := mgcSchemaPkg.NewCOWSchema(value)
	cowSchema, err := mgcSchemaPkg.Transform[*mgcSchemaPkg.COWSchema](t, schema, cowSchema)
	if err != nil {
		return value, err
	}
	return cowSchema.Peek(), nil
}

func transformArraySchema(t mgcSchemaPkg.Transformer[*mgcSchemaPkg.COWSchema], schema *core.Schema, itemSchema *core.Schema, value *mgcSchemaPkg.COWSchema) (*mgcSchemaPkg.COWSchema, error) {
	itemsCow := value.ItemsCOW().ValueCOW()
	_, err := mgcSchemaPkg.Transform(t, itemSchema, itemsCow)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func transformObjectSchema(t mgcSchemaPkg.Transformer[*mgcSchemaPkg.COWSchema], schema *core.Schema, value *mgcSchemaPkg.COWSchema) (*mgcSchemaPkg.COWSchema, error) {
	_, err := mgcSchemaPkg.TransformObjectProperties(
		schema,
		value.PropertiesCOW(),
		func(propName string, propSchema *core.Schema, propertiesCow *utils.COWMapOfCOW[string, *mgcSchemaPkg.SchemaRef, *mgcSchemaPkg.COWSchemaRef],
		) (*utils.COWMapOfCOW[string, *mgcSchemaPkg.SchemaRef, *mgcSchemaPkg.COWSchemaRef], error) {
			propSchemaCow, ok := propertiesCow.GetCOW(propName)
			if !ok {
				return nil, fmt.Errorf("schema missing property %q", propName) // this should never happen
			}

			_, err := mgcSchemaPkg.Transform(t, propSchema, propSchemaCow.ValueCOW())
			if err != nil {
				return nil, err
			}
			return propertiesCow, nil
		},
	)
	if err != nil {
		return value, err
	}
	return value, nil
}

func transformConstraintsSchema(t mgcSchemaPkg.Transformer[*mgcSchemaPkg.COWSchema], kind mgcSchemaPkg.ConstraintKind, schemaRefs mgcSchemaPkg.SchemaRefs, value *mgcSchemaPkg.COWSchema) (result *mgcSchemaPkg.COWSchema, err error) {
	result = value

	if kind == mgcSchemaPkg.ConstraintNot {
		_, err = mgcSchemaPkg.Transform(t, (*mgcSchemaPkg.Schema)(schemaRefs[0].Value), value.NotCOW().ValueCOW())
		return
	}

	var constraintCow *utils.COWSliceOfCOW[*mgcSchemaPkg.SchemaRef, *mgcSchemaPkg.COWSchemaRef]
	switch kind {
	case mgcSchemaPkg.ConstraintAllOf:
		constraintCow = value.AllOfCOW()
	case mgcSchemaPkg.ConstraintAnyOf:
		constraintCow = value.AnyOfCOW()
	case mgcSchemaPkg.ConstraintOneOf:
		constraintCow = value.OneOfCOW()
	default:
		return value, fmt.Errorf("unknown constraint kind: %q", kind)
	}

	constraintCow.ForEachCOW(func(i int, cowRef *mgcSchemaPkg.COWSchemaRef) (run bool) {
		itemSchema := cowRef.Peek()
		if itemSchema == nil {
			return true
		}

		_, err = mgcSchemaPkg.Transform(t, (*mgcSchemaPkg.Schema)(itemSchema.Value), cowRef.ValueCOW())
		return err == nil
	})

	return
}

func reverseTranslate(spec *transformTranslateSpec, value any) (any, error) {
	for _, item := range spec.Translations {
		if reflect.DeepEqual(item.To, value) {
			return item.From, nil
		}
	}
	if spec.AllowMissing {
		return value, nil
	}
	return value, fmt.Errorf("translation not found: %#v", value)
}

func transformTranslateSchema(params map[string]any, schema *mgcSchemaPkg.COWSchema) (result *mgcSchemaPkg.COWSchema, err error) {
	if schema.Default() == nil && len(schema.Enum()) == 0 {
		return schema, nil
	}

	spec, err := utils.DecodeNewValue[transformTranslateSpec](params)
	if err != nil {
		return schema, fmt.Errorf("invalid translation parameters: %w", err)
	}
	if len(spec.Translations) == 0 {
		return schema, fmt.Errorf("invalid translation parameters: missing translations")
	}

	result = schema

	if schema.Default() != nil {
		var schemaDefault any
		schemaDefault, err = reverseTranslate(spec, schema.Default())
		if err != nil {
			return
		}
		schema.SetDefault(schemaDefault)
	}

	enumCow := schema.EnumCOW()
	enumCow.ForEach(func(i int, value any) (run bool) {
		var translatedEnum any
		translatedEnum, err = reverseTranslate(spec, value)
		if err != nil {
			return false
		}
		enumCow.Set(i, translatedEnum)
		return true
	})

	return
}
