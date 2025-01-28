package openapi

import (
	"fmt"

	"slices"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/http"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/jsonpointer"
)

type openapiLinker struct {
	name                 string
	description          string
	owner                *operation
	link                 *openapi3.Link
	additionalParameters *core.Schema
	additionalConfigs    *core.Schema
	target               core.Executor
}

type extraParameterExtension struct {
	Name     string
	Required bool
	Schema   *core.Schema
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
	op *operation,
	oapiName string,
	value core.Value,
	dstParams core.Parameters,
	dstConfigs core.Configs,
) {
	finished, _ := op.parameters.forEach(parametersLocations, insertParameterCb(oapiName, dstParams, value))
	if !finished {
		return
	}
	_, _ = op.parameters.forEach(configLocations, insertParameterCb(oapiName, dstConfigs, value))
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
	o *operation,
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

		insertParameter(o, paramOAPIName, resolved, preparedParams, preparedConfigs)
	}

	return nil
}

// TODO: This function only deals with one-level deep JSON Pointers, we should handle arbitrary depths later on
func (l *openapiLinker) addReqBodyParameters(
	o *operation,
	specResolver *linkSpecResolver,
	preparedParams core.Parameters,
) error {
	// The official OAPI specification for link request body is, for some reason, different from
	// the parameters and, thus, unusable. The issue can be tracked here: https://github.com/OAI/OpenAPI-Specification/issues/1594
	// Until a version of OAPI fixes this, the extension specified by @anentropic will be used.
	// Ref: https://apigraph.readthedocs.io/en/latest/reference/openapi-extensions.html#x-apigraph-requestbodyparameters
	if reqBodyParamsSpec, ok := getExtensionObject(o.extensionPrefix, "requestBodyParameters", l.link.Extensions, nil); ok {
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
		_, _ = o.forEachParameterName(func(externalName, internalName, location string) (run bool, err error) {
			if value, ok := reqBodyParams[internalName]; ok {
				preparedParams[externalName] = value
			}
			return true, nil
		})
	}

	return nil
}

