package schema

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/jsonschema"
)

func ToCoreSchema(s *jsonschema.Schema) (*Schema, error) {
	if s == nil {
		return nil, fmt.Errorf("invalid jsonschema.Schema passed to 'toCoreSchema' function")
	}

	rootSchema := s
	if s.Ref != "" {
		rootDef, ok := lookupDefByPath(s.Definitions, getRefPath(s.Ref))
		if !ok {
			return nil, fmt.Errorf("unable to resolve reference %s when generating schema via 'toCoreSchema'", s.Ref)
		}
		rootSchema = rootDef
	}

	oapiSchema, err := unmarshallIntoOpenapiSchema(rootSchema)

	if err != nil {
		return nil, err
	}

	refCache := map[string]*openapi3.Schema{}

	err = resolveRefs(oapiSchema, s.Definitions, refCache)
	if err != nil {
		return nil, err
	}

	initializeExtensions(oapiSchema)
	return (*Schema)(oapiSchema), nil
}

func unmarshallIntoOpenapiSchema(s *jsonschema.Schema) (*openapi3.Schema, error) {
	json, err := s.MarshalJSON()
	if err != nil {
		return nil, err
	}

	newS := &openapi3.Schema{}
	err = newS.UnmarshalJSON(json)
	if err != nil {
		return nil, err
	}

	return newS, nil
}

func getRefPath(ref string) []string {
	ref = strings.TrimPrefix(ref, "#/$defs/")
	return strings.Split(ref, "/")
}

func lookupDefByPath(defs jsonschema.Definitions, path []string) (*jsonschema.Schema, bool) {
	pathLength := len(path)
	if pathLength == 0 {
		return nil, false
	}

	def, ok := defs[path[0]]
	if !ok {
		return nil, false
	}

	if pathLength > 1 {
		if def.Definitions != nil {
			return lookupDefByPath(def.Definitions, path[1:])
		} else {
			return nil, false
		}
	} else {
		return def, true
	}
}

func initializeExtensions(s *openapi3.Schema) {
	s.Extensions = make(map[string]any)
	_, _ = visitAllSubRefs(s, func(ref *openapi3.SchemaRef) error {
		ref.Value.Extensions = make(map[string]any)

		return nil
	})
}

type subRefVisitor func(ref *openapi3.SchemaRef) error

func visitAllSubRefs(s *openapi3.Schema, visitor subRefVisitor) (bool, error) {
	if s == nil {
		return true, nil
	}

	if s.AdditionalProperties.Schema != nil {
		if err := visitor(s.AdditionalProperties.Schema); err != nil {
			return false, err
		}
		if shouldContinue, err := visitAllSubRefs(s.AdditionalProperties.Schema.Value, visitor); !shouldContinue {
			return false, err
		}
	}

	if s.OneOf != nil {
		for _, s := range s.OneOf {
			if err := visitor(s); err != nil {
				return false, err
			}
			if shouldContinue, err := visitAllSubRefs(s.Value, visitor); !shouldContinue {
				return false, err
			}
		}
	}
	if s.AnyOf != nil {
		for _, s := range s.AnyOf {
			if err := visitor(s); err != nil {
				return false, err
			}
			if shouldContinue, err := visitAllSubRefs(s.Value, visitor); !shouldContinue {
				return false, err
			}
		}
	}
	if s.AllOf != nil {
		for _, s := range s.AllOf {
			if err := visitor(s); err != nil {
				return false, err
			}
			if shouldContinue, err := visitAllSubRefs(s.Value, visitor); !shouldContinue {
				return false, err
			}
		}
	}
	if s.Properties != nil {
		for _, s := range s.Properties {
			if err := visitor(s); err != nil {
				return false, err
			}
			if shouldContinue, err := visitAllSubRefs(s.Value, visitor); !shouldContinue {
				return false, err
			}
		}
	}
	if s.Not != nil {
		if err := visitor(s.Not); err != nil {
			return false, err
		}
		if shouldContinue, err := visitAllSubRefs(s.Not.Value, visitor); !shouldContinue {
			return false, err
		}
	}
	if s.Items != nil {
		if err := visitor(s.Items); err != nil {
			return false, err
		}
		if shouldContinue, err := visitAllSubRefs(s.Items.Value, visitor); !shouldContinue {
			return false, err
		}
	}
	return true, nil
}

func resolveRefs(s *openapi3.Schema, defs jsonschema.Definitions, cache map[string]*openapi3.Schema) error {
	_, err := visitAllSubRefs(s, func(ref *openapi3.SchemaRef) error {
		if ref == nil || ref.Value != nil {
			return nil
		}

		if ref.Ref == "" {
			return fmt.Errorf("schema with empty reference passed to 'toCoreSchema'")
		}

		if def, ok := findRef(getRefPath(ref.Ref), defs, cache); ok {
			ref.Value = def
			return resolveRefs(def, defs, cache)
		}

		return fmt.Errorf("unable to resolve %s reference when generating schema via 'toCoreSchema' function", ref.Ref)
	})
	return err
}

func findRef(refPath []string, defs jsonschema.Definitions, cache map[string]*openapi3.Schema) (*openapi3.Schema, bool) {
	refPathLength := len(refPath)
	if refPathLength == 0 {
		return nil, false
	}

	refName := strings.Join(refPath, "/")
	if def, ok := cache[refName]; ok {
		return def, true
	}

	def, ok := lookupDefByPath(defs, refPath)
	if !ok {
		return nil, false
	}

	oapiDef, err := unmarshallIntoOpenapiSchema(def)
	if err != nil {
		return nil, false
	}

	err = resolveRefs(oapiDef, defs, cache)
	if err != nil {
		return nil, false
	}

	cache[refName] = oapiDef
	return oapiDef, true
}
