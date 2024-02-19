package core

import (
	"context"

	"maps"
)

type LinkExecutor interface {
	Executor
	ExecutorWrapper
	extendParametersAndConfigs(parameters Parameters, configs Configs) (p Parameters, c Configs)
}

type linkExecutor struct {
	Executor
	preparedParameters         Parameters
	preparedConfigs            Configs
	additionalParametersSchema *Schema
	additionalConfigsSchema    *Schema
}

// Wraps the target executor, handling prepared values before calling it.
//
// The returned Executor will expose additionalParametersSchema() as ParametersSchema()
// and additionalConfigsSchema as ConfigsSchema().
//
// Then the Execute() will automatically copy the preparedParameters and preparedConfigs
// in the given parameters and configs (respectively), before calling target.Execute().
//
// Other Executor methods will be passed thru target without further modifications.
func NewLinkExecutor(
	target Executor,
	preparedParameters Parameters,
	preparedConfigs Configs,
	additionalParametersSchema *Schema,
	additionalConfigsSchema *Schema,
) *linkExecutor {
	return &linkExecutor{
		Executor:                   target,
		preparedParameters:         preparedParameters,
		preparedConfigs:            preparedConfigs,
		additionalParametersSchema: additionalParametersSchema,
		additionalConfigsSchema:    additionalConfigsSchema,
	}
}

func (l *linkExecutor) extendParametersAndConfigs(parameters Parameters, configs Configs) (p Parameters, c Configs) {
	if len(l.preparedParameters) == 0 && len(l.preparedConfigs) == 0 {
		return parameters, configs
	}

	p = maps.Clone(parameters)
	if p == nil {
		p = Parameters{}
	}
	c = maps.Clone(configs)
	if c == nil {
		c = Configs{}
	}
	maps.Copy(p, l.preparedParameters)
	maps.Copy(c, l.preparedConfigs)
	return
}

func (l *linkExecutor) ParametersSchema() *Schema {
	return l.additionalParametersSchema
}

func (l *linkExecutor) ConfigsSchema() *Schema {
	return l.additionalConfigsSchema
}

func (l *linkExecutor) Execute(ctx context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	p, c := l.extendParametersAndConfigs(parameters, configs)
	r, e := l.Executor.Execute(ctx, p, c)
	originalSource := ResultSource{l, ctx, parameters, configs}
	return ExecutorWrapResultSource(originalSource, r, e)
}

func (l *linkExecutor) Unwrap() Executor {
	return l.Executor
}

var _ Executor = (*linkExecutor)(nil)
var _ LinkExecutor = (*linkExecutor)(nil)
var _ ExecutorWrapper = (*linkExecutor)(nil)

type linkTerminatorExecutor struct {
	LinkExecutor
	tExec TerminatorExecutor
}

// Wraps a linkExecutor that implements TerminatorExecutor
//
// linkExecutor.Unwrap() (target) must implement the TerminatorExecutor interface, otherwise it will panic
func NewLinkTerminatorExecutor(linkExecutor LinkExecutor) *linkTerminatorExecutor {
	tExec, ok := ExecutorAs[TerminatorExecutor](linkExecutor)
	if !ok {
		panic("linkExecutor target must implement TerminatorExecutor")
	}
	return &linkTerminatorExecutor{linkExecutor, tExec}
}

func (l *linkTerminatorExecutor) Execute(ctx context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	r, e := l.LinkExecutor.Execute(ctx, parameters, configs)
	return ExecutorWrapResult(l, r, e)
}

func (l *linkTerminatorExecutor) ExecuteUntilTermination(ctx context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	p, c := l.extendParametersAndConfigs(parameters, configs)
	r, e := l.tExec.ExecuteUntilTermination(ctx, p, c)
	originalSource := ResultSource{l, ctx, parameters, configs}
	return ExecutorWrapResultSource(originalSource, r, e)
}

func (l *linkTerminatorExecutor) ExecuteUntilTerminationIgnoreErrors(ctx context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	p, c := l.extendParametersAndConfigs(parameters, configs)
	r, e := l.tExec.ExecuteUntilTerminationIgnoreErrors(ctx, p, c)
	originalSource := ResultSource{l, ctx, parameters, configs}
	return ExecutorWrapResultSource(originalSource, r, e)
}

func (l *linkTerminatorExecutor) Unwrap() Executor {
	return l.LinkExecutor
}

var _ Executor = (*linkTerminatorExecutor)(nil)
var _ LinkExecutor = (*linkTerminatorExecutor)(nil)
var _ TerminatorExecutor = (*linkTerminatorExecutor)(nil)
var _ ExecutorWrapper = (*linkTerminatorExecutor)(nil)

type linkConfirmableExecutor struct {
	LinkExecutor
	cExec ConfirmableExecutor
}

// Wraps a linkExecutor that implements ConfirmableExecutor
//
// linkExecutor.Unwrap() (target) must implement the ConfirmableExecutor interface, otherwise it will panic
func NewLinkConfirmableExecutor(linkExecutor LinkExecutor) *linkConfirmableExecutor {
	tConfirm, ok := ExecutorAs[ConfirmableExecutor](linkExecutor)
	if !ok {
		panic("linkExecutor target must implement ConfirmableExecutor")
	}
	return &linkConfirmableExecutor{linkExecutor, tConfirm}
}

func (l *linkConfirmableExecutor) Execute(ctx context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	r, e := l.LinkExecutor.Execute(ctx, parameters, configs)
	return ExecutorWrapResult(l, r, e)
}

func (l *linkConfirmableExecutor) ConfirmPrompt(parameters Parameters, configs Configs) (message string) {
	p, c := l.extendParametersAndConfigs(parameters, configs)
	return l.cExec.ConfirmPrompt(p, c)
}

func (l *linkConfirmableExecutor) Unwrap() Executor {
	return l.LinkExecutor
}

var _ Executor = (*linkConfirmableExecutor)(nil)
var _ LinkExecutor = (*linkConfirmableExecutor)(nil)
var _ ConfirmableExecutor = (*linkConfirmableExecutor)(nil)
var _ ExecutorWrapper = (*linkConfirmableExecutor)(nil)

type linkPromptInputExecutor struct {
	LinkExecutor
	pExec PromptInputExecutor
}

func NewLinkPromptInputExecutor(linkExecutor LinkExecutor) *linkPromptInputExecutor {
	tConfirm, ok := ExecutorAs[PromptInputExecutor](linkExecutor)
	if !ok {
		panic("linkExecutor target must implement PromptInputExecutor")
	}
	return &linkPromptInputExecutor{linkExecutor, tConfirm}
}

func (l *linkPromptInputExecutor) Execute(ctx context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	r, e := l.LinkExecutor.Execute(ctx, parameters, configs)
	return ExecutorWrapResult(l, r, e)
}

func (l *linkPromptInputExecutor) PromptInput(parameters Parameters, configs Configs) (message string, validate func(input string) error) {
	p, c := l.extendParametersAndConfigs(parameters, configs)
	return l.pExec.PromptInput(p, c)
}

func (l *linkPromptInputExecutor) Unwrap() Executor {
	return l.LinkExecutor
}

var _ Executor = (*linkPromptInputExecutor)(nil)
var _ LinkExecutor = (*linkPromptInputExecutor)(nil)
var _ PromptInputExecutor = (*linkPromptInputExecutor)(nil)
var _ ExecutorWrapper = (*linkPromptInputExecutor)(nil)