func opParameterValueResolver(op *operation, paramData core.Parameters) func(location, name string) (core.Value, bool) {
	return func(location, name string) (core.Value, bool) {
		var result core.Value
		notFound, err := op.parameters.forEachWithValue(
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
	return l.name
}

func (l *openapiLinker) Description() string {
	return l.description
}

func (l *openapiLinker) IsInternal() bool {
	return getHiddenExtension(l.owner.extensionPrefix, l.link.Extensions)
}

// Add all extra parameters defined via 'extra-paramters' extension in the Link object. The extension must be an array of
// objects, and each object must match the 'extraParameterExtension' struct. If a parameter has the same name as a
// standard parameter in the target request, it will NOT be added, since that would overshadow the standard parameter
func (l *openapiLinker) addExtraParametersExtension(dst map[string]*core.Schema, required *[]string) {
	if extraParams, ok := getExtensionArray(l.owner.extensionPrefix, "extra-parameters", l.link.Extensions, nil); ok {
		for _, extraSpec := range extraParams {
			param, err := utils.DecodeNewValue[extraParameterExtension](extraSpec)
			if err != nil {
				l.owner.logger.Warnw(
					"unable to decode extra parameter spec for link",
					"link", l.name,
					"spec data", extraSpec,
				)
				continue
			}

			if _, ok := dst[param.Name]; ok {
				l.owner.logger.Debugw(
					"ignoring extra parameter spec, since it overshadows target parameters",
					"link", l.name,
					"parameter name", param.Name,
				)
				continue
			}

			dst[param.Name] = param.Schema
			if param.Required {
				*required = append(*required, param.Name)
			}

			l.owner.logger.Debugw(
				"added extra parameter to link AdditionalParameters schema",
				"link", l.name,
				"parameter name", param.Name,
			)
		}
	}

}

func (l *openapiLinker) AdditionalParametersSchema() *core.Schema {
	if l.additionalParameters == nil {
		target := l.Target()
		op, ok := core.ExecutorAs[*operation](target)
		if !ok {
			l.additionalParameters = mgcSchemaPkg.NewObjectSchema(nil, nil)
			return l.additionalParameters
		}

		targetParameters := target.ParametersSchema()
		props := map[string]*core.Schema{}
		required := []string{}

		_, _ = op.forEachParameterName(func(externalName, internalName, location string) (run bool, err error) {
			if _, ok := l.link.Parameters[internalName]; ok {
				return true, nil
			}

			if reqBodyParamsSpec, ok := getExtensionObject(op.extensionPrefix, "requestBodyParameters", l.link.Extensions, nil); ok {
				for jpStr := range reqBodyParamsSpec {
					jp, err := jsonpointer.New(jpStr)
					if err != nil {
						continue
					}

					tokens := jp.DecodedTokens()
					if len(tokens) == 0 {
						continue
					}

					if tokens[0] == internalName {
						return true, nil
					}
				}
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

		l.addExtraParametersExtension(props, &required)

		l.additionalParameters = mgcSchemaPkg.NewObjectSchema(props, required)
	}
	return l.additionalParameters
}

func (l *openapiLinker) AdditionalConfigsSchema() *core.Schema {
	if l.additionalConfigs == nil {
		target := l.Target()
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

		l.additionalConfigs = mgcSchemaPkg.NewObjectSchema(props, required)
	}
	return l.additionalConfigs
}

func (l *openapiLinker) CreateExecutor(originalResult core.Result) (core.Executor, error) {
	target := l.Target()
	op, ok := core.ExecutorAs[*operation](target)
	if !ok {
		return nil, fmt.Errorf("link '%s' has unexpected target type. Expected *Operation", l.Name())
	}

	httpResult, ok := core.ResultAs[http.HttpResult](originalResult)
	if !ok {
		return nil, fmt.Errorf("result passed to CreateExecutor has unexpected type. Expected HttpResult for link '%s'", l.Name())
	}

	preparedParams := core.Parameters{}
	preparedConfigs := core.Configs{}

	parameterValueResolver := opParameterValueResolver(op, httpResult.Source().Parameters)
	specResolver := linkSpecResolver{httpResult, parameterValueResolver}

	if err := l.addParameters(op, &specResolver, preparedParams, preparedConfigs); err != nil {
		return nil, err
	}

	if err := l.addReqBodyParameters(op, &specResolver, preparedParams); err != nil {
		return nil, err
	}

	fillMissingConfigs(preparedConfigs, target.ConfigsSchema(), httpResult.Source().Configs)

	if wtExt, ok := getExtensionObject(l.owner.extensionPrefix, "wait-termination", l.link.Extensions, nil); ok && wtExt != nil {
		if tExec, err := wrapInTerminatorExecutorWithOwnerResult(target, wtExt, originalResult); err == nil {
			target = tExec
		}
	}

	var exec core.LinkExecutor = core.NewLinkExecutor(target, preparedParams, preparedConfigs, l.AdditionalParametersSchema(), l.AdditionalConfigsSchema())
	if _, ok := core.ExecutorAs[core.TerminatorExecutor](target); ok {
		exec = core.NewLinkTerminatorExecutor(exec)
	}
	if _, ok := core.ExecutorAs[core.ConfirmableExecutor](target); ok {
		exec = core.NewLinkConfirmableExecutor(exec)
	}
	if _, ok := core.ExecutorAs[core.PromptInputExecutor](target); ok {
		exec = core.NewLinkPromptInputExecutor(exec)
	}

	return exec, nil
}

func (l *openapiLinker) ResultSchema() *core.Schema {
	return l.target.ResultSchema()
}

func (l *openapiLinker) Target() core.Executor {
	return l.target
}

func (l *openapiLinker) Links() core.Links {
	return l.target.Links()
}

func (l *openapiLinker) IsTargetTerminatorExecutor() bool {
	if wExt, ok := getExtensionObject(l.owner.extensionPrefix, "wait-termination", l.link.Extensions, nil); ok && wExt != nil {
		return true
	}
	if _, isTerminator := core.ExecutorAs[core.TerminatorExecutor](l.target); isTerminator {
		return true
	}
	return false
}

var _ core.Linker = (*openapiLinker)(nil)

// END: Linker
