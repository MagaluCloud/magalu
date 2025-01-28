package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"slices"

	"maps"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
)

const testPathRef = "/path/to/ref"

func copySchemaRef(input *SchemaRef) *SchemaRef {
	if input == nil {
		return nil
	}
	output := *input
	output.Value = (*openapi3.Schema)(copySchema((*Schema)(input.Value)))
	return &output
}

func copySchemaRefSlice(input SchemaRefs) SchemaRefs {
	if input == nil {
		return nil
	}
	output := make(SchemaRefs, len(input))
	for i, e := range input {
		output[i] = copySchemaRef(e)
	}
	return output
}

func copySchemaRefMap(input map[string]*SchemaRef) map[string]*SchemaRef {
	if input == nil {
		return nil
	}
	output := make(map[string]*SchemaRef, len(input))
	for k, e := range input {
		output[k] = copySchemaRef(e)
	}
	return output
}

func copySchema(input *Schema) *Schema {
	if input == nil {
		return nil
	}
	output := *input
	output.Extensions = maps.Clone(input.Extensions)
	output.OneOf = copySchemaRefSlice(input.OneOf)
	output.AnyOf = copySchemaRefSlice(input.AnyOf)
	output.AllOf = copySchemaRefSlice(input.AllOf)
	output.Enum = slices.Clone(input.Enum)
	output.Items = copySchemaRef(input.Items)
	output.Required = slices.Clone(input.Required)
	output.Properties = copySchemaRefMap(input.Properties)
	output.AdditionalProperties.Schema = copySchemaRef(input.AdditionalProperties.Schema)
	return &output
}

func checkPointers_Equal[T any](t *testing.T, prefix string, expected, got *T) {
	if expected != got {
		t.Errorf("%s. Diverging %T.\nExpected:\n%p\nGot:\n%p\n", prefix, expected, expected, got)
	}
}

func checkNoError(t *testing.T, prefix string, err error) {
	if err != nil {
		t.Errorf("%s. Unexpected error: %s\n", prefix, err)
	}
}

func check_Must_Error(t *testing.T, prefix string, err error) {
	if err == nil {
		t.Errorf("%s. Did expect error, got none!\n", prefix)
	}
}

func checkMustChange[T any](t *testing.T, prefix string, orig, got *T, changed bool) {
	if !changed {
		t.Errorf("%s. Expected changes, but it didn't\n", prefix)
	} else if orig == got {
		t.Errorf("%s. Unexpected same %T.\nOriginal:\n%p\nExpected Changed:\n%p\n", prefix, orig, orig, got)
	}
}

func check_No_Changes[T any](t *testing.T, prefix string, orig, got *T, changed bool) {
	if changed {
		t.Errorf("%s. Expected NO changes, but it did\n", prefix)
	} else if orig != got {
		t.Errorf("%s. Expected same %T.\nOriginal:\n%p\nChanged:\n%p\n", prefix, orig, orig, got)
	}
}

func checkSchemasEqual(t *testing.T, prefix string, expected, got *Schema) {
	if err := CompareJsonSchemas(expected, got); err != nil {
		expectedJson, _ := json.MarshalIndent(expected, "", "  ")
		gotJson, _ := json.MarshalIndent(got, "", "  ")
		t.Errorf("%s. Diverging schemas.\nExpected:\n%s\n%#v\nGot:\n%s\n%#v\nError: %s\n", prefix, expectedJson, expected, gotJson, got, err.Error())
	}
}

func checkSchemaRefsEqual(t *testing.T, prefix string, expected, got *SchemaRef) {
	if expected == nil && got == nil {
		return
	} else if expected == nil {
		t.Errorf("%s. Did not expect a schemaRef, but got one. Got:\n%#v", prefix, got)
	} else if got == nil {
		t.Errorf("%s. Did expect a schemaRef, but got NONE. Expected:\n%#v", prefix, expected)
	} else if expected.Ref != got.Ref {
		t.Errorf("%s. Diverging refs.\nExpected:\n%q\nGot:\n%q\n", prefix, expected.Ref, got.Ref)
	} else {
		checkSchemasEqual(t, prefix, (*Schema)(expected.Value), (*Schema)(got.Value))
	}
}

var nullSchema = &Schema{Type: "null"}

var templateComplexStringSchema = func() *Schema {
	maxLen := uint64(255)
	return &Schema{
		Type:        "string",
		Format:      "uri",
		Description: "some description",
		Enum:        []any{"http://localhost", "https://localhost", "http://server.com"},
		Default:     "http://localhost",
		Example:     "http://localhost/1234",
		MinLength:   1,
		MaxLength:   &maxLen,
		Pattern:     "https?://.*",
	}
}()

