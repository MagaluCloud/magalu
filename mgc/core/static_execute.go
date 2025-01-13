package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

// Given a base spec and types to be reflected, populate with the schemas
func ReflectExecutorSpecSchemas[ParamsT any, ConfigsT any, ResultT any](baseSpec ExecutorSpec) (spec ExecutorSpec, err error) {
	spec = baseSpec

	spec.ParametersSchema, err = mgcSchemaPkg.SchemaFromType[ParamsT]()
	if err != nil {
		err = &ChainedError{Name: "ParamsT", Err: err}
		return
	}
	spec.ConfigsSchema, err = mgcSchemaPkg.SchemaFromType[ConfigsT]()
	if err != nil {
		err = &ChainedError{Name: "ConfigsT", Err: err}
		return
	}
	spec.ResultSchema, err = mgcSchemaPkg.SchemaFromType[ResultT]()
	if err != nil {
		err = &ChainedError{Name: "ResultT", Err: err}
		return
	}

	spec.PositionalArgs = getPositionals(reflect.TypeOf(new(ParamsT)))
	spec.HiddenFlags = getHiddens(reflect.TypeOf(new(ParamsT)))

	return spec, nil
}

func getHiddens(t reflect.Type) []string {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	var hidden []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			hidden = append(hidden, getHiddens(field.Type)...)
			continue
		}

		if !GetMgcTagBool(field.Tag, "hidden") {
			continue
		}

		name := field.Name
		if jsonName := strings.Split(field.Tag.Get("json"), ",")[0]; jsonName != "" {
			name = jsonName
		}

		hidden = append(hidden, name)
	}

	return hidden
}

func getPositionals(t reflect.Type) []string {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	var positionals []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			positionals = append(positionals, getPositionals(field.Type)...)
			continue
		}

		if !GetMgcTagBool(field.Tag, "positional") {
			continue
		}

		name := field.Name
		if jsonName := strings.Split(field.Tag.Get("json"), ",")[0]; jsonName != "" {
			name = jsonName
		}

		positionals = append(positionals, name)
	}

	return positionals
}

func ReflectExecutorSpecFn[ParamsT any, ConfigsT any, ResultT any](
	typedExecute func(context context.Context, params ParamsT, configs ConfigsT) (result ResultT, err error),
) ExecutorSpecFn {
	return func(executor Executor, ctx context.Context, parameters Parameters, configs Configs) (Result, error) {
		typedParams, err := utils.DecodeNewValue[ParamsT](parameters)
		if err != nil {
			return nil, &ChainedError{Name: "parameters", Err: fmt.Errorf("decoding error. Did you forget to set 'json' struct flags for struct %T?: %w", typedParams, err)}
		}

		typedConfigs, err := utils.DecodeNewValue[ConfigsT](configs)
		if err != nil {
			return nil, &ChainedError{Name: "configs", Err: fmt.Errorf("decoding error. Did you forget to set 'json' struct flags for struct %T?: %w", typedConfigs, err)}
		}

		typedResult, err := typedExecute(ctx, *typedParams, *typedConfigs)
		if err != nil {
			return nil, err
		}

		value, err := utils.SimplifyAny(typedResult)
		if err != nil {
			return nil, &ChainedError{Name: "result", Err: fmt.Errorf("error simplifying %T: %w", typedResult, err)}
		}

		source := ResultSource{
			Executor:   executor,
			Context:    ctx,
			Parameters: parameters,
			Configs:    configs,
		}
		return NewSimpleResult(source, executor.ResultSchema(), value), nil
	}
}

func ReflectExecutorSpec[ParamsT any, ConfigsT any, ResultT any](
	baseSpec ExecutorSpec,
	typedExecute func(context context.Context, params ParamsT, configs ConfigsT) (result ResultT, err error),
) (spec ExecutorSpec, err error) {
	spec, err = ReflectExecutorSpecSchemas[ParamsT, ConfigsT, ResultT](baseSpec)
	if err != nil {
		return
	}

	spec.Execute = ReflectExecutorSpecFn[ParamsT, ConfigsT, ResultT](typedExecute)
	return
}

// Go Parameter and Config structs
// Note: we use both 'jsonschema' and 'mapstructure' for this helper. Be careful
// when using struct tags in your Params and Configs structs, as the tags from those
// libraries can't be out of sync when it comes to field names/json names
// See:
// - https://pkg.go.dev/github.com/invopop/jsonschema
// - https://pkg.go.dev/github.com/mitchellh/mapstructure
func NewReflectedSimpleExecutor[ParamsT any, ConfigsT any, ResultT any](
	baseSpec ExecutorSpec,
	execute func(context context.Context, params ParamsT, configs ConfigsT) (result ResultT, err error),
) *SimpleExecutor {
	execSpec, err := ReflectExecutorSpec[ParamsT, ConfigsT, ResultT](baseSpec, execute)
	if err != nil {
		logger().Fatalw("cannot reflect static executor", "err", err)
	}
	return NewSimpleExecutor(execSpec)
}

// Version that generates the schema, but uses baseSpec.Execute as is
func NewReflectedSimpleExecutorSchemas[ParamsT any, ConfigsT any, ResultT any](
	baseSpec ExecutorSpec,
) *SimpleExecutor {
	execSpec, err := ReflectExecutorSpecSchemas[ParamsT, ConfigsT, ResultT](baseSpec)
	if err != nil {
		logger().Fatalw("cannot reflect static executor", "err", err)
	}
	return NewSimpleExecutor(execSpec)
}

// Go Parameter and Config structs
// Note: we use both 'jsonschema' and 'mapstructure' for this helper. Be careful
// when using struct tags in your Params and Configs structs, as the tags from those
// libraries can't be out of sync when it comes to field names/json names
// See:
// - https://pkg.go.dev/github.com/invopop/jsonschema
// - https://pkg.go.dev/github.com/mitchellh/mapstructure
func NewStaticExecute[ParamsT any, ConfigsT any, ResultT any](
	spec DescriptorSpec,
	execute func(context context.Context, params ParamsT, configs ConfigsT) (result ResultT, err error),
) *SimpleExecutor {
	return NewReflectedSimpleExecutor(ExecutorSpec{DescriptorSpec: spec}, execute)
}

// No parameters or configs
func NewStaticExecuteSimple[ResultT any](
	spec DescriptorSpec,
	execute func(context context.Context) (result ResultT, err error),
) *SimpleExecutor {
	return NewReflectedSimpleExecutor(
		ExecutorSpec{DescriptorSpec: spec},
		func(context context.Context, _ struct{}, _ struct{}) (result ResultT, err error) {
			return execute(context)
		},
	)
}
