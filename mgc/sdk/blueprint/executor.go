package blueprint

import (
	"context"

	"go.uber.org/zap"
	"magalu.cloud/core"
	schemaPkg "magalu.cloud/core/schema"
)

type executor struct {
	core.SimpleDescriptor
	spec         *executorSpec
	logger       *zap.SugaredLogger
	refResolver  *core.BoundRefPathResolver
	resolveError error
}

func newExecutor(spec *childSpec, logger *zap.SugaredLogger, refResolver *core.BoundRefPathResolver) (exec core.Executor, err error) {
	logger = logger.Named(spec.Name)
	execSpec := &spec.executorSpec
	exec = &executor{
		SimpleDescriptor: core.SimpleDescriptor{Spec: spec.DescriptorSpec},
		spec:             execSpec,
		logger:           logger,
		refResolver:      refResolver,
	}

	if execSpec.Confirm != "" {
		exec = core.NewConfirmableExecutor(exec, core.ConfirmPromptWithTemplate(execSpec.Confirm))
	}

	if execSpec.PromptInput != nil {
		exec = core.NewPromptInputExecutor(exec, core.NewPromptInput(execSpec.PromptInput.MessageTemplate, execSpec.PromptInput.ConfirmValueTemplate))
	}

	if execSpec.WaitTermination != nil {
		exec, err = execSpec.WaitTermination.Build(exec, func(result core.ResultWithValue) any {
			if result, ok := core.ResultAs[*executorResult](result); ok {
				return result.jsonPathDocumentWithResult()
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return exec, nil
}

func (e *executor) resolve() error {
	if e.refResolver != nil {
		e.resolveError = e.spec.resolve(e.refResolver, e)
		if e.resolveError != nil {
			e.logger.Errorw(
				"failed to resolve blueprint references",
				"name", e.Name(),
				"description", e.Description(),
				"error", e.resolveError,
			)
		}
		e.refResolver = nil
		if e.spec.parametersSchema == nil {
			e.spec.parametersSchema = &schemaPkg.Schema{}

		}
		if e.spec.configsSchema == nil {
			e.spec.configsSchema = &schemaPkg.Schema{}
		}
		if e.spec.resultSchema == nil {
			e.spec.resultSchema = &schemaPkg.Schema{}
		}
	}
	return e.resolveError
}

func (e *executor) ParametersSchema() *core.Schema {
	_ = e.resolve()
	return e.spec.parametersSchema
}

func (e *executor) ConfigsSchema() *core.Schema {
	_ = e.resolve()
	return e.spec.configsSchema
}

func (e *executor) ResultSchema() *core.Schema {
	_ = e.resolve()
	return e.spec.resultSchema
}

func (e *executor) PositionalArgs() []string {
	return e.spec.PositionalArgs
}
func (e *executor) HiddenFlags() []string {
	if e.spec.HiddenFlags == nil {
		return []string{}
	}
	return e.spec.HiddenFlags
}
func (e *executor) executeStep(
	step *executeStep,
	result *executorResult,
) (err error) {
	jsonPathDocument := result.jsonPathDocument()
	logger := e.logger.With(
		"step", step.Id,
		"target", step.Target,
		"jsonPathDocument", jsonPathDocument,
	)

	if shouldExecute, err := step.shouldExecute(jsonPathDocument); err != nil {
		logger.Warnw("failed to evaluate 'if' condition", "if", step.IfCondition)
		return err
	} else if !shouldExecute {
		logger.Debugw("skipping execution", "if", step.IfCondition)
		result.skip(step)
		return nil
	}

	p, err := step.prepareParameters(jsonPathDocument, e.ParametersSchema())
	if err != nil {
		logger.Warnw(
			"failed to get step parameters",
			"parameters", step.Parameters,
			"stepSchema", step.executor.ParametersSchema(),
			"outerSchema", e.ParametersSchema(),
		)
		return err
	}

	c, err := step.prepareConfigs(jsonPathDocument, e.ConfigsSchema())
	if err != nil {
		logger.Warnw(
			"failed to get step configs",
			"configs", step.Configs,
			"stepSchema", step.executor.ConfigsSchema(),
			"outerSchema", e.ConfigsSchema(),
		)
		return err
	}

	logger = logger.With("parameters", p, "configs", c)
	if step.RetryUntil != nil {
		logger = logger.With("retryUntil", step.RetryUntil)
	}

	var cb core.RetryUntilCb
	ctx := result.Context
	if tExec, ok := core.ExecutorAs[core.TerminatorExecutor](step.executor); ok && step.WaitTermination {
		cb = func() (result core.Result, err error) {
			logger.Debugw("execute step (waitTermination)")
			return tExec.ExecuteUntilTermination(ctx, p, c)
		}
	} else {
		cb = func() (result core.Result, err error) {
			logger.Debugw("execute step")
			return step.executor.Execute(ctx, p, c)
		}
	}

	// retryUntil.run() is a safe nil pointer receiver, will execute only once without checks in that case
	execResult, execErr := step.RetryUntil.run(ctx, cb, func(value core.Value) map[string]any {
		return result.jsonPathDocumentWithCurrent(step, p, c, value)
	})

	if execErr != nil {
		logger.Debugw("failed step", "error", execErr)
		result.reportError(step, p, c, execErr)
	} else {
		logger.Debugw("finished step", "result", execResult)

		var v any
		if vResult, ok := core.ResultAs[core.ResultWithValue](execResult); ok {
			v = vResult.Value()
		}

		execErr = step.check(result.jsonPathDocumentWithCurrent(step, p, c, v))
		if execErr != nil {
			logger.Debugw("failed step check", "error", execErr)
			result.reportError(step, p, c, execErr)
		} else {
			result.reportResult(step, p, c, execResult)
		}
	}
	return nil // regardless of execErr as we want to register errors and let further steps to handle them
}

func (e *executor) Execute(
	ctx context.Context,
	parameters core.Parameters,
	configs core.Configs,
) (r core.Result, err error) {
	err = e.resolve()
	if err != nil {
		return
	}

	result := &executorResult{
		ResultSource: core.ResultSource{
			Executor:   e,
			Context:    ctx,
			Parameters: parameters,
			Configs:    configs,
		},
		Steps:          nil,
		Logger:         e.logger,
		ResultJsonPath: e.spec.Result,
	}

	for _, step := range e.spec.Steps {
		err = e.executeStep(step, result)
		if err != nil {
			return nil, err
		}
	}

	r, err = result.finalize()
	if e.spec.OutputFlag != "" {
		if resultWithValue, ok := core.ResultAs[core.ResultWithValue](result); ok {
			r = core.NewResultWithOriginalSource(r.Source(), core.NewResultWithDefaultOutputOptions(resultWithValue, e.spec.OutputFlag))
		}
	}

	return
}

// This map should not be altered externally
func (e *executor) Links() core.Links {
	return e.spec.linkers
}

// This map should not be altered externally
func (e *executor) Related() map[string]core.Executor {
	return e.spec.relatedExecutors
}

func (e *executor) EmptyResult() core.Result {
	return &executorResult{}
}

var _ core.Executor = (*executor)(nil)