// The same value as templateComplexStringSchema, but composed of 2 schemas merged with AllOf:
var templateComplexStringSchemaToBeSimplified = func() *Schema {
	maxLen := uint64(255)
	return NewAllOfSchema(
		&Schema{
			Type:        "string",
			Format:      "uri",
			Description: "some description",
		},
		&Schema{
			Enum:      []any{"http://localhost", "https://localhost", "http://server.com"},
			Default:   "http://localhost",
			Example:   "http://localhost/1234",
			MinLength: 1,
			MaxLength: &maxLen,
			Pattern:   "https?://.*",
		},
	)
}()

var templateComplexStringSchemaNullable = func() (schema *Schema) {
	schema = copySchema(templateComplexStringSchema)
	schema.Nullable = true
	return schema
}()

// The same value as templateComplexStringSchemaNullable, but composed of 2 schemas merged with OneOf:
var templateComplexStringSchemaNullableToBeSimplified = NewOneOfSchema(
	templateComplexStringSchema,
	nullSchema,
)

type testCaseSchemaRef struct {
	orig     *SchemaRef
	expected *SchemaRef
}

func (tc testCaseSchemaRef) isChanged() bool {
	return tc.orig != tc.expected
}

type testCaseSchema struct {
	orig     *Schema
	expected *Schema
}

func (tc testCaseSchema) isError() bool {
	return tc.orig != nil && tc.expected == nil
}

func (tc testCaseSchema) isChanged() bool {
	return !tc.isError() && tc.orig != tc.expected
}

func runSubTests[T any](t *testing.T, prefix string, run func(t *testing.T, prefix string, tc T), tests map[string]T) {
	keys := make([]string, 0, len(tests))
	for name := range tests {
		keys = append(keys, name)
	}

	slices.Sort(keys)

	for _, name := range keys {
		tc := tests[name]
		t.Run(name, func(t *testing.T) {
			run(t, prefix+": "+name, tc)
		})
	}
}

// -- SimplifySchemaRefCOW

func checkSimplifySchemaRefSimplifies(t *testing.T, prefix string, tc testCaseSchemaRef) {
	cow := NewCOWSchemaRef(tc.orig)
	err := SimplifySchemaRefCOW(cow)
	checkNoError(t, prefix, err)

	got, changed := cow.Release()
	if tc.isChanged() {
		checkMustChange(t, prefix, tc.orig, got, changed)
		checkSchemaRefsEqual(t, prefix, tc.expected, got)
	} else {
		check_No_Changes(t, prefix, tc.orig, got, changed)
	}
}

func TestSimplifySchemaRefCOW(t *testing.T) {
	schema := openapi3.NewStringSchema()
	runSubTests(t, "SimplifySchemaRefCOW", checkSimplifySchemaRefSimplifies, map[string]testCaseSchemaRef{
		"NilValue": {
			orig:     nil,
			expected: nil,
		},
		"UnsetRef": {
			orig:     &SchemaRef{Ref: testPathRef},
			expected: &SchemaRef{},
		},
		"SameValueWithoutRef": {
			orig:     &SchemaRef{Ref: testPathRef, Value: schema},
			expected: &SchemaRef{Value: schema},
		},
		"Simplified": {
			orig:     &SchemaRef{Value: (*openapi3.Schema)(copySchema(templateComplexStringSchemaToBeSimplified))},
			expected: &SchemaRef{Value: (*openapi3.Schema)(copySchema(templateComplexStringSchema))},
		},
	})
}

// -- SimplifySchema

func checkSimplifySchemaSimplifies(t *testing.T, prefix string, tc testCaseSchema) {
	got, err := SimplifySchema(tc.orig)
	if tc.isError() {
		check_Must_Error(t, prefix, err)
		return
	} else {
		checkNoError(t, prefix, err)
	}

	if tc.isChanged() {
		checkSchemasEqual(t, prefix, tc.expected, got)
	} else {
		checkPointers_Equal(t, prefix, tc.expected, got)
	}
}

