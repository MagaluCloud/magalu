package openapi

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/jsonpointer"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"magalu.cloud/core"
	"magalu.cloud/core/http"
	"magalu.cloud/core/schema"
)

type openapiLinker struct {
	name                 string
	description          string
	owner                *Operation
	link                 *openapi3.Link
	additionalParameters *core.Schema
	additionalConfigs    *core.Schema
}

func insertParameterCb(
	oapiName string,
	dst map[string]core.Value,
	value core.Value,
) cbForEachParameter {
	// Find equivalent parameter and use its external name for insertion
	return func(externalName string, parameter *openapi3.Parameter) (run bool, err error) {
		isCurrent := oapiName == parameter.Name
		if !isCurrent {
			// According to OpenAPI Spec, link parameter keys can be specified with a
			// location prefix to disambiguate between two target parameters with the same
			// name but different locations, so this also needs to be checked.
			// Ref:
			// The parameter name can be qualified using the parameter location [{in}.]{name} for operations that use the same parameter name in different locations (e.g. path.id).
			isCurrent = oapiName == fmt.Sprintf("%s.%s", parameter.In, parameter.Name)
		}

		if isCurrent {
			// Use external name to be in sync with MGC
			dst[externalName] = value
			return false, nil
		}
		return true, nil
	}
}

func insertParameter(
	op *Operation,
	oapiName string,
	value core.Value,
	dstParams core.Parameters,
	dstConfigs core.Configs,
) {
	finished, _ := op.forEachParameter(parametersLocations, insertParameterCb(oapiName, dstParams, value))
	if !finished {
		return
	}
	_, _ = op.forEachParameter(configLocations, insertParameterCb(oapiName, dstConfigs, value))
}

func fillMissingConfigs(preparedConfigs core.Configs, schema *core.Schema, sourceConfigs core.Configs) {
	for configName := range schema.Properties {
		_, isPresent := preparedConfigs[configName]
		if isPresent {
			continue
		}
		val, ok := sourceConfigs[configName]
		if !ok {
			continue
		}
		preparedConfigs[configName] = val
	}
}

func (l *openapiLinker) addParameters(
	operation *Operation,
	specResolver *linkSpecResolver,
	preparedParams core.Parameters,
	preparedConfigs core.Configs,
) error {
	for paramOAPIName, paramSpec := range l.link.Parameters {
		if paramSpec == nil {
			continue
		}

		resolved, found, err := specResolver.resolve(paramSpec)
		if err != nil {
			return err
		}

		if !found {
			continue
		}

		insertParameter(operation, paramOAPIName, resolved, preparedParams, preparedConfigs)
	}

	return nil
}

// TODO: This function only deals with one-level deep JSON Pointers, we should handle arbitrary depths later on
func (l *openapiLinker) addReqBodyParameters(
	operation *Operation,
	specResolver *linkSpecResolver,
	preparedParams core.Parameters,
) error {
	// The official OAPI specification for link request body is, for some reason, different from
	// the parameters and, thus, unusable. The issue can be tracked here: https://github.com/OAI/OpenAPI-Specification/issues/1594
	// Until a version of OAPI fixes this, the extension specified by @anentropic will be used.
	// Ref: https://apigraph.readthedocs.io/en/latest/reference/openapi-extensions.html#x-apigraph-requestbodyparameters
	if reqBodyParamsSpec, ok := getExtensionObject(operation.extensionPrefix, "requestBodyParameters", l.link.Extensions, nil); ok {
		reqBodyParams := map[string]core.Value{}
		for jpStr, rtExpStr := range reqBodyParamsSpec {
			resolved, found, err := specResolver.resolve(rtExpStr)
			if err != nil {
				return err
			}

			if !found {
				continue
			}

			jp, err := jsonpointer.New(jpStr)
			if err != nil {
				return fmt.Errorf("malformed json pointer: '%s'", jpStr)
			}

			// Set to 'reqBodyParams' instead of 'preparedParams' because name in JSON Pointer is internal,
			// without 'x-' extension transformations
			_, err = jp.Set(reqBodyParams, resolved)
			if err != nil {
				return fmt.Errorf("failed to set jsonpointer '%s' on object %#v using value %#v", jpStr, preparedParams, resolved)
			}
		}

		// Translate names to set to 'preparedParams'
		_, _ = operation.forEachParameterName(func(externalName, internalName, location string) (run bool, err error) {
			if value, ok := reqBodyParams[internalName]; ok {
				preparedParams[externalName] = value
			}
			return true, nil
		})
	}

	return nil
}

func opParameterValueResolver(op *Operation, paramData core.Parameters) func(location, name string) (core.Value, bool) {
	return func(location, name string) (core.Value, bool) {
		var result core.Value
		notFound, err := op.forEachParameterWithValue(
			paramData,
			[]string{location},
			func(externalName string, parameter *openapi3.Parameter, value any) (run bool, err error) {
				if name == parameter.Name {
					result = value
					return false, nil
				}
				return true, nil
			},
		)
		if err != nil || notFound {
			return nil, false
		}
		return result, true
	}
}

