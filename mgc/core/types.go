package core

import (
	"context"
	"fmt"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

// NOTE: TODO: should we duplicate this, or find a more generic package?
type Schema openapi3.Schema

func (s *Schema) VisitJSON(value any, opts ...openapi3.SchemaValidationOption) error {
	return (*openapi3.Schema)(s).VisitJSON(value, opts...)
}

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

type Executor interface {
	Descriptor
	ParametersSchema() *Schema
	ConfigsSchema() *Schema
	ResultSchema() *Schema
	// The maps for the parameters and configs should NOT be modified inside the implementation of 'Execute'
	Execute(context context.Context, parameters map[string]Value, configs map[string]Value) (result Value, err error)
}

type ExecutorWrapper interface {
	Unwrap() Executor
}

func ExecutorAs[T Executor](exec Executor) (T, bool) {
	var zeroT T

	if u, ok := exec.(ExecutorWrapper); ok {
		t, ok := u.(T)
		if ok {
			return t, true
		}

		x := u.Unwrap()
		if x == nil {
			return zeroT, false
		}

		return ExecutorAs[T](t)
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

// Implement this interface in Executor()s that want to provide customized formatting of output.
// It's used by the command line interface (CLI) and possible other tools.
// This is only called if no other explicit formatting is desired
type ExecutorResultFormatter interface {
	Executor
	// NOTE: result is the converted value, such as primitives, map[string]any, []any...
	// Whenever using StaticExecute, it's *NOT* the ResultT (ie: struct)
	DefaultFormatResult(result Value) string
}

type executeFormat struct {
	Executor
	formatter func(result Value) string
}

func (o *executeFormat) DefaultFormatResult(result Value) string {
	return o.formatter(result)
}

var _ ExecutorResultFormatter = (*executeFormat)(nil)

// Wraps (embeds) an executor and add specific result formatting.
func NewExecuteFormat(
	executor Executor,
	formatter func(result Value) string,
) ExecutorResultFormatter {
	return &executeFormat{executor, formatter}
}

// Implement this interface in Executor()s that want to provide default output options.
// It's used by the command line interface (CLI) and possible other tools.
// This is only called if no other explicit options are desired
type ExecutorResultOutputOptions interface {
	Executor
	// The return should be in the same format as CLI -o "VALUE"
	// example: "yaml" or "table=COL:$.path.to[*].element,OTHERCOL:$.path.to[*].other"
	DefaultOutputOptions(result Value) string
}

type executeResultOutputOptions struct {
	Executor
	getOutputOptions func(exec Executor, result Value) string
}

func (o *executeResultOutputOptions) Unwrap() Executor {
	return o.Executor
}

func (o *executeResultOutputOptions) DefaultOutputOptions(result Value) string {
	return o.getOutputOptions(o.Executor, result)
}

var _ ExecutorResultOutputOptions = (*executeResultOutputOptions)(nil)

// Wraps (embeds) an executor and add specific result default output options getter.
func NewExecuteResultOutputOptions(
	executor Executor,
	getOutputOptions func(exec Executor, result Value) string,
) ExecutorResultOutputOptions {
	return &executeResultOutputOptions{executor, getOutputOptions}
}

type LinksVisitor func(link Linker) (run bool, err error)

type Linker interface {
	Name() string
	// The executor to use to run this link. Prepare it's parameters and config using PrepareLink()
	LinkTarget() Executor
	// The link will prepare the parameters and config based on the original executor's parameter/config and result
	// this will return the new parameters/config to give to LinkTarget()
	PrepareLink(originalParameters map[string]Value, originalConfigs map[string]Value, originalResult Value) (preparedParameters map[string]Value, preparedConfigs map[string]Value, err error)
}

type LinkedExecutor interface {
	Executor
	VisitLinks(visitor LinksVisitor) (finished bool, err error)
	GetLinkByName(name string) (child Descriptor, err error)
}

func ExecuteLink(link Linker, ctx context.Context, originalParameters map[string]Value, originalConfigs map[string]Value, originalResult Value) (result Value, err error) {
	preparedParameters, preparedConfigs, err := link.PrepareLink(originalParameters, originalConfigs, originalResult)
	if err != nil {
		return nil, err
	}

	exec := link.LinkTarget()
	return exec.Execute(ctx, preparedParameters, preparedConfigs)
}
