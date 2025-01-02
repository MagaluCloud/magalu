package schema_flags

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
)

func checkError(t *testing.T, expected, got error) {
	if expected != nil {
		if got == nil {
			t.Errorf("expected error %#v, got none", expected)
		} else if !errors.Is(expected, got) {
			t.Errorf("expected error %#v, got %#v", expected, got)
		}
	} else if got != nil {
		t.Errorf("got unexpected error %#v", got)
	}
}

func checkArray(t *testing.T, expected, got []any) {
	if len(expected) != len(got) {
		t.Errorf("expected array length %d, got %d", len(expected), len(got))
		return
	}

	for i, vExpected := range expected {
		vGot := got[i]
		if !reflect.DeepEqual(vExpected, vGot) {
			// print type as it's usually the case, int x int64 x float64...
			t.Errorf("expected index %d:\nexpected: %#v (%T)\ngot.....: %#v (%T)", i, vExpected, vExpected, vGot, vGot)
		}
	}
}

func Test_parseStringItem(t *testing.T) {
	type testCase struct {
		name      string
		rawValue  string
		expected  string
		remaining string
		err       error
	}

	tests := []testCase{
		{
			name:     "empty",
			rawValue: "",
			expected: "",
		},
		{
			name:     "just spaces",
			rawValue: "   ",
			expected: "",
		},
		{
			name:     "word",
			rawValue: "word",
			expected: "word",
		},
		{
			name:     "leading spaces",
			rawValue: "   word",
			expected: "word",
		},
		{
			name:     "trailing spaces",
			rawValue: "word   ",
			expected: "word",
		},
		{
			name:     "inner spaces",
			rawValue: "   keep  inner\t white spaces   ",
			expected: "keep  inner\t white spaces",
		},
		{
			name:      "stop at CSV delimiter",
			rawValue:  "a,remaining",
			expected:  "a",
			remaining: "remaining",
		},
		{
			name:     "simple quoted",
			rawValue: "\"word\"",
			expected: "word",
		},
		{
			name:     "quotes with escape",
			rawValue: "\"word \\\"escaped quotes\\\" usage\"",
			expected: "word \"escaped quotes\" usage",
		},
		{
			name:      "quoted stop at CSV delimiter",
			rawValue:  "\"word\"   ,remaining",
			expected:  "word",
			remaining: "remaining",
		},
	}

	for _, tc := range tests {
		name := tc.name
		if name == "" {
			tc.name = tc.rawValue
		}
		t.Run(name, func(t *testing.T) {
			got, end, err := parseStringItem(tc.rawValue)
			checkError(t, tc.err, err)
			if tc.expected != got {
				t.Errorf("expected value %q, got %#v", tc.expected, got)
			}

			remaining := tc.rawValue[end:]
			if tc.remaining != remaining {
				t.Errorf("expected remaining value %q, got %q", tc.remaining, remaining)
			}
		})
	}
}

func Test_parseAnyItem(t *testing.T) {
	type testCase struct {
		name      string
		rawValue  string
		expected  any
		remaining string
		err       error
	}

	tests := []testCase{
		{
			name:     "empty",
			rawValue: "",
			expected: nil,
		},
		{
			name:     "just spaces",
			rawValue: "   ",
			expected: nil,
		},
		{
			name:     "word",
			rawValue: "word",
			expected: "word",
		},
		{
			name:     "leading spaces",
			rawValue: "   word",
			expected: "word",
		},
		{
			name:     "trailing spaces",
			rawValue: "word   ",
			expected: "word",
		},
		{
			name:     "simple quoted",
			rawValue: "\"word\"",
			expected: "word",
		},
		{
			name:     "quotes with escape",
			rawValue: "\"word \\\"escaped quotes\\\" usage\"",
			expected: "word \"escaped quotes\" usage",
		},
		{
			name:      "quoted stop at CSV delimiter",
			rawValue:  "\"word\"   ,remaining",
			expected:  "word",
			remaining: "remaining",
		},
		{
			name:      "csv number",
			rawValue:  "  1   ,remaining",
			expected:  1.0,
			remaining: "remaining",
		},
		{
			name:      "csv bool",
			rawValue:  "  true   ,remaining",
			expected:  true,
			remaining: "remaining",
		},
		{
			name:      "csv with json array",
			rawValue:  "   [true, 1, \"string\"]   ,remaining",
			expected:  []any{true, 1.0, "string"},
			remaining: "remaining",
		},
		{
			name:      "csv with json object",
			rawValue:  "  { \"k\": 1 }   ,remaining",
			expected:  map[string]any{"k": 1.0},
			remaining: "remaining",
		},
	}

	for _, tc := range tests {
		name := tc.name
		if name == "" {
			tc.name = tc.rawValue
		}
		t.Run(name, func(t *testing.T) {
			got, end, err := parseAnyItem(tc.rawValue)
			checkError(t, tc.err, err)
			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("expected value %#v (%T), got %#v (%T)", tc.expected, tc.expected, got, got)
			}

			remaining := tc.rawValue[end:]
			if tc.remaining != remaining {
				t.Errorf("expected remaining value %q, got %q", tc.remaining, remaining)
			}
		})
	}
}

