package core

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/invopop/jsonschema"
	"github.com/mitchellh/mapstructure"
)

type StaticExecute struct {
	name        string
	version     string
	description string
	parameters  *Schema
	config      *Schema
	result      *Schema
	execute     func(ctx context.Context, parameters map[string]Value, configs map[string]Value) (result Value, err error)
}

// Raw Parameter and Config JSON Schemas
func NewRawStaticExecute(name string, version string, description string, parameters *Schema, config *Schema, result *Schema, execute func(context context.Context, parameters map[string]Value, configs map[string]Value) (result Value, err error)) *StaticExecute {
	return &StaticExecute{name, version, description, parameters, config, result, execute}
}

func convertValue(value Value) (converted Value, err error) {
	if value == nil {
		return nil, nil
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return value, nil

	case reflect.Invalid,
		reflect.Chan,
		reflect.Complex64, reflect.Complex128,
		reflect.Func,
		reflect.UnsafePointer,
		reflect.Interface:
		return nil, fmt.Errorf("Forbidden value type %s", v)

	case reflect.Pointer:
		if v.IsNil() {
			return nil, nil
		}
		return convertValue(v.Elem().Interface())

	case reflect.Array, reflect.Slice:
		if result, ok := value.([]Value); ok {
			return result, nil
		}
		// convert whatever map to []Value
		count := v.Len()
		result := make([]Value, count)
		for i := 0; i < count; i++ {
			subVal := v.Index(i).Interface()
			decoded := map[string]any{}
			err := mapstructure.Decode(subVal, &decoded)
			if err != nil {
				return nil, err
			}
			result[i] = decoded
		}
		return result, err

	case reflect.Map:
		if resultMap, ok := value.(map[string]Value); ok {
			return resultMap, nil
		}
		// convert whatever map to map[string]Value
		resultMap := make(map[string]Value, v.Len())
		err = mapstructure.Decode(value, &resultMap)
		return resultMap, err

	case reflect.Struct:
		resultMap := map[string]Value{}
		err = mapstructure.Decode(value, &resultMap)
		return resultMap, err

	default:
		return nil, fmt.Errorf("Unhandled value type: %s", v)
	}
}

func schemaFromType[T any]() (*Schema, error) {
	t := new(T)
	s, err := toCoreSchema(schemaReflector.Reflect(t))
	if err != nil {
		return nil, fmt.Errorf("unable to create JSON Schema for type '%T': %w", t, err)
	}

	kind := reflect.TypeOf(t).Elem().Kind()
	isArray := kind == reflect.Array || kind == reflect.Slice

	// schemaReflector seems to lose the fact that it's an array, so we bring that back
	if isArray && s.Type == "object" {
		arrSchema := NewArraySchema(s)
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
func NewStaticExecute[ParamsT any, ConfigsT any, ResultT any](
	name string,
	version string,
	description string,
	execute func(context context.Context, params ParamsT, configs ConfigsT) (result ResultT, err error),
) *StaticExecute {
	ps, err := schemaFromType[ParamsT]()
	if err != nil {
		log.Fatal(err)
	}
	cs, err := schemaFromType[ConfigsT]()
	if err != nil {
		log.Fatal(err)
	}
	rs, err := schemaFromType[ResultT]()
	if err != nil {
		log.Fatal(err)
	}

	return NewRawStaticExecute(
		name,
		version,
		description,
		ps,
		cs,
		rs,
		func(ctx context.Context, parameters, configs map[string]any) (Value, error) {
			var paramsStruct ParamsT
			var configsStruct ConfigsT

			err := mapstructure.Decode(parameters, &paramsStruct)
			if err != nil {
				return nil, fmt.Errorf("error when decoding parameters. Did you forget to set 'mapstructure' struct flags?: %w", err)
			}

			err = mapstructure.Decode(configs, &configsStruct)
			if err != nil {
				return nil, fmt.Errorf("error when decoding configs. Did you forget to set 'mapstructure' struct flags?: %w", err)
			}

			result, err := execute(ctx, paramsStruct, configsStruct)
			if err != nil {
				return nil, err
			}

			return convertValue(result)
		},
	)
}

// No parameters or configs
func NewStaticExecuteSimple[ResultT any](
	name string,
	version string,
	description string,
	execute func(ctx context.Context) (result ResultT, err error),
) *StaticExecute {
	return NewStaticExecute(
		name,
		version,
		description,
		func(ctx context.Context, _, _ struct{}) (ResultT, error) {
			return execute(ctx)
		},
	)
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

func (o *StaticExecute) Execute(context context.Context, parameters map[string]Value, configs map[string]Value) (result Value, err error) {
	return o.execute(context, parameters, configs)
}

var _ Executor = (*StaticExecute)(nil)

// END: Executor interface

var schemaReflector *jsonschema.Reflector

func init() {
	schemaReflector = new(jsonschema.Reflector)
}