func TestSimplifySchema(t *testing.T) {
	schema := NewStringSchema()

	// TODO: Uncomment recursive schema tests when 'CompareJsonSchemas' supports recursive comparisons...

	// recursiveUnchangedRootSchema := &Schema{Items: &SchemaRef{}}
	// recursiveUnchangedRootSchema.Items.Value = (*openapi3.Schema)(recursiveUnchangedRootSchema)

	// recursiveUnchangedRootSimplifiedSchema := &Schema{Items: &SchemaRef{Ref: "/"}}
	// recursiveUnchangedRootSimplifiedSchema.Items.Value = (*openapi3.Schema)(recursiveUnchangedRootSimplifiedSchema)

	// recursiveUnchangedChildSchema := &Schema{Items: &SchemaRef{
	// 	Value: &openapi3.Schema{
	// 		OneOf: openapi3.SchemaRefs{
	// 			&SchemaRef{Value: &openapi3.Schema{Type: "string"}},
	// 			&SchemaRef{Value: &openapi3.Schema{Type: "integer"}},
	// 			&SchemaRef{},
	// 		},
	// 	},
	// }}
	// recursiveUnchangedChildSchema.Items.Value.OneOf[2].Value = recursiveUnchangedChildSchema.Items.Value

	// recursiveUnchangedChildSimplifiedSchema := &Schema{Items: &SchemaRef{
	// 	Value: &openapi3.Schema{
	// 		OneOf: openapi3.SchemaRefs{
	// 			&SchemaRef{Value: &openapi3.Schema{Type: "string"}},
	// 			&SchemaRef{Value: &openapi3.Schema{Type: "integer"}},
	// 			&SchemaRef{Ref: "/items"},
	// 		},
	// 	},
	// }}
	// recursiveUnchangedChildSimplifiedSchema.Items.Value.OneOf[2].Value = recursiveUnchangedChildSimplifiedSchema.Items.Value

	// recursiveToBeSimplifiedSchema := &Schema{OneOf: openapi3.SchemaRefs{
	// 	&SchemaRef{},
	// 	&SchemaRef{Value: (*openapi3.Schema)(copySchema(templateComplexStringSchemaToBeSimplified))}},
	// }
	// recursiveToBeSimplifiedSchema.OneOf[0].Value = (*openapi3.Schema)(recursiveToBeSimplifiedSchema)

	// recursiveSimplifiedSchema := &Schema{OneOf: openapi3.SchemaRefs{
	// 	&SchemaRef{Ref: "/"},
	// 	&SchemaRef{Value: (*openapi3.Schema)(copySchema(templateComplexStringSchema))}},
	// 	Type: "string",
	// }
	// recursiveSimplifiedSchema.OneOf[0].Value = (*openapi3.Schema)(recursiveSimplifiedSchema)

	runSubTests(t, "SimplifySchema", checkSimplifySchemaSimplifies, map[string]testCaseSchema{
		"NilValue": {
			orig:     nil,
			expected: nil,
		},
		"Unmodified": {
			orig:     schema,
			expected: schema,
		},
		"Simplified": {
			orig:     copySchema(templateComplexStringSchemaToBeSimplified),
			expected: copySchema(templateComplexStringSchema),
		},
		// "RecursiveUnchangedRoot": {
		// 	orig:     recursiveUnchangedRootSchema,
		// 	expected: recursiveUnchangedRootSimplifiedSchema,
		// },
		// "RecursiveUnchangedChild": {
		// 	orig:     recursiveUnchangedChildSchema,
		// 	expected: recursiveUnchangedChildSimplifiedSchema,
		// },
		// "RecursiveSimplified": {
		// 	orig:     recursiveToBeSimplifiedSchema,
		// 	expected: recursiveSimplifiedSchema,
		// },
	})
}

// -- SimplifySchemaCOW

func checkSimplifySchemaCOWSimplifies(t *testing.T, prefix string, tc testCaseSchema) {
	schemaCow := NewCOWSchema(tc.orig)
	err := SimplifySchemaCOW(schemaCow)
	if tc.isError() {
		check_Must_Error(t, prefix, err)
	} else {
		checkNoError(t, prefix, err)
	}

	got, changed := schemaCow.Release()
	if tc.isChanged() {
		checkMustChange(t, prefix, tc.orig, got, changed)
		checkSchemasEqual(t, prefix, tc.expected, got)
	} else {
		check_No_Changes(t, prefix, tc.orig, got, changed)
	}
}

func addSimplifySchemaCOWTestCases_Enum(m map[string]testCaseSchema) {
	tests := map[string]struct {
		enum         []any
		origType     string
		expectedType string
	}{
		"set-string": {
			enum:         []any{"a", "b"},
			expectedType: "string",
		},
		"set-integer": {
			enum:         []any{1, 2, 3},
			expectedType: "integer",
		},
		"set-boolean": {
			enum:         []any{true, false},
			expectedType: "boolean",
		},
		"set-number": {
			enum:         []any{1.2, 3.4},
			expectedType: "number",
		},
		"mixed": {
			enum:         []any{"a", 1, true},
			expectedType: "",
		},
		"unchanged-number": {
			enum:         []any{1.2, 3.4},
			origType:     "number",
			expectedType: "number",
		},
	}

	for name, e := range tests {
		tc := testCaseSchema{orig: &Schema{Enum: e.enum, Type: e.origType}}
		if e.expectedType == e.origType {
			tc.expected = tc.orig
		} else {
			expected := copySchema(tc.orig)
			expected.Type = e.expectedType
			tc.expected = expected
		}
		m["Enum/"+name] = tc
	}
}