func Test_parseArrayFlagValueSingle(t *testing.T) {
	type testCase struct {
		name        string
		itemsSchema *core.Schema
		rawValue    string
		expected    []any
		err         error
	}

	tests := []testCase{
		{
			name:        "json",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    "[11,22,33]",
			expected:    []any{11.0, 22.0, 33.0}, // JSON parses numbers as float64
		},
		{
			name:        "empty",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    "",
			expected:    nil,
		},
		{
			name:        "just spaces",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    "   ",
			expected:    nil,
		},
		{
			name:        "single element",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    "42",
			expected:    []any{42},
		},
		{
			name:        "single element with leading delimiter",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    ",42",
			expected:    []any{42},
		},
		{
			name:        "single element with trailing delimiter",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    "42,",
			expected:    []any{42},
		},
		{
			name:        "integers with comma delimiter",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    "1,2,3",
			expected:    []any{1, 2, 3},
		},
		{
			name:        "integers with spaces and comma delimiter",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    "   10, 20, 30  ",
			expected:    []any{10, 20, 30},
		},
		{
			name:        "integers with semi-colon delimiter",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    "100;200;300",
			expected:    []any{100, 200, 300},
		},
		{
			name:        "integers with spaces delimiters",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    "1000 2000 3000",
			expected:    []any{1000, 2000, 3000},
		},
		{
			name:        "integers with mixed delimiters",
			itemsSchema: mgcSchemaPkg.NewIntegerSchema(),
			rawValue:    "1000; 2000, 3000 4000",
			expected:    []any{1000, 2000, 3000, 4000},
		},
		{
			name:        "numbers",
			itemsSchema: mgcSchemaPkg.NewNumberSchema(),
			rawValue:    "1,2.3",
			expected:    []any{1.0, 2.3},
		},
		{
			name:        "booleans",
			itemsSchema: mgcSchemaPkg.NewBooleanSchema(),
			rawValue:    "true,false",
			expected:    []any{true, false},
		},
		{
			name:        "simple string",
			itemsSchema: mgcSchemaPkg.NewStringSchema(),
			rawValue:    "apple,banana",
			expected:    []any{"apple", "banana"},
		},
		{
			name:        "simple string with trailing spaces",
			itemsSchema: mgcSchemaPkg.NewStringSchema(),
			rawValue:    "  apple , banana  ",
			expected:    []any{"apple", "banana"},
		},
		{
			name:        "simple string with inner spaces",
			itemsSchema: mgcSchemaPkg.NewStringSchema(),
			rawValue:    "  keep  inner  spaces , x  ",
			expected:    []any{"keep  inner  spaces", "x"},
		},
		{
			name:        "quoted strings",
			itemsSchema: mgcSchemaPkg.NewStringSchema(),
			rawValue:    "  \"keep  \\\"inner\\\"  spaces\" , \"x\"  ",
			expected:    []any{"keep  \"inner\"  spaces", "x"},
		},
	}

	for _, tc := range tests {
		name := tc.name
		if name == "" {
			tc.name = tc.rawValue
		}
		t.Run(name, func(t *testing.T) {
			got, err := parseArrayFlagValueSingle(tc.itemsSchema, tc.rawValue)
			checkError(t, tc.err, err)
			checkArray(t, tc.expected, got)
		})
	}
}

