package core

import (
	"context"
	"errors"
	"fmt"
)

type LinksSpecFn func(e Executor) Links
type RelatedSpecFn func() map[string]Executor
type ExecutorSpecFn func(executor Executor, context context.Context, parameters Parameters, configs Configs) (result Result, err error)

type ExecutorSpec struct {
	DescriptorSpec
	ParametersSchema *Schema
	ConfigsSchema    *Schema
	ResultSchema     *Schema
	Links            LinksSpecFn
	Related          RelatedSpecFn
	PositionalArgs   []string
	HiddenFlags      []string
	Execute          ExecutorSpecFn
}

var errNil = errors.New("cannot be nil")

func (s *ExecutorSpec) Validate() (err error) {
	err = s.DescriptorSpec.Validate()
	if err != nil {
		return
	}

	if s.ParametersSchema == nil {
		return &ChainedError{Name: "ParametersSchema", Err: errNil}
	} else if s.ParametersSchema.Type != nil && !s.ParametersSchema.Type.Includes("object") {
		return &ChainedError{Name: "ParametersSchema", Err: fmt.Errorf("want object, got %q", s.ParametersSchema.Type)}
	}

	if s.ConfigsSchema == nil {
		return &ChainedError{Name: "ConfigsSchema", Err: errNil}
	} else if s.ConfigsSchema.Type != nil && !s.ConfigsSchema.Type.Includes("object") {
		return &ChainedError{Name: "ConfigsSchema", Err: fmt.Errorf("want object, got %q", s.ConfigsSchema.Type)}
	}

	if s.ResultSchema == nil {
		return &ChainedError{Name: "ResultSchema", Err: errNil}
	}

	if s.Execute == nil {
		return &ChainedError{Name: "Execute", Err: errNil}
	}

	return nil
}

type SimpleExecutor struct {
	SimpleDescriptor

	parametersSchema *Schema
	configsSchema    *Schema
	resultSchema     *Schema
	links            LinksSpecFn
	related          RelatedSpecFn
	positionalArgs   []string
	hiddenFlags      []string
	execute          ExecutorSpecFn
}

func NewSimpleExecutor(spec ExecutorSpec) *SimpleExecutor {
	err := spec.Validate()
	if err != nil {
		logger().Fatalw("invalid spec", "spec", spec, "err", err)
	}
	return &SimpleExecutor{
		SimpleDescriptor{spec.DescriptorSpec},
		spec.ParametersSchema,
		spec.ConfigsSchema,
		spec.ResultSchema,
		spec.Links,
		spec.Related,
		spec.PositionalArgs,
		spec.HiddenFlags,
		spec.Execute,
	}
}

func (e *SimpleExecutor) ParametersSchema() *Schema {
	return e.parametersSchema
}

func (e *SimpleExecutor) ConfigsSchema() *Schema {
	return e.configsSchema
}

func (e *SimpleExecutor) ResultSchema() *Schema {
	return e.resultSchema
}

func (e *SimpleExecutor) Links() Links {
	if e.links == nil {
		return nil
	}
	return e.links(e)
}

func (e *SimpleExecutor) Related() map[string]Executor {
	if e.related == nil {
		return nil
	}
	return e.related()
}

func (e *SimpleExecutor) PositionalArgs() []string {
	return e.positionalArgs
}

func (e *SimpleExecutor) HiddenFlags() []string {
	if e.hiddenFlags == nil {
		return []string{}
	}
	return e.hiddenFlags
}

func (*SimpleExecutor) EmptyResult() Result {
	return NewSimpleResult(ResultSource{}, nil, nil)
}

func (e *SimpleExecutor) Execute(context context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	return e.execute(e, context, parameters, configs)
}

var _ Executor = (*SimpleExecutor)(nil)