func TestSimplifySchemaCOW(t *testing.T) {
	simpleSchema := NewStringSchema()
	tests := map[string]testCaseSchema{
		"NilValue": {
			orig:     nil,
			expected: nil,
		},
		"Unmodified": {
			orig:     simpleSchema,
			expected: simpleSchema,
		},
		"Not": {
			orig:     &Schema{Not: &SchemaRef{Value: openapi3.NewStringSchema()}},
			expected: nil, // should error
		},
		"Items": {
			orig:     NewArraySchema(copySchema(templateComplexStringSchemaToBeSimplified)),
			expected: NewArraySchema(templateComplexStringSchema),
		},
	}

	addSimplifySchemaCOWTestCases_Enum(tests)
	addSimplifySchemaCOWTestCases_OneOf(tests)
	addSimplifySchemaCOWTestCases_AnyOf(tests)
	addSimplifySchemaCOWTestCases_AllOf(tests)
	addSimplifySchemaCOWTestCases_Properties(tests)
	addSimplifySchemaCOWTestCases_AdditionalProperties(tests)

	runSubTests(t, "SimplifySchemaCOW", checkSimplifySchemaCOWSimplifies, tests)
}

type testCaseSchemaChildren struct {
	origChildren     []*Schema
	expectedChildren []*Schema
	origType         string
	expectedType     string
}

func addSimplifySchemaCOWTestCases_Children(prefix string, creator func(...*Schema) *Schema, tests map[string]testCaseSchema, children map[string]testCaseSchemaChildren) {
	for name, e := range children {
		orig := creator(e.origChildren...)
		orig.Type = e.origType
		tc := testCaseSchema{orig: orig}
		if e.expectedType == e.origType && utils.IsSameValueOrPointer(e.origChildren, e.expectedChildren) {
			tc.expected = tc.orig
		} else {
			expected := creator(e.expectedChildren...)
			expected.Type = e.expectedType
			tc.expected = expected
		}
		k := prefix + "/" + name
		if e, ok := tests[k]; ok {
			panic(fmt.Sprint("duplicated test name: ", k, ", has: ", e, ", new: ", tc))
		}
		tests[k] = tc
	}
}

type testCaseSchemaChildrenMerged struct {
	origChildren []*Schema
	expected     *Schema
	origType     string
}

func addSimplifySchemaCOWTestCases_ChildrenMerged(prefix string, creator func(...*Schema) *Schema, tests map[string]testCaseSchema, children map[string]testCaseSchemaChildrenMerged) {
	for name, e := range children {
		orig := creator(e.origChildren...)
		orig.Type = e.origType
		tc := testCaseSchema{orig: orig, expected: e.expected}
		k := prefix + "/" + name
		if e, ok := tests[k]; ok {
			panic(fmt.Sprint("duplicated test name: ", k, ", has: ", e, ", new: ", tc))
		}
		tests[k] = tc
	}
}

func addSimplifySchemaCOWTestCases_OneOf(tests map[string]testCaseSchema) {
	simpleMixedChildren := []*Schema{NewStringSchema(), NewIntegerSchema()}
	simpleSameChildren := []*Schema{NewStringSchema(), NewStringSchema()}

	prefix := "OneOf"
	creator := NewOneOfSchema

	addSimplifySchemaCOWTestCases_Children(prefix, creator, tests, map[string]testCaseSchemaChildren{
		"Unmodified": {
			origChildren:     simpleMixedChildren,
			expectedChildren: simpleMixedChildren,
		},
		"SimplifiesChildren": {
			origChildren:     []*Schema{copySchema(templateComplexStringSchemaToBeSimplified), NewIntegerSchema()},
			expectedChildren: []*Schema{templateComplexStringSchema, NewIntegerSchema()},
		},
		"SetTypeFromChildren": {
			origChildren:     simpleSameChildren,
			expectedChildren: simpleSameChildren,
			expectedType:     "string",
		},
	})

	addSimplifySchemaCOWTestCases_ChildrenMerged(prefix, creator, tests, map[string]testCaseSchemaChildrenMerged{
		"PromoteSingleChild": {
			origChildren: []*Schema{copySchema(templateComplexStringSchema)},
			expected:     templateComplexStringSchema,
		},
		"PromoteChildNullable": {
			origChildren: []*Schema{copySchema(templateComplexStringSchema), copySchema(nullSchema)},
			expected:     templateComplexStringSchemaNullable,
		},
	})
}

