package transform

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/maps"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
)

type transformer interface {
	TransformValue(value any) (any, error)
	TransformSchema(value *mgcSchemaPkg.COWSchema) (*mgcSchemaPkg.COWSchema, error)
}

type transformerCreator func(spec *transformSpec) (transformer, error)

type transformerRegistryEntry struct {
	name    string
	schema  *mgcSchemaPkg.Schema
	creator transformerCreator
}

var registryInits = map[string]func() *transformerRegistryEntry{}
var m sync.RWMutex
var registry = utils.NewLazyLoader(func() map[string]*transformerRegistryEntry {
	reg := map[string]*transformerRegistryEntry{}
	m.RLock()
	for n, init := range registryInits {
		reg[n] = init()
	}
	m.RUnlock()
	maps.Clear(registryInits)
	return reg
})

type transformSpec struct {
	Type string `json:"type"`
	// See more about the 'remain' directive here: https://pkg.go.dev/github.com/mitchellh/mapstructure#hdr-Remainder_Values
	Parameters map[string]any `json:",remain"` // nolint
	Schema     *openapi3.Schema
}

func newTransformSpecFromString(s string) *transformSpec {
	if len(s) == 0 {
		return nil
	}
	return &transformSpec{Type: s}
}

func newTransformSpecFromMap(m map[string]any) *transformSpec {
	if len(m) == 0 {
		return nil
	}
	spec, err := utils.DecodeNewValue[transformSpec](m)
	if err != nil || len(spec.Type) == 0 {
		return nil
	}
	return spec
}

// spec must be string or map
func newTransformSpec(spec any) *transformSpec {
	if s, ok := spec.(string); ok {
		return newTransformSpecFromString(s)
	} else if m, ok := spec.(map[string]any); ok {
		return newTransformSpecFromMap(m)
	}
	return nil
}

func newTransformSpecSlice(sl []any) []*transformSpec {
	ret := make([]*transformSpec, 0, len(sl))
	for _, spec := range sl {
		if ts := newTransformSpec(spec); ts != nil {
			ret = append(ret, ts)
		}
	}
	if len(ret) == 0 {
		return nil
	}
	return ret
}

func getTransformationSpecs(extensions map[string]any, transformationKey string) []*transformSpec {
	if spec, ok := extensions[transformationKey]; !ok {
		return nil
	} else if sl, ok := spec.([]any); ok {
		return newTransformSpecSlice(sl)
	} else if ts := newTransformSpec(spec); ts != nil {
		return []*transformSpec{ts}
	} else {
		return nil
	}
}

func getTransformers(extensions map[string]any, transformationKey string) (transformers []transformer, err error) {
	specs := getTransformationSpecs(extensions, transformationKey)
	for i, spec := range specs {
		entry := registry()[spec.Type]
		if entry == nil {
			return nil, &utils.ChainedError{Name: fmt.Sprint(i), Err: fmt.Errorf("unknown transformer type: %q", spec.Type)}
		}
		t, err := entry.creator(spec)
		if err != nil {
			return nil, &utils.ChainedError{
				Name: fmt.Sprint(i), Err: &utils.ChainedError{
					Name: spec.Type,
					Err:  fmt.Errorf("failed to create: %#v", err),
				},
			}
		}
		transformers = append(transformers, t)
	}
	return
}

var baseSchema = utils.NewLazyLoader(func() *mgcSchemaPkg.Schema {
	schema, err := mgcSchemaPkg.SchemaFromType[transformSpec]()
	if err != nil {
		panic(err.Error())
	}
	return schema
})

func addTransformer[T any](creator transformerCreator, names ...string) {
	if len(names) == 0 {
		panic("programming error: missing names")
	}
	rt := reflect.TypeOf(*new(T))

	initSchema := func() *mgcSchemaPkg.Schema {
		var schema *mgcSchemaPkg.Schema = baseSchema()
		if rt.Kind() == reflect.Struct && rt.NumField() > 0 {
			transformerSchema, err := mgcSchemaPkg.SchemaFromType[T]()
			if err != nil {
				panic(err.Error())
			}
			schema, err = mgcSchemaPkg.SimplifySchema(mgcSchemaPkg.NewAllOfSchema(schema, transformerSchema))
			if err != nil {
				panic(err.Error())
			}
		}
		return schema
	}

	for _, n := range names {
		if existing := registryInits[n]; existing != nil {
			panic(fmt.Sprintf("existing transformer: %q => creator=%#v, schema=%#v\n", n, existing().creator, existing().schema))
		}
		registryInits[n] = func() *transformerRegistryEntry {
			return &transformerRegistryEntry{n, initSchema(), creator}
		}
	}
}
