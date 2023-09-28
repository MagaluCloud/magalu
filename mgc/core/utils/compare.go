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
	if vA.Comparable() {
		return a == b
	} else {
		return vA.UnsafePointer() == vB.UnsafePointer()
	}
}
