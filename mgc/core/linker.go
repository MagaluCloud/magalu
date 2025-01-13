package core

import (
	"slices"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type Linker interface {
	Name() string
	Description() string
	IsInternal() bool
	// Describes the additional parameters required by the created executor.
	//
	// This will match CreateExecutor().ParametersSchema()
	AdditionalParametersSchema() *Schema
	// Describes the additional configuration required by the created executor.
	//
	// This will match CreateExecutor().ConfigsSchema()
	AdditionalConfigsSchema() *Schema
	ResultSchema() *Schema
	// Create an executor based on a result.
	//
	// The returned executor will have ParametersSchema() matching AdditionalParametersSchema()
	// and ConfigsSchema() matching AdditionalConfigsSchema()
	CreateExecutor(originalResult Result) (exec Executor, err error)
	Links() Links

	IsTargetTerminatorExecutor() bool
}

type Links map[string]Linker

type simpleLink struct {
	name                          string
	description                   string
	owner                         Executor
	target                        Executor
	fromOwner                     map[string]string
	fromResult                    map[string]string
	getAdditionalParametersSchema func() *Schema
	getAdditionalConfigsSchema    func() *Schema
}

var _ Linker = (*simpleLink)(nil)

type SimpleLinkSpec struct {
	Owner  Executor
	Target Executor
	// Maps parameters from Owner into Target
	FromOwner map[string]string
	// Maps parameters from Owner result into Target
	FromResult map[string]string
}

func NewSimpleLink(s SimpleLinkSpec) *simpleLink {
	l := &simpleLink{
		name:        s.Target.Name(),
		description: s.Target.Description(),
		owner:       s.Owner,
		target:      s.Target,
		fromOwner:   s.FromOwner,
		fromResult:  s.FromResult,
	}

	l.getAdditionalParametersSchema = utils.NewLazyLoader[*Schema](
		func() *Schema {
			return l.additionalParams()
		},
	)

	l.getAdditionalConfigsSchema = utils.NewLazyLoader[*Schema](
		func() *Schema {
			return additionalProps(l.target.ConfigsSchema(), l.owner.ResultSchema(), l.owner.ConfigsSchema())
		},
	)

	return l
}

func (l *simpleLink) Name() string {
	return l.name
}

func (l *simpleLink) Description() string {
	return l.description
}

func (l *simpleLink) IsInternal() bool {
	return l.target.IsInternal()
}

func additionalProps(target *Schema, sources ...*Schema) *Schema {
	additional := map[string]*Schema{}
	required := []string{}

TargetLoop:
	for propName, propRef := range target.Properties {
		for _, source := range sources {
			if _, ok := source.Properties[propName]; ok {
				continue TargetLoop
			}
		}

		additional[propName] = (*Schema)(propRef.Value)
		if slices.Contains(target.Required, propName) {
			required = append(required, propName)
		}
	}

	return mgcSchemaPkg.NewObjectSchema(additional, required)
}

func (l *simpleLink) additionalParams() *Schema {
	additional := map[string]*Schema{}
	required := []string{}

TargetLoop:
	for propName, propRef := range l.target.ParametersSchema().Properties {
		for _, v := range l.fromResult {
			if propName == v {
				continue TargetLoop
			}
		}

		for _, v := range l.fromOwner {
			if propName == v {
				continue TargetLoop
			}
		}

		additional[propName] = (*Schema)(propRef.Value)
		if slices.Contains(l.target.ParametersSchema().Required, propName) {
			required = append(required, propName)
		}
	}

	return mgcSchemaPkg.NewObjectSchema(additional, required)
}

func (l *simpleLink) AdditionalParametersSchema() *Schema {
	return l.getAdditionalParametersSchema()
}

func (l *simpleLink) AdditionalConfigsSchema() *Schema {
	return l.getAdditionalConfigsSchema()
}

func (l *simpleLink) ResultSchema() *Schema {
	return l.target.ResultSchema()
}

func (l *simpleLink) Links() Links {
	return l.target.Links()
}

func injectValues(dst map[string]any, schema *Schema, sources ...map[string]any) {
SchemaLoop:
	for propName, propSchema := range schema.Properties {
		for _, source := range sources {
			if value, ok := source[propName]; ok {
				if err := propSchema.Value.VisitJSON(value); err != nil {
					// TODO: Should this be an error?
					continue
				}
				dst[propName] = value
				continue SchemaLoop
			}
		}
	}
}

func (l *simpleLink) CreateExecutor(originalResult Result) (Executor, error) {
	preparedParams := Parameters{}
	preparedConfigs := Configs{}

	var resultMap map[string]any
	if result, ok := ResultAs[ResultWithValue](originalResult); ok {
		resultMap, _ = result.Value().(map[string]any)
	}

	for sourceName, linkName := range l.fromResult {
		value, ok := resultMap[sourceName]
		if !ok {
			panic("parameter not found in original result")
		}

		preparedParams[linkName] = value
	}

	for sourceName, linkName := range l.fromOwner {
		value, ok := originalResult.Source().Parameters[sourceName]
		if !ok {
			panic("parameter not found in original result")
		}

		preparedParams[linkName] = value
	}

	injectValues(preparedConfigs, l.target.ConfigsSchema(), resultMap, originalResult.Source().Configs)

	var exec LinkExecutor = NewLinkExecutor(l.target, preparedParams, preparedConfigs, l.AdditionalParametersSchema(), l.AdditionalConfigsSchema())

	if _, ok := ExecutorAs[ConfirmableExecutor](l.target); ok {
		exec = NewLinkConfirmableExecutor(exec)
	}

	if _, ok := ExecutorAs[PromptInputExecutor](l.target); ok {
		exec = NewLinkPromptInputExecutor(exec)
	}

	if _, ok := ExecutorAs[TerminatorExecutor](l.target); ok {
		exec = NewLinkTerminatorExecutor(exec)
	}

	return exec, nil
}

func (l *simpleLink) IsTargetTerminatorExecutor() bool {
	_, isTerminator := ExecutorAs[TerminatorExecutor](l.target)
	return isTerminator
}
