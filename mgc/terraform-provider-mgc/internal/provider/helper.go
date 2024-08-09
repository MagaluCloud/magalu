package provider

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StructToObjectValue(s interface{}) (types.Object, error) {
	return structToObjectValueRecursive(reflect.ValueOf(s))
}

func structToObjectValueRecursive(v reflect.Value) (types.Object, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return types.ObjectNull(map[string]attr.Type{}), fmt.Errorf("input must be a struct or pointer to struct")
	}

	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		attrValue, attrType, err := valueToAttr(fieldValue)
		if err != nil {
			return types.ObjectNull(map[string]attr.Type{}), fmt.Errorf("error processing field %s: %v", field.Name, err)
		}

		attrTypes[field.Name] = attrType
		attrValues[field.Name] = attrValue
	}

	return types.ObjectValueMust(attrTypes, attrValues), nil
}

func valueToAttr(v reflect.Value) (attr.Value, attr.Type, error) {
	switch v.Kind() {
	case reflect.Bool:
		return types.BoolValue(v.Bool()), types.BoolType, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return types.Int64Value(v.Int()), types.Int64Type, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return types.Int64Value(int64(v.Uint())), types.Int64Type, nil
	case reflect.Float32, reflect.Float64:
		return types.Float64Value(v.Float()), types.Float64Type, nil
	case reflect.String:
		return types.StringValue(v.String()), types.StringType, nil
	case reflect.Slice:
		return sliceToAttr(v)
	case reflect.Map:
		return mapToAttr(v)
	case reflect.Struct:
		obj, err := structToObjectValueRecursive(v)
		return obj, types.ObjectType{AttrTypes: obj.AttributeTypes(nil)}, err
	case reflect.Ptr:
		if v.IsNil() {
			return types.StringNull(), types.StringType, nil
		}
		return valueToAttr(v.Elem())
	default:
		return nil, nil, fmt.Errorf("unsupported type: %v", v.Kind())
	}
}

func sliceToAttr(v reflect.Value) (attr.Value, attr.Type, error) {
	if v.Len() == 0 {
		return types.ListNull(types.StringType), types.ListType{ElemType: types.StringType}, nil
	}

	elemType := v.Index(0).Type()
	listType, err := getListElemType(elemType)
	if err != nil {
		return nil, nil, err
	}

	elements := make([]attr.Value, v.Len())
	for i := 0; i < v.Len(); i++ {
		elem, _, err := valueToAttr(v.Index(i))
		if err != nil {
			return nil, nil, err
		}
		elements[i] = elem
	}

	return types.ListValueMust(listType, elements), types.ListType{ElemType: listType}, nil
}

func mapToAttr(v reflect.Value) (attr.Value, attr.Type, error) {
	if v.Len() == 0 {
		return types.MapNull(types.StringType), types.MapType{ElemType: types.StringType}, nil
	}

	elemType := v.Type().Elem()
	mapType, err := getMapElemType(elemType)
	if err != nil {
		return nil, nil, err
	}

	elements := make(map[string]attr.Value)
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		elem, _, err := valueToAttr(iter.Value())
		if err != nil {
			return nil, nil, err
		}
		elements[key] = elem
	}

	return types.MapValueMust(mapType, elements), types.MapType{ElemType: mapType}, nil
}

func getListElemType(t reflect.Type) (attr.Type, error) {
	switch t.Kind() {
	case reflect.Bool:
		return types.BoolType, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return types.Int64Type, nil
	case reflect.Float32, reflect.Float64:
		return types.Float64Type, nil
	case reflect.String:
		return types.StringType, nil
	case reflect.Struct:
		return types.ObjectType{AttrTypes: make(map[string]attr.Type)}, nil
	default:
		return nil, fmt.Errorf("unsupported slice element type: %v", t.Kind())
	}
}

func getMapElemType(t reflect.Type) (attr.Type, error) {
	switch t.Kind() {
	case reflect.Bool:
		return types.BoolType, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return types.Int64Type, nil
	case reflect.Float32, reflect.Float64:
		return types.Float64Type, nil
	case reflect.String:
		return types.StringType, nil
	case reflect.Struct:
		return types.ObjectType{AttrTypes: make(map[string]attr.Type)}, nil
	default:
		return nil, fmt.Errorf("unsupported map element type: %v", t.Kind())
	}
}