package utils

import "reflect"

// Checks if the value is exactly the same, it's basically a "=="
// but handles map/slice and other non-comparable types as pointers,
// that is if their underlying pointer is the same.
// It DOES NOT check for map/slice similarities, only their addresses.
func IsSameValueOrPointer(a, b any) bool {
	if a == nil {
		return b == nil
	}

	vA := reflect.ValueOf(a)
	vB := reflect.ValueOf(b)
	if vA.Type() != vB.Type() {
		return false
	}

	vA = valueDereference(vA)
	vB = valueDereference(vB)

	if vA.Comparable() {
		return vA.Interface() == vB.Interface()
	} else {
		return vA.UnsafePointer() == vB.UnsafePointer()
	}
}

func IsComparableEqual[V comparable](a, b V) bool {
	return a == b
}

func IsComparablePointerEqual[V comparable](a, b *V) bool {
	return IsPointerEqualFunc[V](a, b, func(v1, v2 *V) bool {
		return IsComparableEqual[V](*a, *b)
	})
}

// equals is only called if both pointers are non-nil
func IsPointerEqualFunc[V any](a, b *V, equals func(*V, *V) bool) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return equals(a, b)
}

func valueDereference(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Pointer && !value.IsNil() {
		value = value.Elem()
	}
	return value
}
