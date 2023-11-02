package transform

import (
	"fmt"
	"strings"

	"github.com/stoewer/go-strcase"
	mgcSchemaPkg "magalu.cloud/core/schema"
)

type transformerString struct {
	transform func(string) string
}

var _ transformer = (*transformerString)(nil)

func (t *transformerString) TransformValue(value any) (transformedValue any, err error) {
	s, ok := value.(string)
	if !ok {
		err = fmt.Errorf("expected string, got %T", value)
		return
	}

	return t.transform(s), nil
}

func (t *transformerString) TransformSchema(value *mgcSchemaPkg.COWSchema) (*mgcSchemaPkg.COWSchema, error) {
	return value, nil
}

func init() {
	transformers := []struct {
		f     func(string) string
		names []string
	}{
		{strings.ToUpper, []string{"uppercase", "upper-case", "upper"}},
		{strings.ToLower, []string{"lowercase", "lower-case", "lower"}},
		{strcase.KebabCase, []string{"kebabcase", "kebab-case", "kebab"}},
		{strcase.KebabCase, []string{"snakecase", "snake-case", "snake"}},
		{strcase.UpperCamelCase, []string{"pascalcase", "pascal-case", "pascal", "upper-camel"}},
		{strcase.UpperCamelCase, []string{"camelcase", "camel-case", "camel", "lower-camel"}},
	}
	for _, t := range transformers {
		addTransformer[struct{}](func(spec *transformSpec) (transformer, error) {
			return &transformerString{t.f}, nil
		}, t.names...)
	}
}
