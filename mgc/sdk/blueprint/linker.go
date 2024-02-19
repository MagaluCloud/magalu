package blueprint

import (
	"fmt"

	"slices"

	"github.com/PaesslerAG/gval"
	"magalu.cloud/core"
	schemaPkg "magalu.cloud/core/schema"
)

type linker struct {
	spec                 *linkSpec
	owner                core.Executor
	additionalParameters *schemaPkg.Schema
	additionalConfigs    *schemaPkg.Schema
}

func (l *linker) Name() string {
	return l.spec.name
}

func (l *linker) Description() string {
	return l.spec.Description
}

func (l *linker) IsInternal() bool {
	return l.spec.IsInternal
}

func createAdditionalSchema(origSchema *schemaPkg.Schema, providedValues map[string]gval.Evaluable, providedProperties map[string]*schemaPkg.SchemaRef) *schemaPkg.Schema {
	required := []string{}
	properties := map[string]*schemaPkg.Schema{}
	for k, propRef := range origSchema.Properties {
		if _, isProvided := providedValues[k]; isProvided {
			continue
		}

		propSchema := (*schemaPkg.Schema)(propRef.Value)

		if providedRef, isProvided := providedProperties[k]; isProvided {
			providedSchema := (*schemaPkg.Schema)(providedRef.Value)
			if schemaPkg.CheckSimilarJsonSchemas(propSchema, providedSchema) {
				continue
			}
		}

		properties[k] = propSchema
		if slices.Contains(origSchema.Required, k) {
			required = append(required, k)
		}
	}
	return schemaPkg.NewObjectSchema(properties, required)
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

func (l *linker) AdditionalParametersSchema() *schemaPkg.Schema {
	if l.additionalParameters == nil {
		l.additionalParameters = createAdditionalSchema(l.spec.executor.ParametersSchema(), l.spec.parametersJSONPath, nil)
	}
	return l.additionalParameters
}

func (l *linker) AdditionalConfigsSchema() *schemaPkg.Schema {
	if l.additionalConfigs == nil {
		l.additionalConfigs = createAdditionalSchema(l.spec.executor.ConfigsSchema(), l.spec.configsJSONPath, l.owner.ConfigsSchema().Properties)
	}
	return l.additionalConfigs
}

func (l *linker) ResultSchema() *schemaPkg.Schema {
	return l.spec.executor.ResultSchema()
}

func (l *linker) Links() core.Links {
	return l.spec.executor.Links()
}

func (l *linker) CreateExecutor(originalResult core.Result) (target core.Executor, err error) {
	target = l.spec.executor

	result, ok := core.ResultAs[*executorResult](originalResult)
	if !ok {
		return nil, fmt.Errorf("result passed to CreateExecutor has unexpected type. Expected blueprint.executorResult for link '%s'", l.Name())
	}

	jsonPathDocument := result.jsonPathDocument()

	preparedParams, err := prepareMapFromRules(jsonPathDocument, l.spec.parametersJSONPath)
	if err != nil {
		return nil, fmt.Errorf("could not prepare parameters: %w", err)
	}

	preparedConfigs, err := prepareMapFromRules(jsonPathDocument, l.spec.configsJSONPath)
	if err != nil {
		return nil, fmt.Errorf("could not prepare configs: %w", err)
	}
	fillMissingConfigs(preparedConfigs, target.ConfigsSchema(), originalResult.Source().Configs)

	if l.spec.WaitTermination != nil {
		target, err = l.spec.WaitTermination.Build(target, func(result core.ResultWithValue) any {
			if result, ok := core.ResultAs[*executorResult](result); ok {
				return result.jsonPathDocumentWithResult()
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	var exec core.LinkExecutor = core.NewLinkExecutor(target, preparedParams, preparedConfigs, l.additionalParameters, l.additionalConfigs)
	if _, ok := core.ExecutorAs[core.TerminatorExecutor](target); ok {
		exec = core.NewLinkTerminatorExecutor(exec)
	}
	if _, ok := core.ExecutorAs[core.ConfirmableExecutor](target); ok {
		exec = core.NewLinkConfirmableExecutor(exec)
	}

	return exec, nil
}

func (l *linker) IsTargetTerminatorExecutor() bool {
	if l.spec.WaitTermination != nil {
		return true
	}
	if _, isTerminator := core.ExecutorAs[core.TerminatorExecutor](l.spec.executor); isTerminator {
		return true
	}
	return false
}

var _ core.Linker = (*linker)(nil)