func addSimplifySchemaCOWTestCases_AnyOf(tests map[string]testCaseSchema) {
	simpleMixedChildren := []*Schema{NewStringSchema(), NewIntegerSchema()}
	simpleSameChildren := []*Schema{NewStringSchema(), NewStringSchema()}

	prefix := "AnyOf"
	creator := NewAnyOfSchema

	addSimplifySchemaCOWTestCases_Children(prefix, creator, tests, map[string]testCaseSchemaChildren{
		"Unmodified": {
			origChildren:     simpleMixedChildren,
			expectedChildren: simpleMixedChildren,
		},
		"SimplifiesChildren": {
			origChildren:     []*Schema{copySchema(templateComplexStringSchemaToBeSimplified), NewIntegerSchema()},
			expectedChildren: []*Schema{templateComplexStringSchema, NewIntegerSchema()},
		},
		"SetTypeFromChildren": {
			origChildren:     simpleSameChildren,
			expectedChildren: simpleSameChildren,
			expectedType:     "string",
		},
	})

	addSimplifySchemaCOWTestCases_ChildrenMerged(prefix, creator, tests, map[string]testCaseSchemaChildrenMerged{
		"PromoteSingleChild": {
			origChildren: []*Schema{copySchema(templateComplexStringSchema)},
			expected:     templateComplexStringSchema,
		},
		"PromoteChildNullable": {
			origChildren: []*Schema{copySchema(templateComplexStringSchema), copySchema(nullSchema)},
			expected:     templateComplexStringSchemaNullable,
		},
	})
}

func addSimplifySchemaCOWTestCases_AllOf(tests map[string]testCaseSchema) {
	simpleMixedChildren := []*Schema{NewStringSchema(), NewIntegerSchema()}
	simpleStringChildren := []*Schema{NewStringSchema(), NewStringSchema()}

	mergeIntoParentExpected := copySchema(templateComplexStringSchema)
	mergeIntoParentExpected.MinItems = 123 // MinItems is not string, but we just want to check if it was copied as well

	prefix := "AllOf"
	creator := NewAllOfSchema

	addSimplifySchemaCOWTestCases_ChildrenMerged(prefix, creator, tests, map[string]testCaseSchemaChildrenMerged{
		"ErrorDifferentChildren": {
			origChildren: simpleMixedChildren,
			expected:     nil, // should error
		},
		"ErrorDifferentChildrenAndType": {
			origChildren: simpleStringChildren,
			origType:     "number",
			expected:     nil, // should error
		},
		"MergeIntoParent": {
			origChildren: []*Schema{
				{MinItems: mergeIntoParentExpected.MinItems},
				copySchema(templateComplexStringSchema),
			},
			expected: mergeIntoParentExpected,
		},
		"SimplifyAndMerge": {
			origChildren: []*Schema{copySchema(templateComplexStringSchemaNullableToBeSimplified)},
			expected:     templateComplexStringSchemaNullable,
		},
	})
}

func addSimplifySchemaCOWTestCases_Properties(tests map[string]testCaseSchema) {
	simpleObject := NewObjectSchema(map[string]*Schema{"simpleProp": NewStringSchema()}, []string{"simpleProp"})

	prefix := "Properties"

	tests[prefix+"/Unmodified"] = testCaseSchema{
		orig:     simpleObject,
		expected: simpleObject,
	}
	tests[prefix+"/SimplifiesChildren"] = testCaseSchema{
		orig: NewObjectSchema(
			map[string]*Schema{
				"simpleProp":   NewStringSchema(),
				"optionalProp": NewIntegerSchema(),
				"complexProp":  copySchema(templateComplexStringSchemaNullableToBeSimplified),
			},
			[]string{"simpleProp", "complexProp"},
		),
		expected: NewObjectSchema(
			map[string]*Schema{
				"simpleProp":   NewStringSchema(),
				"optionalProp": NewIntegerSchema(),
				"complexProp":  templateComplexStringSchemaNullable,
			},
			[]string{"simpleProp", "complexProp"},
		),
	}
}

func addSimplifySchemaCOWTestCases_AdditionalProperties(tests map[string]testCaseSchema) {
	simpleAdditionalProperties := NewObjectSchema(nil, nil)

	has := true
	simpleAdditionalProperties.AdditionalProperties.Has = &has
	simpleAdditionalProperties.AdditionalProperties.Schema = &SchemaRef{Value: (*openapi3.Schema)(NewIntegerSchema())}

	origComplexAdditionalProperties := NewObjectSchema(nil, nil)
	origComplexAdditionalProperties.AdditionalProperties.Has = &has
	origComplexAdditionalProperties.AdditionalProperties.Schema = &SchemaRef{Value: (*openapi3.Schema)(copySchema(templateComplexStringSchemaToBeSimplified))}

	expectedComplexAdditionalProperties := NewObjectSchema(nil, nil)
	expectedComplexAdditionalProperties.AdditionalProperties.Has = &has
	expectedComplexAdditionalProperties.AdditionalProperties.Schema = &SchemaRef{Value: (*openapi3.Schema)(templateComplexStringSchema)}

	prefix := "AdditionalProperties"

	tests[prefix+"/Unmodified"] = testCaseSchema{
		orig:     simpleAdditionalProperties,
		expected: simpleAdditionalProperties,
	}
	tests[prefix+"/SimplifiesChildren"] = testCaseSchema{
		orig:     origComplexAdditionalProperties,
		expected: expectedComplexAdditionalProperties,
	}
}

