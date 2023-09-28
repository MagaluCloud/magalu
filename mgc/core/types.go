package core

import (
	"context"
	"fmt"
	"time"

	"magalu.cloud/core/schema"
)

type Schema = schema.Schema

// NOTE: This is so 'jsonschema' doesn't generate a schema with type string and format
// 'date-time'. We want the raw object schema for later validation
type Time time.Time

// General interface that describes both Executor and Grouper
type Descriptor interface {
	Name() string
	Version() string
	Description() string
}

type DescriptorVisitor func(child Descriptor) (run bool, err error)

type Grouper interface {
	Descriptor
	VisitChildren(visitor DescriptorVisitor) (finished bool, err error)
	GetChildByName(name string) (child Descriptor, err error)
}

// contextKey is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type contextKey string

// grouperContextKey is the key for sdk.Grouper values in Contexts. It is
// unexported; clients use NewGrouperContext() and GrouperFromContext()
// instead of using this key directly.
var grouperContextKey contextKey = "magalu.cloud/core/Grouper"

func NewGrouperContext(parent context.Context, group Grouper) context.Context {
	return context.WithValue(parent, grouperContextKey, group)
}

func GrouperFromContext(ctx context.Context) Grouper {
	if value, ok := ctx.Value(grouperContextKey).(Grouper); !ok {
		return nil
	} else {
		return value
	}
}

// Type comes from the Schema
type Value = any

// Type comes from the Schema
type Example = Value

// TODO: Evaluate if the typealias/type assign is needed at all. If a type assign is needed for some reason,
// the kin-openapi lib will need to be patched to accept type assigns of the correct structure on VisitJSON
// (likely through reflection). As it is now, validation fails with type assigns
type Parameters = map[string]Value
type Configs = Parameters

type Linker interface {
	Name() string
	Description() string
	AdditionalParametersSchema() *Schema
	AdditionalConfigsSchema() *Schema
	PrepareLink(originalResult Result, additionalParameters Parameters, additionalConfigs Configs) (preparedParams Parameters, preparedConfigs Configs, err error)
	Target() Executor
}

type Executor interface {
	Descriptor
	ParametersSchema() *Schema
	ConfigsSchema() *Schema
	// The general schema this executor can produce. It may be oneOf/anyOf with multiple schemas.
	// The Result.Schema() may be a subset of the schema, if multiple were available.
	ResultSchema() *Schema
	// This map should not be altered externally
	Links() map[string]Linker
	// The maps for the parameters and configs should NOT be modified inside the implementation of 'Execute'
	Execute(context context.Context, parameters Parameters, configs Configs) (result Result, err error)
}

// NOTE: whenever you wrap an executor remember to also wrap the result with
// ExecutorWrapResult() so the outmost executor is given as source
type ExecutorWrapper interface {
	Executor
	Unwrap() Executor
}

func ExecutorWrapResult(executorWrapper ExecutorWrapper, result Result, err error) (Result, error) {
	if result != nil {
		result = NewResultWithOriginalExecutor(executorWrapper, result)
	}
	return result, err
}

func ExecutorAs[T Executor](exec Executor) (T, bool) {
	var zeroT T

	for {
		if t, ok := exec.(T); ok {
			return t, true
		}

		if u, ok := exec.(ExecutorWrapper); ok {
			exec = u.Unwrap()
		} else {
			break
		}
	}

	return zeroT, false
}

func VisitAllExecutors(child Descriptor, path []string, visitExecutor func(executor Executor, path []string) (bool, error)) (bool, error) {
	if executor, ok := child.(Executor); ok {
		return visitExecutor(executor, path)
	} else if group, ok := child.(Grouper); ok {
		return group.VisitChildren(func(child Descriptor) (run bool, err error) {
			size := len(path)
			path = append(path, child.Name())
			run, err = VisitAllExecutors(child, path, visitExecutor)
			path = path[:size]

			return run, err
		})
	} else {
		return false, fmt.Errorf("child %v not group/executor", child)
	}
}

type executeResultWrapper struct {
	Executor
	wrapResult func(wrapperExecutor ExecutorWrapper, originalResult Result) (wrappedResult Result, err error)
}

func (o *executeResultWrapper) Execute(ctx context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	result, err = o.Executor.Execute(ctx, parameters, configs)
	if err != nil {
		return
	}
	result, err = o.wrapResult(o, result)
	return ExecutorWrapResult(o, result, err)
}

func (o *executeResultWrapper) Unwrap() Executor {
	return o.Executor
}

var _ Executor = (*executeResultWrapper)(nil)
var _ ExecutorWrapper = (*executeResultWrapper)(nil)

// Wraps (embeds) an executor and wrap its result.
// This may be used to add extra interfaces to a result, such as formatting or output options
func NewExecuteResultWrapper(
	executor Executor,
	wrapResult func(wrapperExecutor ExecutorWrapper, originalResult Result) (wrappedResult Result, err error),
) Executor {
	return &executeResultWrapper{executor, wrapResult}
}

// Wraps (embeds) an executor and add specific result formatting.
func NewExecuteFormat(
	executor Executor,
	getFormatter func(exec Executor, result Result) string,
) Executor {
	return NewExecuteResultWrapper(executor, func(wrapperExecutor ExecutorWrapper, originalResult Result) (wrappedResult Result, err error) {
		result, ok := ResultAs[ResultWithValue](originalResult)
		if !ok {
			return nil, fmt.Errorf("result is not core.ResultWithValue: %T %+v", originalResult, originalResult)
		}
		return NewResultWithDefaultFormatter(result, getFormatter(executor, originalResult)), nil
	})
}

// Wraps (embeds) an executor and add specific result default output options getter.
func NewExecuteResultOutputOptions(
	executor Executor,
	getOutputOptions func(exec Executor, result Result) string,
) Executor {
	return NewExecuteResultWrapper(executor, func(wrapperExecutor ExecutorWrapper, originalResult Result) (wrappedResult Result, err error) {
		result, ok := ResultAs[ResultWithValue](originalResult)
		if !ok {
			return nil, fmt.Errorf("result is not core.ResultWithValue: %T %+v", originalResult, originalResult)
		}
		return NewResultWithDefaultOutputOptions(result, getOutputOptions(executor, originalResult)), nil
	})
}
