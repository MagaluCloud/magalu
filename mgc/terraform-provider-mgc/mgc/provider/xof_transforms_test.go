package provider

import (
	"encoding/json"
	"testing"

	mgcSchemaPkg "magalu.cloud/core/schema"
)

func marshalSchema(s *mgcSchemaPkg.Schema) string {
	marshaled, _ := json.Marshal(s)
	return string(marshaled)
}

func marshalXOfChildren(xOfChildren []xOfChild) []string {
	marshaledSchemas := make([]string, 0, len(xOfChildren))
	for _, xOf := range xOfChildren {
		marshaledSchema, _ := json.Marshal(xOf.s)
		marshaledSchemas = append(marshaledSchemas, string(marshaledSchema))
	}
	return marshaledSchemas
}

func TestCanPromoteXOfProps(t *testing.T) {
	type testCase struct {
		parent     *mgcSchemaPkg.Schema
		xOfs       []xOfChild
		canPromote bool
	}

	testCases := []testCase{
		{
			mgcSchemaPkg.NewAnyOfSchema(),
			[]xOfChild{
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"id":    mgcSchemaPkg.NewStringSchema(),
						"other": mgcSchemaPkg.NewIntegerSchema(),
					}, nil),
					"",
				},
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewStringSchema(),
					}, nil),
					"",
				},
			},
			true,
		},
		{
			mgcSchemaPkg.NewAnyOfSchema(),
			[]xOfChild{
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewStringSchema(),
					}, nil),
					"",
				},
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewStringSchema(),
					}, nil),
					"",
				},
			},
			true,
		},
		{
			mgcSchemaPkg.NewAnyOfSchema(),
			[]xOfChild{
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewStringSchema(),
					}, nil),
					"",
				},
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewIntegerSchema(),
					}, nil),
					"",
				},
			},
			false,
		},
	}

	for _, testCase := range testCases {
		parentCOW := mgcSchemaPkg.NewCOWSchema(testCase.parent)
		if canPromoteXOfChildrenProps(parentCOW, testCase.xOfs) != testCase.canPromote {
			t.Errorf("expected canPromote == %v for %s with children %s", testCase.canPromote, marshalSchema(testCase.parent), marshalXOfChildren(testCase.xOfs))
		}
	}
}

func TestPromoteXOfChildren(t *testing.T) {
	type testCase struct {
		parent       *mgcSchemaPkg.Schema
		xOfs         []xOfChild
		expected     *mgcSchemaPkg.Schema
		promotedKeys []string
	}

	testCases := []testCase{
		{
			mgcSchemaPkg.NewAnyOfSchema(),
			[]xOfChild{
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"id":    mgcSchemaPkg.NewStringSchema(),
						"other": mgcSchemaPkg.NewIntegerSchema(),
					}, nil),
					"",
				},
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewStringSchema(),
					}, nil),
					"",
				},
			},
			mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
				"id":    mgcSchemaPkg.NewStringSchema(),
				"other": mgcSchemaPkg.NewIntegerSchema(),
				"name":  mgcSchemaPkg.NewStringSchema(),
			}, nil),
			[]string{"id", "other", "name"},
		},
		{
			mgcSchemaPkg.NewAnyOfSchema(),
			[]xOfChild{
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewStringSchema(),
					}, nil),
					"",
				},
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewStringSchema(),
					}, nil),
					"",
				},
			},
			mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
				"name": mgcSchemaPkg.NewStringSchema(),
			}, nil),
			[]string{"name"},
		},
		{
			mgcSchemaPkg.NewObjectSchema(nil, nil),
			[]xOfChild{
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewStringSchema(),
					}, nil),
					"",
				},
				{
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewIntegerSchema(),
					}, nil),
					"",
				},
			},
			mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
				"object1": mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
					"name": mgcSchemaPkg.NewStringSchema(),
				}, nil),
				"object2": mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
					"name": mgcSchemaPkg.NewIntegerSchema(),
				}, nil),
			}, nil),
			[]string{"object1", "object2"},
		},
	}

	for _, testCase := range testCases {
		parentCOW := mgcSchemaPkg.NewCOWSchema(testCase.parent)
		err := promoteXOfChildren(parentCOW, testCase.xOfs)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !mgcSchemaPkg.CheckSimilarJsonSchemas(parentCOW.Peek(), testCase.expected) {
			t.Errorf("expected %s for %s with children %s, got %s", marshalSchema(testCase.expected), marshalSchema(testCase.parent), marshalXOfChildren(testCase.xOfs), marshalSchema(parentCOW.Peek()))
		}
	}
}

func TestGetXOfObjectSchemaTransformed(t *testing.T) {
	type testCase struct {
		input  *mgcSchemaPkg.Schema
		output *mgcSchemaPkg.Schema
	}

	testCases := []testCase{
		{
			mgcSchemaPkg.NewAnyOfSchema(
				mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
					"id":    mgcSchemaPkg.NewStringSchema(),
					"other": mgcSchemaPkg.NewIntegerSchema(),
				}, nil),
				mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
					"name": mgcSchemaPkg.NewStringSchema(),
				}, nil),
			),
			mgcSchemaPkg.NewObjectSchema(
				map[string]*mgcSchemaPkg.Schema{
					"id":    mgcSchemaPkg.NewStringSchema(),
					"other": mgcSchemaPkg.NewIntegerSchema(),
					"name":  mgcSchemaPkg.NewStringSchema(),
				},
				nil,
			),
		},
		{
			mgcSchemaPkg.NewAnyOfSchema(
				mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
					"name": mgcSchemaPkg.NewStringSchema(),
				}, nil),
				mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
					"name": mgcSchemaPkg.NewStringSchema(),
				}, nil),
			),
			mgcSchemaPkg.NewObjectSchema(
				map[string]*mgcSchemaPkg.Schema{
					"name": mgcSchemaPkg.NewStringSchema(),
				},
				nil,
			),
		},
		{
			mgcSchemaPkg.NewAnyOfSchema(
				mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
					"name": mgcSchemaPkg.NewStringSchema(),
				}, nil),
				mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
					"name": mgcSchemaPkg.NewIntegerSchema(),
				}, nil),
			),
			mgcSchemaPkg.NewObjectSchema(
				map[string]*mgcSchemaPkg.Schema{
					"object1": mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewStringSchema(),
					}, nil),
					"object2": mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"name": mgcSchemaPkg.NewIntegerSchema(),
					}, nil),
				},
				nil,
			),
		},
		{
			mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
				"prop": mgcSchemaPkg.NewAnyOfSchema(
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"propSub1": mgcSchemaPkg.NewStringSchema(),
					}, []string{"propSub1"}),
					mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
						"propSub2": mgcSchemaPkg.NewStringSchema(),
					}, []string{"propSub2"}),
				),
			}, []string{"prop"}),
			mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
				"prop": mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
					"propSub1": mgcSchemaPkg.NewStringSchema(),
					"propSub2": mgcSchemaPkg.NewStringSchema(),
				}, nil),
			}, []string{"prop"}),
		},
	}

	for _, testCase := range testCases {
		transformed, err := getXOfObjectSchemaTransformed(testCase.input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !mgcSchemaPkg.CheckSimilarJsonSchemas(transformed, testCase.output) {
			t.Errorf("expected %s for %s, but got %s", marshalSchema(testCase.output), marshalSchema(testCase.input), marshalSchema(transformed))
		}
	}
}