// -- mergeEnum

type testCaseMergeEnum struct {
	orig     []any
	target   []any
	expected []any
}

func (tc testCaseMergeEnum) isError() bool {
	return len(tc.orig) > 0 && len(tc.expected) == 0
}

func (tc testCaseMergeEnum) isChanged() bool {
	return !tc.isError() && !utils.IsSameValueOrPointer(tc.orig, tc.expected)
}

func checkMergeEnum(t *testing.T, prefix string, tc testCaseMergeEnum) {
	orig := &Schema{Enum: tc.orig}
	cow := NewCOWSchema(orig)
	err := mergeEnum(cow, tc.target)
	if tc.isError() {
		check_Must_Error(t, prefix, err)
	} else {
		checkNoError(t, prefix, err)
	}

	got, changed := cow.Release()
	if tc.isChanged() {
		checkMustChange(t, prefix, orig, got, changed)
		if !reflect.DeepEqual(got.Enum, tc.expected) {
			t.Errorf("%s. Diverging enums.\nExpected:\n%#v\nGot:\n%#v\n", prefix, tc.expected, got.Enum)
		}
	} else {
		check_No_Changes(t, prefix, orig, got, changed)
	}
}

func Test_mergeEnum(t *testing.T) {
	enum := []any{1, "a", true}

	runSubTests(t, "mergeEnum", checkMergeEnum, map[string]testCaseMergeEnum{
		"Unmodified/Nil": {
			orig:     enum,
			target:   nil,
			expected: enum,
		},
		"Unmodified/Empty": {
			orig:     enum,
			target:   []any{},
			expected: enum,
		},
		"Unmodified/Same": {
			orig:     enum,
			target:   slices.Clone(enum),
			expected: enum,
		},
		"Unmodified/Contained": {
			orig:     enum,
			target:   []any{1},
			expected: enum,
		},
		"SetAsTarget/Nil": {
			orig:     nil,
			target:   enum,
			expected: enum,
		},
		"SetAsTarget/Empty": {
			orig:     []any{},
			target:   enum,
			expected: enum,
		},
		"Diverging": {
			orig:     slices.Clone(enum),
			target:   []any{1, "a", true, false, "b", 0}, // extra values are forbidden
			expected: nil,
		},
	})
}

// -- mergeRequired

type testCaseMergeRequired struct {
	orig     []string
	target   []string
	expected []string
}

func (tc testCaseMergeRequired) isChanged() bool {
	return !utils.IsSameValueOrPointer(tc.orig, tc.expected)
}

func checkMergeRequired(t *testing.T, prefix string, tc testCaseMergeRequired) {
	orig := &Schema{Required: tc.orig}
	cow := NewCOWSchema(orig)
	mergeRequired(cow, tc.target)

	got, changed := cow.Release()
	if tc.isChanged() {
		checkMustChange(t, prefix, orig, got, changed)
		if !reflect.DeepEqual(got.Required, tc.expected) {
			t.Errorf("%s. Diverging required.\nExpected:\n%#v\nGot:\n%#v\n", prefix, tc.expected, got.Required)
		}
	} else {
		check_No_Changes(t, prefix, orig, got, changed)
	}
}

func Test_mergeRequired(t *testing.T) {
	required := []string{"a", "b", "xpto"}

	runSubTests(t, "mergeRequired", checkMergeRequired, map[string]testCaseMergeRequired{
		"Unmodified/Nil": {
			orig:     required,
			target:   nil,
			expected: required,
		},
		"Unmodified/Empty": {
			orig:     required,
			target:   []string{},
			expected: required,
		},
		"Unmodified/Same": {
			orig:     required,
			target:   slices.Clone(required),
			expected: required,
		},
		"Unmodified/Contained": {
			orig:     required,
			target:   []string{"xpto"},
			expected: required,
		},
		"SetAsTarget/Nil": {
			orig:     nil,
			target:   required,
			expected: required,
		},
		"SetAsTarget/Empty": {
			orig:     []string{},
			target:   required,
			expected: required,
		},
		"Merge": {
			orig:     slices.Clone(required),
			target:   []string{"banana", "apple"},
			expected: []string{"a", "b", "xpto", "banana", "apple"},
		},
	})
}

// -- mergeProperties

type testCaseMergeProperties struct {
	orig     openapi3.Schemas
	target   openapi3.Schemas
	expected openapi3.Schemas
}