// START: Linker

func (l *openapiLinker) Name() string {
	if l.Target() == nil {
		return "missing"
	}
	return l.name
}

func (l *openapiLinker) Description() string {
	return l.description
}

func (l *openapiLinker) AdditionalParametersSchema() *core.Schema {
	if l.additionalParameters == nil {
		// TODO: Handle errors in a better, safer way
		target := l.Target()
		if target == nil {
			return nil
		}
		op, ok := core.ExecutorAs[*Operation](target)
		if !ok {
			return nil
		}

		targetParameters := target.ParametersSchema()
		props := map[string]*core.Schema{}
		required := []string{}

		_, _ = op.forEachParameterName(func(externalName, internalName, location string) (run bool, err error) {
			if _, ok := l.link.Parameters[internalName]; ok {
				return true, nil
			}

			// The parameter name can be qualified using the parameter location [{in}.]{name} for
			// operations that use the same parameter name in different locations (e.g. path.id).
			if _, ok := l.link.Parameters[fmt.Sprintf("%s.%s", location, internalName)]; ok {
				return true, nil
			}

			props[externalName] = (*core.Schema)(targetParameters.Properties[externalName].Value)
			if slices.Contains(targetParameters.Required, externalName) {
				required = append(required, externalName)
			}

			return true, nil
		})
		l.additionalParameters = schema.NewObjectSchema(props, required)
	}
	return l.additionalParameters
}

func (l *openapiLinker) AdditionalConfigsSchema() *core.Schema {
	if l.additionalConfigs == nil {
		target := l.Target()
		if target == nil {
			return nil
		}

		targetConfigs := target.ConfigsSchema()
		props := map[string]*core.Schema{}
		required := []string{}

		for targetConfigName, targetConfigRef := range targetConfigs.Properties {
			if _, ok := l.owner.ConfigsSchema().Properties[targetConfigName]; ok {
				continue
			}

			props[targetConfigName] = (*core.Schema)(targetConfigRef.Value)
			if slices.Contains(targetConfigs.Required, targetConfigName) {
				required = append(required, targetConfigName)
			}
		}

		l.additionalConfigs = schema.NewObjectSchema(props, required)
	}
	return l.additionalConfigs
}

func (l *openapiLinker) PrepareLink(
	originalResult core.Result,
	additionalParameters core.Parameters,
	additionalConfigs core.Configs,
) (core.Parameters, core.Configs, error) {
	target := l.Target()
	if target == nil {
		return nil, nil, fmt.Errorf("missing target")
	}
	op, ok := core.ExecutorAs[*Operation](target)
	if !ok {
		return nil, nil, fmt.Errorf("link '%s' has unexpected target type. Expected *Operation", l.Name())
	}

	err := l.AdditionalParametersSchema().VisitJSON(additionalParameters, openapi3.MultiErrors())
	if err != nil {
		return nil, nil, fmt.Errorf("additional parameters passed to PrepareLink are invalid: %w", err)
	}

	err = l.AdditionalConfigsSchema().VisitJSON(additionalConfigs, openapi3.MultiErrors())
	if err != nil {
		return nil, nil, fmt.Errorf("additional configs passed to PrepareLink are invalid: %w", err)
	}

	preparedParams := core.Parameters{}
	preparedConfigs := core.Configs{}

	httpResult, ok := core.ResultAs[http.HttpResult](originalResult)
	if !ok {
		return nil, nil, fmt.Errorf("result passed to PrepareLink has unexpected type. Expected HttpResult for link '%s'", l.Name())
	}

	parameterValueResolver := opParameterValueResolver(op, originalResult.Source().Parameters)
	specResolver := linkSpecResolver{httpResult, parameterValueResolver}

	if err := l.addParameters(op, &specResolver, preparedParams, preparedConfigs); err != nil {
		return nil, nil, err
	}

	if err := l.addReqBodyParameters(op, &specResolver, preparedParams); err != nil {
		return nil, nil, err
	}

	fillMissingConfigs(preparedConfigs, target.ConfigsSchema(), originalResult.Source().Configs)

	maps.Copy(preparedParams, additionalParameters)
	maps.Copy(preparedConfigs, additionalConfigs)

	return preparedParams, preparedConfigs, nil
}

func (l *openapiLinker) Target() core.Executor {
	exec := l.owner.execResolver.get(l.link.OperationID)
	if exec == nil {
		exec, _ = l.owner.execResolver.resolve(l.link.OperationRef)
		if exec == nil {
			return nil
		}
	}
	return exec
}

var _ core.Linker = (*openapiLinker)(nil)

// END: Linker
