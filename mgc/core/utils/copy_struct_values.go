package utils

import (
	"fmt"
	"reflect"
)

func isStructPointer(sP any) bool {
	t := reflect.TypeOf(sP)
	return t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct
}

// Must receive two struct pointers
//
// Copy the value from the src field to the dst field if
// the field exists in both with the same name and the same type.
// The field must also be exported
func CopyStructValues(srcP, dstP any) error {
	if !isStructPointer(srcP) || !isStructPointer(dstP) {
		return fmt.Errorf("src and dst should be a pointer to a struct")
	}

	srcV := reflect.ValueOf(srcP).Elem()
	dstV := reflect.ValueOf(dstP).Elem()

	for i := 0; i < dstV.NumField(); i++ {
		fieldName := dstV.Type().Field(i).Name
		dstField := dstV.FieldByName(fieldName)
		srcField := srcV.FieldByName(fieldName)

		if _, srcFieldFounded := srcV.Type().FieldByName(fieldName); !srcFieldFounded || !dstField.CanSet() || dstField.Type() != srcField.Type() {
			continue
		}

		dstField.Set(srcField)
	}

	return nil
}
