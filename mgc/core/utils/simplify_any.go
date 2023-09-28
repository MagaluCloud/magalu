package utils

import (
	"fmt"
	"reflect"
)

// This function takes any Go value (simple types, arrays, maps, structs, pointers...) and converts
// it into the simplest possible native types. Pointers are converted into their respective underlying
// value, structs are converted into map[string]any, etc. This operation is recursive and will convert
// the passed value in its entirety, and not just the first depth level
func SimplifyAny(value any) (converted any, err error) {
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
		return nil, fmt.Errorf("forbidden value type %s", v)

	case reflect.Pointer:
		if v.IsNil() {
			return nil, nil
		}
		return SimplifyAny(v.Elem().Interface())

	case reflect.Array, reflect.Slice:
		return simplifyArray(v)
	case reflect.Map:
		return simplifyMap(v)
	case reflect.Struct:
		resultMap := map[string]any{}
		err := DecodeValue(value, &resultMap)
		if err != nil {
			return nil, err
		}
		return simplifyMap(reflect.ValueOf(resultMap))

	default:
		return nil, fmt.Errorf("unhandled value type: %s", v)
	}
}

func simplifyArray(v reflect.Value) ([]any, error) {
	// convert whatever map to []Value
	count := v.Len()
	result := make([]any, count)
	for i := 0; i < count; i++ {
		subVal := v.Index(i)
		subConverted, err := SimplifyAny(subVal.Interface())
		if err != nil {
			return nil, err
		}
		result[i] = subConverted
	}
	return result, nil
}

func simplifyMap(v reflect.Value) (map[string]any, error) {
	result := make(map[string]any, v.Len())
	keys := v.MapKeys()
	for _, key := range keys {
		sub := v.MapIndex(key)
		subConverted, err := SimplifyAny(sub.Interface())
		if err != nil {
			return nil, err
		}
		keyConverted, err := SimplifyAny(key.Interface())
		if err != nil {
			return nil, err
		}
		keyStr, err := DecodeNewValue[string](keyConverted)
		if err != nil {
			return nil, err
		}
		result[*keyStr] = subConverted
	}
	return result, nil
}