func Test_parseArrayFlagValue(t *testing.T) {
	// most tests were already done before, just ensure it concatenates
	got, err := parseArrayFlagValue(
		mgcSchemaPkg.NewIntegerSchema(),
		[]string{
			"1,2,3",
			"4",
			"",
			",5,6",
		},
	)
	expected := []any{1, 2, 3, 4, 5, 6}

	checkError(t, nil, err)
	checkArray(t, expected, got)
}

func checkObject(t *testing.T, message string, expected, got map[string]any) {
	if len(expected) != len(got) {
		t.Errorf("%s expected object length %d, got %d.\nexpected: %#v\ngot.....: %#v", message, len(expected), len(got), expected, got)
		return
	}

	for k, vExpected := range expected {
		if vGot, ok := got[k]; !ok {
			t.Errorf("%s expected key %q value %#v (%T), got nothing", message, k, vExpected, vExpected)
		} else if !reflect.DeepEqual(vExpected, vGot) {
			if mExpected, expectedMap := vExpected.(map[string]any); expectedMap {
				if mGot, gotMap := vGot.(map[string]any); gotMap {
					checkObject(t, fmt.Sprintf("%s/%s", message, k), mExpected, mGot)
					continue
				}
			}
			// print type as it's usually the case, int x int64 x float64...
			t.Errorf("%s expected key %q value %#v (%T), got %#v (%T)", message, k, vExpected, vExpected, vGot, vGot)
		}
	}
}

