package openapi

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/stoewer/go-strcase"
	"magalu.cloud/core"
	schemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
)

type transformSpec struct {
	Type string `json:"type" yaml:"type"`
	// See more about the 'remain' directive here: https://pkg.go.dev/github.com/mitchellh/mapstructure#hdr-Remainder_Values
	Parameters map[string]any `json:",remain"` // nolint
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

func doTransformsToValue(specs []*transformSpec, value any) (result any, err error) {
	result = value
	for _, spec := range specs {
		result, err = doTransformValue(spec, result)
		if err != nil {
			logger().Debugf("attempted to transform %#v but failed. Transformation type was %s", value, spec.Type)
			return
		}
	}
	logger().Debugf("transformed %#v into %#v", value, result)
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
func createTransform[T any](schema *core.Schema, extensionPrefix *string) func(value T) (T, error) {
	transformationKey := getTransformKey(extensionPrefix)
	if transformationKey == "" {
		return nil
	}

	needs, err := needsTransformation(schema, transformationKey)
	if err != nil {
		// TODO: pass error along
		logger().Errorw("failed to check if needs transformation", "err", err)
		return nil
	}
	if !needs {
		return nil
	}

	return func(value T) (converted T, err error) {
		r, err := transformValue(schema, transformationKey, value)
		if err != nil {
			return
		}
		converted, ok := r.(T)
		if !ok {
			err = fmt.Errorf("invalid conversion result, expected %T, got %+v", converted, r)
			return
		}
		return
	}
}

// Common pattern that checks existing specs, if they exist then call transformSpecs(),
// otherwise process Arrays and Objects.
//
// Scalars are passed thru while Constraints() are recursively processed.
type commonSchemaTransformer[T any] struct {
	tKey                 string
	transformSpecs       func(specs []*transformSpec, value T) (T, error)
	transformArray       func(t schemaPkg.Transformer[T], schema *core.Schema, itemSchema *core.Schema, value T) (T, error)
	transformObject      func(t schemaPkg.Transformer[T], schema *core.Schema, value T) (T, error)
	transformConstraints func(t schemaPkg.Transformer[T], kind schemaPkg.ConstraintKind, schemaRefs schemaPkg.SchemaRefs, value T) (T, error)
}

func (t *commonSchemaTransformer[T]) Transform(schema *core.Schema, value T) (T, error) {
	specs := getTransformationSpecs(schema.Extensions, t.tKey)
	var err error
	if len(specs) > 0 {
		value, err = t.transformSpecs(specs, value)
		if err == nil {
			err = schemaPkg.TransformStop
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

func (t *commonSchemaTransformer[T]) Constraints(kind schemaPkg.ConstraintKind, schemaRefs schemaPkg.SchemaRefs, value T) (T, error) {
	return t.transformConstraints(t, kind, schemaRefs, value)
}

func (t *commonSchemaTransformer[T]) Object(schema *core.Schema, value T) (T, error) {
	return t.transformObject(t, schema, value)
}

var _ schemaPkg.Transformer[any] = (*commonSchemaTransformer[any])(nil)

// Recursively checks whenever the given schema needs transformation
func needsTransformation(schema *core.Schema, transformationKey string) (bool, error) {
	t := &commonSchemaTransformer[bool]{
		tKey:                 transformationKey,
		transformSpecs:       func(specs []*transformSpec, value bool) (bool, error) { return true, nil },
		transformArray:       transformArrayNeedsTransformation,
		transformObject:      transformObjectNeedsTransformation,
		transformConstraints: transformConstraintsNeedsTransformation,
	}
	return schemaPkg.Transform[bool](t, schema, false)
}

func transformArrayNeedsTransformation(t schemaPkg.Transformer[bool], schema *core.Schema, itemSchema *core.Schema, value bool) (bool, error) {
	if itemSchema == nil {
		return value, nil
	}
	return schemaPkg.Transform(t, itemSchema, value)
}

func transformObjectNeedsTransformation(t schemaPkg.Transformer[bool], schema *core.Schema, value bool) (bool, error) {
	return schemaPkg.TransformObjectProperties(schema, value, func(propName string, propSchema *core.Schema, value bool) (bool, error) {
		value, err := schemaPkg.Transform(t, propSchema, value)
		if err != nil {
			return value, err
		}
		if value {
			return true, schemaPkg.TransformStop
		}
		return false, nil
	})
}

func transformConstraintsNeedsTransformation(t schemaPkg.Transformer[bool], kind schemaPkg.ConstraintKind, schemaRefs schemaPkg.SchemaRefs, value bool) (bool, error) {
	value, err := schemaPkg.TransformSchemasArray(t, schemaRefs, value)
	if err != nil {
		return value, err
	}
	if value {
		return true, schemaPkg.TransformStop
	}
	return false, nil

}

// Recursively transforms the value based on the schema that may contain transformations
// If the schema doesn't contain any transformation, then the value is unchanged
func transformValue(schema *core.Schema, transformationKey string, value any) (any, error) {
	t := &commonSchemaTransformer[any]{
		tKey:                 transformationKey,
		transformSpecs:       doTransformsToValue,
		transformArray:       transformArrayValue,
		transformObject:      transformObjectValue,
		transformConstraints: transformConstraintsValue,
	}
	return schemaPkg.Transform[any](t, schema, value)
}

func transformArrayValue(t schemaPkg.Transformer[any], schema *core.Schema, itemSchema *core.Schema, value any) (any, error) {
	valueSlice, ok := value.([]any)
	if !ok {
		return value, fmt.Errorf("expected []any, got %T %#v", value, value)
	}

	cs := utils.NewCOWSliceFunc(valueSlice, utils.IsSameValueOrPointer)
	for i, itemValue := range valueSlice {
		convertedValue, err := schemaPkg.Transform(t, itemSchema, itemValue)
		if err != nil {
			return value, err
		}
		cs.Set(i, convertedValue)
	}

	valueSlice, _ = cs.Release()
	return valueSlice, nil
}

func transformObjectValue(t schemaPkg.Transformer[any], schema *core.Schema, value any) (any, error) {
	valueMap, ok := value.(map[string]any)
	if !ok {
		return value, fmt.Errorf("expected map[string]any, got %T %#v", value, value)
	}
	cm, err := schemaPkg.TransformObjectProperties(
		schema,
		utils.NewCOWMapFunc(valueMap, utils.IsSameValueOrPointer),
		func(propName string, propSchema *core.Schema, cm *utils.COWMap[string, any],
		) (*utils.COWMap[string, any], error) {
			propValue, ok := valueMap[propName]
			if !ok {
				return cm, nil
			}

			convertedFieldValue, err := schemaPkg.Transform(t, propSchema, propValue)
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

func transformConstraintsValue(t schemaPkg.Transformer[any], kind schemaPkg.ConstraintKind, schemaRefs schemaPkg.SchemaRefs, value any) (any, error) {
	// TODO: handle kind properly, see https://swagger.io/docs/specification/data-models/oneof-anyof-allof-not/
	return schemaPkg.TransformSchemasArray(t, schemaRefs, value)
}
