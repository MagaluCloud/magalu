package openapi

import (
	"fmt"
	"strings"

	"slices"

	"github.com/getkin/kin-openapi/openapi3"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
)

type parameterWithName struct {
	name      string
	parameter *openapi3.Parameter
}

func collectParameters(byNameAndLocation map[string]map[string]*parameterWithName, parameters openapi3.Parameters, extensionPrefix *string) {
	for _, ref := range parameters {
		// "A unique parameter is defined by a combination of a name and location."
		parameter := ref.Value

		if getHiddenExtension(extensionPrefix, parameter.Extensions) {
			continue
		}

		byLocation, exists := byNameAndLocation[parameter.Name]
		if !exists {
			byLocation = map[string]*parameterWithName{}
			byNameAndLocation[parameter.Name] = byLocation
		}

		name := getNameExtension(extensionPrefix, parameter.Extensions, "")
		byLocation[parameter.In] = &parameterWithName{name, parameter}
	}
}

func finalizeParameters(byNameAndLocation map[string]map[string]*parameterWithName) []*parameterWithName {
	parameters := []*parameterWithName{}

	for name, byLocation := range byNameAndLocation {
		for location, pn := range byLocation {
			if pn.name == "" {
				if len(byLocation) == 1 {
					pn.name = name
				} else {
					pn.name = fmt.Sprintf("%s-%s", location, name)
				}
			}
			parameters = append(parameters, pn)
		}
	}

	return parameters
}

type parameters struct {
	getParameters   func() []*parameterWithName
	getPositionals  func() []string
	getHiddenFlags  func() []string
	extensionPrefix *string
}

func newParameters(key string, parentParameters openapi3.Parameters, opParameters openapi3.Parameters, extensionPrefix *string) *parameters {
	p := &parameters{
		getParameters: utils.NewLazyLoader(func() []*parameterWithName {
			// operation parameters take precedence over path:
			// https://spec.openapis.org/oas/latest.html#fixed-fields-7
			// "the new definition will override it but can never remove it"
			// "A unique parameter is defined by a combination of a name and location."
			m := map[string]map[string]*parameterWithName{}
			collectParameters(m, parentParameters, extensionPrefix)
			collectParameters(m, opParameters, extensionPrefix)
			return finalizeParameters(m)
		}),
		extensionPrefix: extensionPrefix,
	}

	p.getPositionals = utils.NewLazyLoader(func() []string {
		type paramIndex struct {
			pos          int
			externalName string
		}
		var params []paramIndex
		_, _ = p.forEach([]string{openapi3.ParameterInPath}, func(externalName string, parameter *openapi3.Parameter) (run bool, err error) {
			pos := strings.Index(key, "{"+parameter.Name+"}")
			params = append(params, paramIndex{pos, externalName})
			return true, nil
		})
		if len(params) == 0 {
			return nil
		}

		slices.SortFunc(params, func(a, b paramIndex) int {
			if a.pos < b.pos {
				return -1
			} else if a.pos > b.pos {
				return 1
			}
			return 0
		})
		positionals := make([]string, len(params))
		for i, p := range params {
			positionals[i] = p.externalName
		}
		return positionals
	})

	return p
}

type cbForEachParameter func(externalName string, parameter *openapi3.Parameter) (run bool, err error)

func (p *parameters) forEach(locations []string, cb cbForEachParameter) (finished bool, err error) {
	for _, pn := range p.getParameters() {
		name := pn.name
		parameter := pn.parameter

		if parameter.Schema == nil || parameter.Schema.Value == nil {
			continue
		}

		if !slices.Contains(locations, parameter.In) {
			continue
		}
		if parameter.In == openapi3.ParameterInHeader && strings.HasPrefix(strings.ToLower(parameter.Name), "content-") {
			continue
		}

		run, err := cb(name, parameter)
		if err != nil {
			return false, err
		}
		if !run {
			return false, nil
		}
	}

	return true, nil
}

type cbForEachParameterWithValue func(externalName string, parameter *openapi3.Parameter, value any) (run bool, err error)

func (p *parameters) forEachWithValue(values map[string]any, locations []string, cb cbForEachParameterWithValue) (finished bool, err error) {
	return p.forEach(locations, func(externalName string, parameter *openapi3.Parameter) (run bool, err error) {
		value, ok := values[externalName]
		if !ok {
			value = parameter.Schema.Value.Default
			if value == nil {
				return true, nil
			}
		}
		return cb(externalName, parameter, value)
	})
}

func (p *parameters) addToSchema(schema *mgcSchemaPkg.Schema, locations []string) error {
	_, err := p.forEach(locations, func(externalName string, parameter *openapi3.Parameter) (run bool, err error) {
		paramSchemaRef := mgcSchemaPkg.NewCOWSchemaRef(parameter.Schema)
		paramSchema := paramSchemaRef.ValueCOW()

		desc := getDescriptionExtension(p.extensionPrefix, parameter.Extensions, parameter.Description)
		if desc == "" {
			desc = getDescriptionExtension(p.extensionPrefix, paramSchema.Extensions(), paramSchema.Description())
		}

		if desc != "" {
			paramSchema.SetDescription(desc)
		}

		if existing := schema.Properties[externalName]; existing != nil {
			return false, &utils.ChainedError{Name: externalName, Err: fmt.Errorf("already exists as schema %v", existing)}
		}

		schema.Properties[externalName] = paramSchemaRef.Peek()

		if parameter.Required && !slices.Contains(schema.Required, externalName) {
			schema.Required = append(schema.Required, externalName)
		}

		return true, nil
	})

	slices.Sort(schema.Required)

	return err
}
