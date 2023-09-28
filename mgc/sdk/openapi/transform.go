package openapi

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
	"magalu.cloud/core"
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

func doTransform(spec *transformSpec, value any) (any, error) {
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

func doTransforms(specs []*transformSpec, value any) (result any, err error) {
	result = value
	for _, spec := range specs {
		result, err = doTransform(spec, result)
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

	if !needsTransformation(schema, transformationKey) {
		return nil
	}

	return func(value T) (converted T, err error) {
		r, err := transform(schema, transformationKey, value)
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

func needsTransformation(schema *core.Schema, transformationKey string) bool {
	specs := getTransformationSpecs(schema.Extensions, transformationKey)
	if specs != nil {
		return true
	}

	switch schema.Type {
	case "string", "number", "integer", "boolean", "null":
		return false

	case "object":
		for _, ref := range schema.Properties {
			propSchema := (*core.Schema)(ref.Value)
			if propSchema != nil {
				if needsTransformation(propSchema, transformationKey) {
					return true
				}
			}
		}
		return false

	case "array":
		if schema.Items != nil && schema.Items.Value != nil {
			itemSchema := (*core.Schema)(schema.Items.Value)
			if needsTransformation(itemSchema, transformationKey) {
				return true
			}
		}
		return false

	default:
		sub := []openapi3.SchemaRefs{schema.AllOf, schema.AnyOf, schema.OneOf}
		for _, refs := range sub {
			for _, ref := range refs {
				if ref.Value != nil {
					subSchema := (*core.Schema)(ref.Value)
					if needsTransformation(subSchema, transformationKey) {
						return true
					}
				}
			}
		}
		return false
	}
}

func transform(schema *core.Schema, transformationKey string, value any) (any, error) {
	specs := getTransformationSpecs(schema.Extensions, transformationKey)
	if specs != nil {
		return doTransforms(specs, value)
	}

	switch schema.Type {
	case "string", "number", "integer", "boolean", "null":
		return value, nil

	case "object":
		valueMap, ok := value.(map[string]any)
		if !ok {
			return value, nil
		}
		cm := utils.NewCOWMapFunc(valueMap, utils.IsSameValueOrPointer)
		for k, ref := range schema.Properties {
			propSchema := (*core.Schema)(ref.Value)
			if propSchema != nil {
				if propValue, ok := valueMap[k]; ok {
					convertedValue, err := transform(propSchema, transformationKey, propValue)
					if err != nil {
						return value, err
					}
					cm.Set(k, convertedValue)
				}
			}
		}
		valueMap, _ = cm.Release()
		return valueMap, nil

	case "array":
		valueSlice, ok := value.([]any)
		if !ok {
			return value, nil
		}
		cs := utils.NewCOWSliceFunc(valueSlice, utils.IsSameValueOrPointer)
		if schema.Items != nil && schema.Items.Value != nil {
			itemSchema := (*core.Schema)(schema.Items.Value)
			for i, itemValue := range valueSlice {
				convertedValue, err := transform(itemSchema, transformationKey, itemValue)
				if err != nil {
					return value, err
				}
				cs.Set(i, convertedValue)
			}
		}
		valueSlice, _ = cs.Release()
		return valueSlice, nil

	default:
		sub := []openapi3.SchemaRefs{schema.AllOf, schema.AnyOf, schema.OneOf}
		for _, refs := range sub {
			for _, ref := range refs {
				if ref.Value != nil {
					subSchema := (*core.Schema)(ref.Value)
					var err error
					value, err = transform(subSchema, transformationKey, value)
					if err != nil {
						return value, err
					}
				}
			}
		}
		return value, nil
	}
}
