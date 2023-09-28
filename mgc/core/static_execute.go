package core

import (
	"context"
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/jsonschema"
	"go.uber.org/zap"
	coreLogger "magalu.cloud/core/logger"
	"magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
)

type StaticExecute struct {
	name        string
	version     string
	description string
	parameters  *Schema
	config      *Schema
	result      *Schema
	links       map[string]Linker
	execute     func(ctx context.Context, parameters Parameters, configs Configs) (value Value, err error)
}

var corePkgLogger *zap.SugaredLogger

func logger() *zap.SugaredLogger {
	if corePkgLogger == nil {
		corePkgLogger = coreLogger.New[StaticExecute]()
	}
	return corePkgLogger
}

// Raw Parameter and Config JSON Schemas
func NewRawStaticExecute(name string, version string, description string, parameters *Schema, config *Schema, result *Schema, links map[string]Linker, execute func(context context.Context, parameters Parameters, configs Configs) (value Value, err error)) *StaticExecute {
	return &StaticExecute{name, version, description, parameters, config, result, links, execute}
}

func newAnySchema() *Schema {
	s := &openapi3.Schema{
		Nullable: true,
		AnyOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{Value: &openapi3.Schema{Type: "null", Nullable: true}},
			&openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
			&openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
			&openapi3.SchemaRef{Value: openapi3.NewFloat64Schema()},
			&openapi3.SchemaRef{Value: openapi3.NewIntegerSchema()},
			&openapi3.SchemaRef{Value: openapi3.NewArraySchema()},
			&openapi3.SchemaRef{Value: openapi3.NewObjectSchema().WithAnyAdditionalProperties()},
		},
	}

	return (*Schema)(s)
}

func schemaFromType[T any]() (*Schema, error) {
	t := new(T)
	tp := reflect.TypeOf(t).Elem()
	kind := tp.Kind()
	if tp.Name() == "" && kind == reflect.Interface {
		return newAnySchema(), nil
	}

	s, err := schema.ToCoreSchema(schemaReflector.Reflect(t))
	if err != nil {
		return nil, fmt.Errorf("unable to create JSON Schema for type '%T': %w", t, err)
	}

	isArray := kind == reflect.Array || kind == reflect.Slice

	// schemaReflector seems to lose the fact that it's an array, so we bring that back
	if isArray && s.Type == "object" {
		arrSchema := schema.NewArraySchema(s)
		s = arrSchema
	}

	return s, nil
}

// Go Parameter and Config structs
// Note: we use both 'jsonschema' and 'mapstructure' for this helper. Be careful
// when using struct tags in your Params and Configs structs, as the tags from those
// libraries can't be out of sync when it comes to field names/json names
// See:
// - https://pkg.go.dev/github.com/invopop/jsonschema
// - https://pkg.go.dev/github.com/mitchellh/mapstructure
func NewStaticExecuteWithLinks[ParamsT any, ConfigsT any, ResultT any](
	name string,
	version string,
	description string,
	links map[string]Linker,
	execute func(context context.Context, params ParamsT, configs ConfigsT) (result ResultT, err error),
) *StaticExecute {
	ps, err := schemaFromType[ParamsT]()
	if err != nil {
		logger().Fatal(err)
	}
	cs, err := schemaFromType[ConfigsT]()
	if err != nil {
		logger().Fatal(err)
	}
	rs, err := schemaFromType[ResultT]()
	if err != nil {
		logger().Fatal(err)
	}

	return NewRawStaticExecute(
		name,
		version,
		description,
		ps,
		cs,
		rs,
		links,
		func(ctx context.Context, parameters Parameters, configs Configs) (Value, error) {
			paramsStruct, err := utils.DecodeNewValue[ParamsT](parameters)
			if err != nil {
				return nil, fmt.Errorf("error when decoding parameters. Did you forget to set 'json' struct flags for struct %T?: %w", paramsStruct, err)
			}

			configsStruct, err := utils.DecodeNewValue[ConfigsT](configs)
			if err != nil {
				return nil, fmt.Errorf("error when decoding configs. Did you forget to set 'json' struct flags for struct %T?: %w", paramsStruct, err)
			}

			value, err := execute(ctx, *paramsStruct, *configsStruct)
			if err != nil {
				return nil, err
			}

			return utils.SimplifyAny(value)
		},
	)
}

// Go Parameter and Config structs
// Note: we use both 'jsonschema' and 'mapstructure' for this helper. Be careful
// when using struct tags in your Params and Configs structs, as the tags from those
// libraries can't be out of sync when it comes to field names/json names
// See:
// - https://pkg.go.dev/github.com/invopop/jsonschema
// - https://pkg.go.dev/github.com/mitchellh/mapstructure
func NewStaticExecute[ParamsT any, ConfigsT any, ResultT any](
	name string,
	version string,
	description string,
	execute func(context context.Context, params ParamsT, configs ConfigsT) (result ResultT, err error),
) *StaticExecute {
	return NewStaticExecuteWithLinks(name, version, description, nil, execute)
}

// No parameters or configs
func NewStaticExecuteSimpleWithLinks[ResultT any](
	name string,
	version string,
	description string,
	links map[string]Linker,
	execute func(ctx context.Context) (result ResultT, err error),
) *StaticExecute {
	return NewStaticExecuteWithLinks(
		name,
		version,
		description,
		links,
		func(ctx context.Context, _, _ struct{}) (ResultT, error) {
			return execute(ctx)
		},
	)
}

// No parameters or configs
func NewStaticExecuteSimple[ResultT any](
	name string,
	version string,
	description string,
	execute func(context context.Context) (result ResultT, err error),
) *StaticExecute {
	return NewStaticExecuteSimpleWithLinks(name, version, description, nil, execute)
}

// BEGIN: Descriptor interface:

func (o *StaticExecute) Name() string {
	return o.name
}

func (o *StaticExecute) Version() string {
	return o.version
}

func (o *StaticExecute) Description() string {
	return o.description
}

// END: Descriptor interface

// BEGIN: Executor interface:

func (o *StaticExecute) ParametersSchema() *Schema {
	return o.parameters
}

func (o *StaticExecute) ConfigsSchema() *Schema {
	return o.config
}

func (o *StaticExecute) ResultSchema() *Schema {
	return o.result
}

func (o *StaticExecute) Execute(context context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	value, err := o.execute(context, parameters, configs)
	if err != nil {
		return nil, err
	}
	source := ResultSource{
		Executor:   o,
		Context:    context,
		Parameters: parameters,
		Configs:    configs,
	}
	return NewSimpleResult(source, o.result, value), nil
}

func (o *StaticExecute) Links() map[string]Linker {
	return o.links
}

var _ Executor = (*StaticExecute)(nil)

// END: Executor interface

var schemaReflector *jsonschema.Reflector

func init() {
	schemaReflector = &jsonschema.Reflector{
		DoNotReference: false,
	}
}