func (tc testCaseMergeProperties) isError() bool {
	return len(tc.orig) > 0 && len(tc.expected) == 0
}

func (tc testCaseMergeProperties) isChanged() bool {
	return !tc.isError() && !utils.IsSameValueOrPointer(tc.orig, tc.expected)
}

func checkMergeProperties(t *testing.T, prefix string, tc testCaseMergeProperties) {
	orig := &Schema{Properties: tc.orig}
	cow := NewCOWSchema(orig)
	err := mergeProperties(cow, tc.target)
	if tc.isError() {
		check_Must_Error(t, prefix, err)
	} else {
		checkNoError(t, prefix, err)
	}

	got, changed := cow.Release()
	if tc.isChanged() {
		checkMustChange(t, prefix, orig, got, changed)
		if !reflect.DeepEqual(got.Properties, tc.expected) {
			t.Errorf("%s. Diverging properties.\nExpected:\n%#v\nGot:\n%#v\n", prefix, tc.expected, got.Properties)
		}
	} else {
		check_No_Changes(t, prefix, orig, got, changed)
	}
}

func Test_mergeProperties(t *testing.T) {
	properties := openapi3.Schemas{
		"stringProp":  {Value: openapi3.NewStringSchema()},
		"integerProp": {Value: openapi3.NewIntegerSchema()},
	}

	runSubTests(t, "mergeProperties", checkMergeProperties, map[string]testCaseMergeProperties{
		"Unmodified/Nil": {
			orig:     properties,
			target:   nil,
			expected: properties,
		},
		"Unmodified/Empty": {
			orig:     properties,
			target:   openapi3.Schemas{},
			expected: properties,
		},
		"Unmodified/Same": {
			orig:     properties,
			target:   maps.Clone(properties),
			expected: properties,
		},
		"Unmodified/Contained": {
			orig:     properties,
			target:   openapi3.Schemas{"stringProp": {Value: openapi3.NewStringSchema()}},
			expected: properties,
		},
		"SetAsTarget/Nil": {
			orig:     nil,
			target:   properties,
			expected: properties,
		},
		"SetAsTarget/Empty": {
			orig:     openapi3.Schemas{},
			target:   properties,
			expected: properties,
		},
		"Merge": {
			orig:   maps.Clone(properties),
			target: openapi3.Schemas{"newProp": {Value: (*openapi3.Schema)(templateComplexStringSchema)}},
			expected: openapi3.Schemas{
				"stringProp":  {Value: openapi3.NewStringSchema()},
				"integerProp": {Value: openapi3.NewIntegerSchema()},
				"newProp":     {Value: (*openapi3.Schema)(templateComplexStringSchema)},
			},
		},
		"Diverging": {
			orig:     maps.Clone(properties),
			target:   openapi3.Schemas{"stringProp": {Value: openapi3.NewIntegerSchema()}}, // extra values are forbidden
			expected: nil,
		},
	})
}

// -- mergeProperties

type testCaseMergeExtensions struct {
	orig     map[string]any
	target   map[string]any
	expected map[string]any
}

func (tc testCaseMergeExtensions) isError() bool {
	return len(tc.orig) > 0 && len(tc.expected) == 0
}

func (tc testCaseMergeExtensions) isChanged() bool {
	return !tc.isError() && !utils.IsSameValueOrPointer(tc.orig, tc.expected)
}

func checkMergeExtensions(t *testing.T, prefix string, tc testCaseMergeExtensions) {
	orig := &Schema{Extensions: tc.orig}
	cow := NewCOWSchema(orig)
	err := mergeExtensions(cow, tc.target)
	if tc.isError() {
		check_Must_Error(t, prefix, err)
	} else {
		checkNoError(t, prefix, err)
	}

	got, changed := cow.Release()
	if tc.isChanged() {
		checkMustChange(t, prefix, orig, got, changed)
		if !reflect.DeepEqual(got.Extensions, tc.expected) {
			t.Errorf("%s. Diverging extensions.\nExpected:\n%#v\nGot:\n%#v\n", prefix, tc.expected, got.Extensions)
		}
	} else {
		check_No_Changes(t, prefix, orig, got, changed)
	}
}

