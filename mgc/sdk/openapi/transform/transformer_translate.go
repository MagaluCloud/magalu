package transform

import (
	"fmt"
	"reflect"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type transformTranslateSpecItem struct {
	From any `json:"from"`
	To   any `json:"to"`
}

type transformTranslateSpec struct {
	Translations []transformTranslateSpecItem `json:"translations"`
	AllowMissing bool                         `json:"allowMissing,omitempty"`
}

func newTransformerTranslate(tSpec *transformSpec) (transformer, error) {
	t := &transformerTranslate{}
	var err error
	t.spec, err = utils.DecodeNewValue[transformTranslateSpec](tSpec.Parameters)
	if err != nil {
		return nil, fmt.Errorf("invalid translation parameters: %w", err)
	}
	if len(t.spec.Translations) == 0 {
		return nil, fmt.Errorf("invalid translation parameters: missing translations")
	}

	return t, nil
}

type transformerTranslate struct {
	spec *transformTranslateSpec
}

var _ transformer = (*transformerTranslate)(nil)

func (t *transformerTranslate) TransformValue(value any) (any, error) {
	for _, item := range t.spec.Translations {
		if reflect.DeepEqual(item.From, value) {
			return item.To, nil
		}
	}
	if t.spec.AllowMissing {
		return value, nil
	}
	return value, fmt.Errorf("translation not found: %+v", value)
}

func (t *transformerTranslate) TransformSchema(schema *mgcSchemaPkg.COWSchema) (result *mgcSchemaPkg.COWSchema, err error) {
	if schema.Default() == nil && len(schema.Enum()) == 0 {
		return schema, nil
	}

	result = schema

	if schema.Default() != nil {
		var schemaDefault any
		schemaDefault, err = reverseTranslate(t.spec, schema.Default())
		if err != nil {
			return
		}
		schema.SetDefault(schemaDefault)
	}

	enumCow := schema.EnumCOW()
	enumCow.ForEach(func(i int, value any) (run bool) {
		var translatedEnum any
		translatedEnum, err = reverseTranslate(t.spec, value)
		if err != nil {
			return false
		}
		enumCow.Set(i, translatedEnum)
		return true
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

func init() {
	addTransformer[transformTranslateSpec](newTransformerTranslate, "translate")
}