func Test_parseObjectFlagValueSingle(t *testing.T) {
	schema := mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
		"int":    mgcSchemaPkg.NewIntegerSchema(),
		"bool":   mgcSchemaPkg.NewBooleanSchema(),
		"number": mgcSchemaPkg.NewNumberSchema(),
		"str":    mgcSchemaPkg.NewStringSchema(),
		"any":    mgcSchemaPkg.NewAnySchema(),
		"obj": mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			"root": mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
				"child": mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
					"int": mgcSchemaPkg.NewIntegerSchema(),
					"str": mgcSchemaPkg.NewStringSchema(),
				}, nil),
			}, nil),
		}, nil),
	}, nil)

	type testCase struct {
		name     string
		rawValue string
		expected map[string]any
		err      error
	}

	tests := []testCase{
		{
			name:     "json",
			rawValue: `{"int": 1, "bool": true, "number": 2.0, "str": "word", "any": [true, false]}`,
			expected: map[string]any{
				"int":    1.0, // JSON parses numbers as float64
				"bool":   true,
				"number": 2.0,
				"str":    "word",
				"any":    []any{true, false},
			},
		},
		{
			name:     "empty",
			rawValue: "",
			expected: nil,
		},
		{
			name:     "just spaces",
			rawValue: "   ",
			expected: nil,
		},
		{
			name:     "single element",
			rawValue: "int=1",
			expected: map[string]any{"int": 1},
		},
		{
			name:     "single element with leading delimiter",
			rawValue: ",int=1",
			expected: map[string]any{"int": 1},
		},
		{
			name:     "single element with trailing delimiter",
			rawValue: "int=1,",
			expected: map[string]any{"int": 1},
		},
		{
			name:     "with comma delimiter",
			rawValue: "int=1,number=2,bool=TRUE,str=word,any=[true,false]",
			expected: map[string]any{
				"int":    1,
				"bool":   true,
				"number": 2.0,
				"str":    "word",
				"any":    []any{true, false},
			},
		},
		{
			name:     "with mixed delimiters and spaces",
			rawValue: "   int = 1; number: 2,bool=TRUE, str=word; any=[true,false]",
			expected: map[string]any{
				"int":    1,
				"bool":   true,
				"number": 2.0,
				"str":    "word",
				"any":    []any{true, false},
			},
		},
		{
			name:     "single nested object with int children",
			rawValue: `obj=root=child=int=123`,
			expected: map[string]any{
				"obj": map[string]any{
					"root": map[string]any{
						"child": map[string]any{
							"int": 123,
						},
					},
				},
			},
		},
		{
			name:     "single nested object with string children",
			rawValue: "obj=root=child=str=hello world",
			expected: map[string]any{
				"obj": map[string]any{
					"root": map[string]any{
						"child": map[string]any{
							"str": "hello world",
						},
					},
				},
			},
		},
		{
			name:     "single nested object with quoted strings",
			rawValue: `"obj"="root"="child"="str"="hello world"`,
			expected: map[string]any{
				"obj": map[string]any{
					"root": map[string]any{
						"child": map[string]any{
							"str": "hello world",
						},
					},
				},
			},
		},
		{
			name:     "single nested object with spaces",
			rawValue: `  obj = "root" = child = str  = "hello world"  `,
			expected: map[string]any{
				"obj": map[string]any{
					"root": map[string]any{
						"child": map[string]any{
							"str": "hello world",
						},
					},
				},
			},
		},
		{
			name:     "multiple nested object with spaces",
			rawValue: `  obj = "root" = child = str  = "hello world";  obj=root=child=int=123  `,
			expected: map[string]any{
				"obj": map[string]any{
					"root": map[string]any{
						"child": map[string]any{
							"int": 123,
							"str": "hello world",
						},
					},
				},
			},
		},
		{
			name:     "multiple nested object with json",
			rawValue: `  obj=root=child=int=123; obj = "root" = child = {"str": "hello world"}  `,
			expected: map[string]any{
				"obj": map[string]any{
					"root": map[string]any{
						"child": map[string]any{
							"int": 123,
							"str": "hello world",
						},
					},
				},
			},
		},
		{
			name:     "wrong object key value",
			rawValue: `  obj = root = 42  `,
			expected: map[string]any{
				"obj": map[string]any{
					"root": 42.0, // if we fail to find a matching property, we just parse it as json and keep going
				},
			},
		},
		{
			name:     "non-existing object key",
			rawValue: `  obj = bug = other = 1  `,
			expected: map[string]any{
				"obj": "bug = other = 1", // same for unknown properties, parse as json and keep going
			},
		},
		{
			name:     "single path object with int children",
			rawValue: "obj=root.child.int=123",
			expected: map[string]any{
				"obj": map[string]any{
					"root": map[string]any{
						"child": map[string]any{
							"int": 123,
						},
					},
				},
			},
		},
		{
			name:     "single quoted path object with int children",
			rawValue: `"obj"."root"."child"."int"=123`, // obj.root or obj=root should be the same
			expected: map[string]any{
				"obj": map[string]any{
					"root": map[string]any{
						"child": map[string]any{
							"int": 123,
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		name := tc.name
		if name == "" {
			tc.name = tc.rawValue
		}
		t.Run(name, func(t *testing.T) {
			got, err := parseObjectFlagValueSingle(schema, tc.rawValue)
			checkError(t, tc.err, err)
			checkObject(t, "", tc.expected, got)
		})
	}
}

func Test_parseObjectFlagValue(t *testing.T) {
	schema := mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
		"int":    mgcSchemaPkg.NewIntegerSchema(),
		"bool":   mgcSchemaPkg.NewBooleanSchema(),
		"number": mgcSchemaPkg.NewNumberSchema(),
		"str":    mgcSchemaPkg.NewStringSchema(),
		"any":    mgcSchemaPkg.NewAnySchema(),
	}, nil)

	// most tests were already done before, just ensure it concatenates
	got, err := parseObjectFlagValue(
		schema,
		[]string{
			"int=1, number=2",
			"bool=TRUE,str=word",
			"",
			",any=[true,false]",
		},
	)
	expected := map[string]any{
		"int":    1,
		"bool":   true,
		"number": 2.0,
		"str":    "word",
		"any":    []any{true, false},
	}

	checkError(t, nil, err)
	checkObject(t, "", expected, got)
}

func Test_help(t *testing.T) {
	types := []string{"array", "boolean", "integer", "string", "object", ""}
	for _, typ := range types {
		t.Run(typ, func(t *testing.T) {
			fv := newSchemaFlagValue(SchemaFlagValueDesc{
				Schema:   &core.Schema{Type: typ},
				PropName: "name",
				FlagName: "name",
			})

			_ = fv.Set("help")
			_, err := fv.Parse()
			checkError(t, ErrWantHelp, err)
		})
	}
}
