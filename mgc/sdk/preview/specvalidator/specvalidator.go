// Package specvalidator checks that a preview package holds a usable spec,
// so a bad file fails with a clear message instead of an empty group.
package specvalidator

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"gopkg.in/yaml.v3"
)

var validName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

var operationMethods = []string{"get", "put", "post", "delete", "options", "head", "patch", "trace"}

type Validator struct{}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(level preview.Level, pkg preview.Package) error {
	if !validName.MatchString(pkg.Name) {
		return fmt.Errorf("%q is not a valid %s module name: use letters, digits, '-' and '_'", pkg.Name, level.Name)
	}
	if len(pkg.Spec) == 0 {
		return errors.New("spec is empty")
	}

	var spec struct {
		Paths map[string]map[string]any `yaml:"paths"`
	}
	if err := yaml.Unmarshal(pkg.Spec, &spec); err != nil {
		return fmt.Errorf("spec does not parse: %w", err)
	}
	if !hasOperations(spec.Paths) {
		return errors.New("spec has no operations")
	}
	return nil
}

func hasOperations(paths map[string]map[string]any) bool {
	for _, item := range paths {
		for _, method := range operationMethods {
			if _, ok := item[method]; ok {
				return true
			}
		}
	}
	return false
}

var _ preview.Validator = (*Validator)(nil)