func Test_mergeExtensions(t *testing.T) {
	extensions := map[string]any{
		"stringProp":  "abc",
		"integerProp": 123,
	}

	runSubTests(t, "mergeExtensions", checkMergeExtensions, map[string]testCaseMergeExtensions{
		"Unmodified/Nil": {
			orig:     extensions,
			target:   nil,
			expected: extensions,
		},
		"Unmodified/Empty": {
			orig:     extensions,
			target:   map[string]any{},
			expected: extensions,
		},
		"Unmodified/Same": {
			orig:     extensions,
			target:   maps.Clone(extensions),
			expected: extensions,
		},
		"Unmodified/Contained": {
			orig:     extensions,
			target:   map[string]any{"stringProp": "abc"},
			expected: extensions,
		},
		"SetAsTarget/Nil": {
			orig:     nil,
			target:   extensions,
			expected: extensions,
		},
		"SetAsTarget/Empty": {
			orig:     map[string]any{},
			target:   extensions,
			expected: extensions,
		},
		"Merge": {
			orig:   maps.Clone(extensions),
			target: map[string]any{"newProp": true},
			expected: map[string]any{
				"stringProp":  "abc",
				"integerProp": 123,
				"newProp":     true,
			},
		},
		"Diverging": {
			orig:     maps.Clone(extensions),
			target:   map[string]any{"stringProp": 123}, // extra values are forbidden
			expected: nil,
		},
	})
}

// -- mergeComparable

func Test_mergeComparable_EmptyTarget(t *testing.T) {
	var err error

	err = mergeComparable(
		func() int {
			t.Error("mergeComparable/EmptyTarget/Scalar: shouldn't get")
			return 1
		},
		func(v int) bool {
			t.Error("mergeComparable/EmptyTarget/Scalar: shouldn't set")
			return false
		},
		0,
	)
	checkNoError(t, "mergeComparable/EmptyTarget/Scalar", err)

	err = mergeComparable(
		func() *int {
			t.Error("mergeComparable/EmptyTarget/Pointer: shouldn't get")
			i := 1
			return &i
		},
		func(v *int) bool {
			t.Error("mergeComparable/EmptyTarget/Pointer: shouldn't set")
			return false
		},
		nil,
	)
	checkNoError(t, "mergeComparable/EmptyTarget/Pointer", err)
}

func Test_mergeComparable_SetIfEmpty(t *testing.T) {
	var err error

	refValue := 123
	setScalar := -1
	err = mergeComparable(
		func() int {
			return 0
		},
		func(v int) bool {
			setScalar = v
			return true
		},
		refValue,
	)
	checkNoError(t, "mergeComparable/SetIfEmpty/Scalar", err)
	if setScalar != refValue {
		t.Errorf("mergeComparable/SetIfEmpty/Scalar: should set to %d, got %d", refValue, setScalar)
	}

	refPointer := 42
	var setPointer *int
	err = mergeComparable(
		func() *int {
			return nil
		},
		func(v *int) bool {
			setPointer = v
			return true
		},
		&refPointer,
	)
	checkNoError(t, "mergeComparable/SetIfEmpty/Pointer", err)
	if setPointer != &refPointer {
		t.Errorf("mergeComparable/SetIfEmpty/Pointer: should set to %p, got %p", &refPointer, setPointer)
	}
}

func Test_mergeComparable_Equal(t *testing.T) {
	var err error

	refValue := 123
	setScalar := -1
	err = mergeComparable(
		func() int {
			return refValue
		},
		func(v int) bool {
			t.Error("mergeComparable/Equal/Scalar: shouldn't set")
			setScalar = v
			return true
		},
		refValue,
	)
	checkNoError(t, "mergeComparable/Equal/Scalar", err)
	if setScalar == refValue {
		t.Errorf("mergeComparable/Equal/Scalar: should NOT set to %d", refValue)
	}

	refPointer := 42
	var setPointer *int
	err = mergeComparable(
		func() *int {
			return &refPointer
		},
		func(v *int) bool {
			t.Error("mergeComparable/Equal/Pointer: shouldn't set")
			setPointer = v
			return true
		},
		&refPointer,
	)
	checkNoError(t, "mergeComparable/Equal/Pointer", err)
	if setPointer == &refPointer {
		t.Errorf("mergeComparable/Equal/Pointer: should NOT set to %p", &refPointer)
	}
}

func Test_mergeComparable_Diverging(t *testing.T) {
	var err error

	refValue := 123
	setScalar := -1
	err = mergeComparable(
		func() int {
			return 2
		},
		func(v int) bool {
			t.Error("mergeComparable/Diverging/Scalar: shouldn't set")
			setScalar = v
			return true
		},
		refValue,
	)
	check_Must_Error(t, "mergeComparable/Diverging/Scalar", err)
	if setScalar == refValue {
		t.Errorf("mergeComparable/Diverging/Scalar: should NOT set to %d", refValue)
	}

	refPointer := 42
	var setPointer *int
	err = mergeComparable(
		func() *int {
			return new(int)
		},
		func(v *int) bool {
			t.Error("mergeComparable/Diverging/Pointer: shouldn't set")
			setPointer = v
			return true
		},
		&refPointer,
	)
	check_Must_Error(t, "mergeComparable/Diverging/Pointer", err)
	if setPointer == &refPointer {
		t.Errorf("mergeComparable/Diverging/Pointer: should NOT set to %p", &refPointer)
	}
}
