package schema

import (
	"testing"
)

var nonNullableData = NewAnyOfSchema(
	NewIntegerSchema(),
	NewStringSchema(),
)
var nullableData = NewAnyOfSchema(
	NewIntegerSchema(),
	NewStringSchema(),
	NewNullSchema(),
)
var refData = NewOneOfSchema(
	NewNullSchema(),
	NewObjectSchema(map[string]*Schema{
		"flex":  NewBooleanSchema(),
		"brand": NewStringSchema(),
	}, []string{"flex", "brand"}),
)

func TestIsSchemaNullable_NotNullable(t *testing.T) {
	isNullable := IsSchemaNullable(nonNullableData)
	if isNullable {
		t.Error("Should detect that type is not nullable")
	}
}

func TestIsSchemaNullable_AnyOf(t *testing.T) {
	isAnyNullable := IsSchemaNullable(nullableData)
	if !isAnyNullable {
		t.Error("Should detect that AnyOf has null type")
	}
}
func TestIsSchemaNullable_OneOf(t *testing.T) {
	isOneNullable := IsSchemaNullable(nullableData)
	if !isOneNullable {
		t.Error("Should detect that OneOf has null type")
	}
}
func TestIsSchemaNullable_AllOf(t *testing.T) {
	isAnyNullable := IsSchemaNullable(nullableData)
	if !isAnyNullable {
		t.Error("Should detect that AllOf has null type")
	}
}

func TestIsSchemaNullable_RefWithNull(t *testing.T) {
	data := NewOneOfSchema(
		refData,
		NewStringSchema(),
	)
	hasRefNullType := IsSchemaNullable(data)
	if !hasRefNullType {
		t.Error("Should detect that refered type can be a null type")
	}
}
func TestIsSchemaNullable_RefNullable(t *testing.T) {
	refData.Nullable = true
	data := NewOneOfSchema(
		refData,
		NewStringSchema(),
	)
	hasRefNullType := IsSchemaNullable(data)
	if !hasRefNullType {
		t.Error("Should detect that refered type is nullable")
	}
}
