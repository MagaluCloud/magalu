package transform

import (
	"fmt"
	"strings"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/stoewer/go-strcase"
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
		{strcase.SnakeCase, []string{"snakecase", "snake-case", "snake"}},
		{strcase.UpperCamelCase, []string{"pascalcase", "pascal-case", "pascal", "upper-camel"}},
		{strcase.LowerCamelCase, []string{"camelcase", "camel-case", "camel", "lower-camel"}},
	}
	for _, t := range transformers {
		transformerDef := t
		addTransformer[struct{}](func(spec *transformSpec) (transformer, error) {
			return &transformerString{transformerDef.f}, nil
		}, transformerDef.names...)
	}
}
