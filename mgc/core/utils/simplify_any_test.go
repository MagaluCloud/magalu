package utils

import (
	"reflect"
	"testing"

	"golang.org/x/exp/slices"
)

type testCase struct {
	Input, Output interface{}
}

type simpleStruct struct {
	StrField  string
	IntField  int `json:"int_field_with_json_name"`
	BoolField bool
}

type complexStruct struct {
	StructField    simpleStruct
	StructPtrField *simpleStruct
	BoolField      bool
}

var strPtr = new(string)
var boolPtr = new(bool)
var int8Ptr = new(int8)
var int16Ptr = new(int16)
var int32Ptr = new(int32)
var int64Ptr = new(int64)
var uint8Ptr = new(uint8)
var uint16Ptr = new(uint16)
var uint32Ptr = new(uint32)
var uint64Ptr = new(uint64)
var float32Ptr = new(float32)
var float64Ptr = new(float64)

var dataResultArr = []testCase{
	{Input: strPtr, Output: "hello"},
	{Input: boolPtr, Output: true},
	{Input: int8Ptr, Output: int8(123)},
	{Input: int16Ptr, Output: int16(123)},
	{Input: int32Ptr, Output: int32(123)},
	{Input: int64Ptr, Output: int64(123)},
	{Input: uint8Ptr, Output: uint8(123)},
	{Input: uint16Ptr, Output: uint16(123)},
	{Input: uint32Ptr, Output: uint32(123)},
	{Input: uint64Ptr, Output: uint64(123)},
	{Input: float32Ptr, Output: float32(123.45)},
	{Input: float64Ptr, Output: float64(123.45)},
	{Input: "hello", Output: "hello"},
	{Input: 123, Output: 123},
	{Input: 123.45, Output: 123.45},
	{Input: []string{"one", "two", "three"}, Output: []any{"one", "two", "three"}},
	{Input: []*string{strPtr}, Output: []any{"hello"}},
	{Input: []int{1, 2, 3}, Output: []any{1, 2, 3}},
	{
		Input:  simpleStruct{StrField: "hello", IntField: 123, BoolField: true},
		Output: map[string]any{"StrField": "hello", "int_field_with_json_name": 123, "BoolField": true},
	},
	{
		Input:  &simpleStruct{StrField: "hello", IntField: 123, BoolField: true},
		Output: map[string]any{"StrField": "hello", "int_field_with_json_name": 123, "BoolField": true},
	},
	{
		Input: complexStruct{
			StructField:    simpleStruct{StrField: "hola", IntField: 321, BoolField: false},
			StructPtrField: &simpleStruct{StrField: "oi", IntField: 000, BoolField: true},
			BoolField:      false,
		},
		Output: map[string]any{
			"StructField":    map[string]any{"StrField": "hola", "int_field_with_json_name": 321, "BoolField": false},
			"StructPtrField": map[string]any{"StrField": "oi", "int_field_with_json_name": 000, "BoolField": true},
			"BoolField":      false,
		},
	},
	{
		Input: &complexStruct{
			StructField:    simpleStruct{StrField: "hola", IntField: 321, BoolField: false},
			StructPtrField: &simpleStruct{StrField: "oi", IntField: 000, BoolField: true},
			BoolField:      false,
		},
		Output: map[string]any{
			"StructField":    map[string]any{"StrField": "hola", "int_field_with_json_name": 321, "BoolField": false},
			"StructPtrField": map[string]any{"StrField": "oi", "int_field_with_json_name": 000, "BoolField": true},
			"BoolField":      false,
		},
	},
	{
		Input: []complexStruct{
			{
				StructField:    simpleStruct{StrField: "hola", IntField: 321, BoolField: false},
				StructPtrField: &simpleStruct{StrField: "oi", IntField: 000, BoolField: true},
				BoolField:      false,
			},
			{
				StructField:    simpleStruct{StrField: "ciao", IntField: 456, BoolField: true},
				StructPtrField: &simpleStruct{StrField: "oi", IntField: 999, BoolField: false},
				BoolField:      true,
			},
		},
		Output: []any{
			map[string]any{
				"StructField":    map[string]any{"StrField": "hola", "int_field_with_json_name": 321, "BoolField": false},
				"StructPtrField": map[string]any{"StrField": "oi", "int_field_with_json_name": 000, "BoolField": true},
				"BoolField":      false,
			},
			map[string]any{
				"StructField":    map[string]any{"StrField": "ciao", "int_field_with_json_name": 456, "BoolField": true},
				"StructPtrField": map[string]any{"StrField": "oi", "int_field_with_json_name": 999, "BoolField": false},
				"BoolField":      true,
			},
		},
	},
	{
		Input: []*complexStruct{
			{
				StructField:    simpleStruct{StrField: "hola", IntField: 321, BoolField: false},
				StructPtrField: &simpleStruct{StrField: "oi", IntField: 000, BoolField: true},
				BoolField:      false,
			},
			{
				StructField:    simpleStruct{StrField: "ciao", IntField: 456, BoolField: true},
				StructPtrField: &simpleStruct{StrField: "oi", IntField: 999, BoolField: false},
				BoolField:      true,
			},
		},
		Output: []any{
			map[string]any{
				"StructField":    map[string]any{"StrField": "hola", "int_field_with_json_name": 321, "BoolField": false},
				"StructPtrField": map[string]any{"StrField": "oi", "int_field_with_json_name": 000, "BoolField": true},
				"BoolField":      false,
			},
			map[string]any{
				"StructField":    map[string]any{"StrField": "ciao", "int_field_with_json_name": 456, "BoolField": true},
				"StructPtrField": map[string]any{"StrField": "oi", "int_field_with_json_name": 999, "BoolField": false},
				"BoolField":      true,
			},
		},
	},
	{Input: func() {}, Output: nil},
}

var expectedErrors = []reflect.Kind{
	reflect.Invalid,
	reflect.Chan,
	reflect.Complex64,
	reflect.Complex128,
	reflect.Func,
	reflect.UnsafePointer,
	reflect.Interface,
}

func TestSimplifyAny(t *testing.T) {
	for _, testCase := range dataResultArr {
		simple, err := SimplifyAny(testCase.Input)
		if err != nil && !slices.Contains(expectedErrors, reflect.TypeOf(testCase.Input).Kind()) {
			t.Error(err)
		}
		if !reflect.DeepEqual(simple, testCase.Output) {
			t.Errorf("Failed to convert value %+v, of type %T, to %+v. Got %+v instead", testCase.Input, testCase.Input, testCase.Output, simple)
		}
	}
}

func init() {
	*strPtr = "hello"
	*boolPtr = true
	*int8Ptr = 123
	*int16Ptr = int16(123)
	*int32Ptr = int32(123)
	*int64Ptr = int64(123)
	*uint8Ptr = uint8(123)
	*uint16Ptr = uint16(123)
	*uint32Ptr = uint32(123)
	*uint64Ptr = uint64(123)
	*float32Ptr = float32(123.45)
	*float64Ptr = float64(123.45)
}
