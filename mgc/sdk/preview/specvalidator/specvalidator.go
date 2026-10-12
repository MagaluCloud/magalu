package specvalidator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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
		Paths      map[string]map[string]any `yaml:"paths"`
		Security   any                       `yaml:"security"`
		Components struct {
			SecuritySchemes map[string]any `yaml:"securitySchemes"`
		} `yaml:"components"`
	}
	if err := yaml.Unmarshal(pkg.Spec, &spec); err != nil {
		return fmt.Errorf("spec does not parse: %w", err)
	}
	if !hasOperations(spec.Paths) {
		return errors.New("spec has no operations")
	}
	if scope, ok := internalScope(spec.Paths, spec.Security, spec.Components.SecuritySchemes); ok {
		return fmt.Errorf("spec uses the internal scope %q: internal APIs are not published as preview commands", scope)
	}
	return nil
}

func internalScope(paths map[string]map[string]any, security any, schemes map[string]any) (string, bool) {
	scopes := requirementScopes(security)
	for _, item := range paths {
		for _, method := range operationMethods {
			if op, ok := item[method].(map[string]any); ok {
				scopes = append(scopes, requirementScopes(op["security"])...)
			}
		}
	}
	for _, scheme := range schemes {
		scopes = append(scopes, schemeScopes(scheme)...)
	}

	for _, scope := range scopes {
		if strings.Contains(strings.ToLower(scope), "internal") {
			return scope, true
		}
	}
	return "", false
}

func requirementScopes(security any) []string {
	var scopes []string
	requirements, _ := security.([]any)
	for _, requirement := range requirements {
		schemes, _ := requirement.(map[string]any)
		for _, values := range schemes {
			list, _ := values.([]any)
			for _, value := range list {
				if scope, ok := value.(string); ok {
					scopes = append(scopes, scope)
				}
			}
		}
	}
	return scopes
}

func schemeScopes(scheme any) []string {
	var scopes []string
	fields, _ := scheme.(map[string]any)
	flows, _ := fields["flows"].(map[string]any)
	for _, flow := range flows {
		flowFields, _ := flow.(map[string]any)
		declared, _ := flowFields["scopes"].(map[string]any)
		for scope := range declared {
			scopes = append(scopes, scope)
		}
	}
	return scopes
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
