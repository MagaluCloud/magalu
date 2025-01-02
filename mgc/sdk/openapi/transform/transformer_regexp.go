package transform

import (
	"fmt"
	"regexp"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type transformRegExpSpec struct {
	// Regular Expression as per https://pkg.go.dev/regexp#Compile
	Pattern string `json:"pattern"`
	// Replacement Template as per https://pkg.go.dev/regexp#Regexp.Expand
	Replacement string `json:"replacement"`
}

func newTransformerRegExp(tSpec *transformSpec) (transformer, error) {
	t := &transformerRegExp{}
	var err error
	t.spec, err = utils.DecodeNewValue[transformRegExpSpec](tSpec.Parameters)
	if err != nil {
		return nil, fmt.Errorf("invalid regexp parameters: %w", err)
	}
	if len(t.spec.Pattern) == 0 {
		return nil, fmt.Errorf("invalid regexp parameters: missing pattern")
	}
	t.re, err = regexp.Compile(t.spec.Pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regexp pattern %q: %w", t.spec.Pattern, err)
	}

	return t, nil
}

type transformerRegExp struct {
	spec *transformRegExpSpec
	re   *regexp.Regexp
}

var _ transformer = (*transformerRegExp)(nil)

func (t *transformerRegExp) TransformValue(value any) (transformedValue any, err error) {
	s, ok := value.(string)
	if !ok {
		err = fmt.Errorf("expected string, got %T", value)
		return
	}

	b := []byte{}
	for _, submatches := range t.re.FindAllStringSubmatchIndex(s, -1) {
		b = t.re.ExpandString(b, t.spec.Replacement, s, submatches)
	}
	return string(b), nil
}

func (t *transformerRegExp) TransformSchema(value *mgcSchemaPkg.COWSchema) (*mgcSchemaPkg.COWSchema, error) {
	return value, nil
}

func init() {
	addTransformer[transformRegExpSpec](newTransformerRegExp, "regexp", "regexp-replace")
}
