package provider

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	mgcSchemaPkg "magalu.cloud/core/schema"
	mgcSdk "magalu.cloud/sdk"
)

type tfStateConverter struct {
	ctx      context.Context
	diag     *diag.Diagnostics
	tfSchema *schema.Schema
}

func newTFStateConverter(ctx context.Context, diag *diag.Diagnostics, tfSchema *schema.Schema) tfStateConverter {
	return tfStateConverter{
		ctx:      ctx,
		diag:     diag,
		tfSchema: tfSchema,
	}
}

func areObjectsSubsets(ctx context.Context, objA *mgcSdk.Schema, objB *mgcSdk.Schema) (bool, error) {
	superset := objB
	subset := objA
	if len(objA.Properties) > len(objB.Properties) {
		superset = objA
		subset = objB
	}

	for propKey, subsetPropRef := range subset.Properties {
		if prop, ok := superset.Properties[propKey]; !ok {
			return false, nil
		} else {
			if slices.Contains(superset.Required, propKey) != slices.Contains(subset.Required, propKey) {
				return false, nil
			}

			if !checkSimilarJsonSchemas(ctx, (*mgcSchemaPkg.Schema)(subsetPropRef.Value), (*mgcSchemaPkg.Schema)(prop.Value)) {
				return false, nil
			}
		}

	}
	return true, nil
}

func isAnyOfElementsSubsets(ctx context.Context, v *mgcSdk.Schema) (string, *mgcSdk.Schema, error) {
	types := map[string]struct{}{}
	maxPropsQty := 0
	var biggestObject *mgcSdk.Schema = nil

	if len(v.AnyOf) == 0 {
		return "", nil, fmt.Errorf("schema doesn't contain anyOf values")
	}

	for _, ref := range v.AnyOf {
		c := (*mgcSdk.Schema)(ref.Value)
		t, err := getJsonType(c)
		if err != nil {
			return "", nil, fmt.Errorf("anyOf element type %q cannot be identified. %w", t, err)
		}

		types[t] = struct{}{}
		if len(types) > 1 {
			// Different types cannot be subsets
			return "", nil, fmt.Errorf("anyOf element of different types cannot be subsets. types: %#v", maps.Keys(types))
		}

		// IS THIS ANOTHER FUNCTION?
		if _, ok := types["object"]; ok {
			if biggestObject == nil {
				biggestObject = c
				maxPropsQty = len(c.Properties)
				continue
			}

			propsQty := len(c.Properties)
			if check, err := areObjectsSubsets(ctx, biggestObject, c); err != nil {
				return "", nil, err
			} else if !check {
				return "", nil, fmt.Errorf("anyOf elements are not subsets")
			}
			if propsQty > maxPropsQty {
				maxPropsQty = propsQty
				biggestObject = c
			}
		}
	}

	if biggestObject != nil {
		// TODO: Handle anyOfs nested in anyOfs? Maybe a weird case of single element anyOf with anyOfs?
		t, err := getJsonType(biggestObject)
		return t, biggestObject, err
	}

	return maps.Keys(types)[0], (*mgcSdk.Schema)(v.AnyOf[0].Value), nil
}

func getJsonEnumType(v *mgcSdk.Schema) (string, error) {
	types := []string{}
	for _, v := range v.Enum {
		var t string
		switch v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			t = "integer"
		case float32, float64:
			t = "number"
		case string:
			t = "string"
		case bool:
			t = "boolean"
		default:
			return "", fmt.Errorf("unsupported enum value: %+v", v)
		}
		if !slices.Contains(types, t) {
			types = append(types, t)
		}
	}
	if len(types) != 1 {
		return "", fmt.Errorf("must provide values of a single type in a enum, got %+v", types)
	}

	return types[0], nil
}

func getJsonType(v *mgcSdk.Schema) (string, error) {
	if v.Type == "" {
		if len(v.Enum) != 0 {
			return getJsonEnumType(v)
		}

		if len(v.AnyOf) != 0 {
			// SHOULD HANDLE SUBSET IN OTHER PLACES? JUST RETURN ANYOF HERE AND DEAL WITH TYPE ELSEWHERE?
			return "anyOf", nil
			// return isAnyOfElementsSubsets(v)

			// if value, err := isAnyOfElementsSubsets(v); err != nil {
			// 	return "", fmt.Errorf("anyOf types are not subsets, they can't be handled. Error: %w", err)
			// } else if value != "" {
			// 	return value, nil
			// }
		}

		return "", fmt.Errorf("unable to find schema %+v type", v)
	}
	return v.Type, nil
}

func (c *tfStateConverter) toMgcSchemaValue(atinfo *attribute, tfValue tftypes.Value, ignoreUnknown bool, filterUnset bool) (mgcValue any, isKnown bool) {
	tflog.Debug(
		c.ctx,
		"[convert] starting conversion from TF state value to mgc value",
		map[string]any{"mgcName": atinfo.mgcName, "tfName": atinfo.tfName, "value": tfValue},
	)
	mgcSchema := atinfo.mgcSchema
	if mgcSchema == nil {
		c.diag.AddError("Invalid schema", "null schema provided to convert state to go values")
		return nil, false
	}

	if !tfValue.IsKnown() {
		if !ignoreUnknown {
			c.diag.AddError(
				"Unable to convert unknown value",
				fmt.Sprintf("[convert] unable to convert %q since value is unknown: value %+v - schema: %+v", atinfo.mgcName, tfValue, mgcSchema),
			)
			return nil, false
		}
		return nil, false
	}

	if tfValue.IsNull() {
		if atinfo.tfSchema.IsOptional() && !atinfo.tfSchema.IsComputed() {
			// Optional values that aren't computed will never be unknown
			// this means they will be null in the state
			return nil, true
		} else if !mgcSchemaPkg.IsSchemaNullable(mgcSchema) {
			c.diag.AddError(
				"Unable to convert non nullable value",
				fmt.Sprintf("[convert] unable to convert %q since value is null and not nullable by the schema: value %+v - schema: %+v", atinfo.mgcName, tfValue, mgcSchema),
			)
			return nil, true
		}
		return nil, true
	}

	t, err := getJsonType(mgcSchema)
	if err != nil {
		c.diag.AddError(fmt.Sprintf("Unable to get schema type for attribute %q", atinfo.mgcName), err.Error())
		return nil, false
	}

	switch t {
	case "string":
		var state string
		err := tfValue.As(&state)
		if err != nil {
			c.diag.AddError(
				"Unable to convert value to string",
				fmt.Sprintf("[convert] unable to convert %q with value %+v to schema %+v - error: %s", atinfo.mgcName, tfValue, mgcSchema, err.Error()),
			)
			return nil, true
		}
		tflog.Debug(c.ctx, "[convert] finished conversion to string", map[string]any{"resulting value": state})
		return state, true
	case "number":
		var state big.Float
		err := tfValue.As(&state)
		if err != nil {
			c.diag.AddError(
				"Unable to convert value to number",
				fmt.Sprintf("[convert] unable to convert %q with value %+v to schema %+v - error: %s", atinfo.mgcName, tfValue, mgcSchema, err.Error()),
			)
			return nil, true
		}

		result, accuracy := state.Float64()
		if accuracy != big.Exact {
			c.diag.AddError(
				"Unable to convert value to float",
				fmt.Sprintf("[convert] %q with value %+v lost accuracy in conversion to %+v", atinfo.mgcName, state, result),
			)
			return nil, true
		}
		tflog.Debug(c.ctx, "[convert] finished conversion to number", map[string]any{"resulting value": result})
		return result, true
	case "integer":
		var state big.Float
		err := tfValue.As(&state)
		if err != nil {
			c.diag.AddError(
				"Unable to convert value to integer",
				fmt.Sprintf("[convert] unable to convert %q with value %+v to schema %+v - error: %s", atinfo.mgcName, tfValue, mgcSchema, err.Error()),
			)
			return nil, true
		}

		result, accuracy := state.Int64()
		if accuracy != big.Exact {
			c.diag.AddError(
				"Unable to convert value to integer",
				fmt.Sprintf("[convert] %q with value %+v lost accuracy in conversion to %+v", atinfo.mgcName, state, result),
			)
			return nil, true
		}
		tflog.Debug(c.ctx, "[convert] finished conversion to integer", map[string]any{"resulting value": result})
		return result, true
	case "boolean":
		var state bool
		err := tfValue.As(&state)
		if err != nil {
			c.diag.AddError(
				"Unable to convert value to boolean",
				fmt.Sprintf("[convert] unable to convert %q with value %+v to schema %+v - error: %s", atinfo.mgcName, tfValue, mgcSchema, err.Error()),
			)
			return nil, true
		}
		tflog.Debug(c.ctx, "[convert] finished conversion to bool", map[string]any{"resulting value": state})
		return state, true
	case "array":
		return c.toMgcSchemaArray(atinfo, tfValue, ignoreUnknown, filterUnset)
	case "object":
		return c.toMgcSchemaMap(atinfo, tfValue, ignoreUnknown, filterUnset)
	default:
		c.diag.AddError("Unknown value", fmt.Sprintf("[convert] unable to convert %q with value %+v to schema %+v", atinfo.mgcName, tfValue, mgcSchema))
		return nil, false
	}
}

func (c *tfStateConverter) toMgcSchemaArray(atinfo *attribute, tfValue tftypes.Value, ignoreUnknown bool, filterUnset bool) (mgcArray []any, isKnown bool) {
	tflog.Debug(
		c.ctx,
		"[convert] starting conversion from TF state value to mgc array",
		map[string]any{"mgcName": atinfo.mgcName, "tfName": atinfo.tfName, "value": tfValue},
	)
	mgcSchema := atinfo.mgcSchema
	var tfArray []tftypes.Value
	err := tfValue.As(&tfArray)
	if err != nil {
		c.diag.AddError(
			"Unable to convert value to list",
			fmt.Sprintf("[convert] unable to convert %q with value %+v to schema %+v - error: %s", atinfo.mgcName, tfValue, mgcSchema, err.Error()),
		)
		return nil, false
	}

	// TODO: Handle attribute information in list - it should be mapped to "0" key
	itemAttr := atinfo.attributes["0"]
	mgcArray = make([]any, len(tfArray))
	isKnown = true
	for i, tfItem := range tfArray {
		mgcItem, isItemKnown := c.toMgcSchemaValue(itemAttr, tfItem, ignoreUnknown, filterUnset)
		if c.diag.HasError() {
			c.diag.AddError("Unable to convert array", fmt.Sprintf("unknown value inside %q array at %v", atinfo.mgcName, i))
			return nil, isItemKnown
		}
		if !isItemKnown {
			// TODO: confirm this logic, should we just keep going?
			c.diag.AddWarning("Unknown list item", fmt.Sprintf("Item %d in %q is unknown: %+v", i, atinfo.mgcName, tfItem))
			isKnown = false
			return
		}
		mgcArray[i] = mgcItem
	}
	tflog.Debug(c.ctx, "[convert] finished conversion to array", map[string]any{"resulting value": mgcArray})
	return
}

func (c *tfStateConverter) toMgcSchemaMap(atinfo *attribute, tfValue tftypes.Value, ignoreUnknown bool, filterUnset bool) (mgcMap map[string]any, isKnown bool) {
	tflog.Debug(
		c.ctx,
		"[convert] starting conversion from TF state value to mgc map",
		map[string]any{"mgcName": atinfo.mgcName, "tfName": atinfo.tfName, "value": tfValue},
	)
	mgcSchema := atinfo.mgcSchema
	var tfMap map[string]tftypes.Value
	err := tfValue.As(&tfMap)
	if err != nil {
		c.diag.AddError(
			"Unable to convert value to map",
			fmt.Sprintf("[convert] unable to convert %q with value %+v to schema %+v - error: %s", atinfo.mgcName, tfValue, mgcSchema, err.Error()),
		)
		return nil, false
	}

	mgcMap = map[string]any{}
	isKnown = true
	for attr := range mgcSchema.Properties {
		mgcName := mgcName(attr)
		itemAttr := atinfo.attributes[mgcName]
		if itemAttr == nil {
			c.diag.AddError(
				"Schema attribute missing from attribute information",
				fmt.Sprintf("[convert] schema attribute %q doesn't have attribute information", mgcName),
			)
			continue
		}

		tfName := itemAttr.tfName
		tfItem, ok := tfMap[string(tfName)]
		if !ok {
			title := "Schema attribute missing from state value"
			msg := fmt.Sprintf("[convert] schema attribute %q with info `%+v` missing from state %+v", mgcName, atinfo, tfMap)
			if itemAttr.tfSchema.IsRequired() {
				c.diag.AddError(title, msg)
				return
			}
			tflog.Debug(c.ctx, msg)
			continue
		}

		mgcItem, isItemKnown := c.toMgcSchemaValue(itemAttr, tfItem, ignoreUnknown, filterUnset)
		if c.diag.HasError() {
			return nil, false
		}

		if !isItemKnown && ignoreUnknown {
			continue
		}
		if mgcItem == nil && filterUnset {
			continue
		}

		mgcMap[string(mgcName)] = mgcItem
	}
	tflog.Debug(c.ctx, "[convert] finished conversion to map", map[string]any{"resulting value": mgcMap})
	return
}

// Read values from tfValue into a map suitable to MGC
func (c *tfStateConverter) readMgcMap(mgcSchema *mgcSdk.Schema, attributes mgcAttributes, tfState tfsdk.State) (mgcMap map[string]any) {
	attr := &attribute{
		tfName:     "inputSchemasInfo",
		mgcName:    "inputSchemasInfo",
		mgcSchema:  mgcSchema,
		attributes: attributes,
	}

	m, _ := c.toMgcSchemaMap(attr, tfState.Raw, true, true)
	return m
}

func (c *tfStateConverter) applyMgcMap(mgcMap map[string]any, attributes mgcAttributes, ctx context.Context, tfState *tfsdk.State, path path.Path) {
	for mgcName, attr := range attributes {
		mgcValue, ok := mgcMap[string(mgcName)]
		if !ok {
			// Ignore non existing values
			continue
		}

		tflog.Debug(ctx, fmt.Sprintf("applying %q attribute in state", mgcName), map[string]any{"value": mgcValue})

		attrPath := path.AtName(string(attr.tfName))
		c.applyValueToState(mgcValue, attr, ctx, tfState, attrPath)

		if c.diag.HasError() {
			attrSchema, _ := tfState.Schema.AttributeAtPath(ctx, attrPath)
			c.diag.AddAttributeError(
				attrPath,
				"unable to convert value",
				fmt.Sprintf("path: %#v - value: %#v - tfschema: %#v", attrPath, mgcValue, attrSchema),
			)
			return
		}
	}
}

func (c *tfStateConverter) applyMgcList(mgcList []any, attributes mgcAttributes, ctx context.Context, tfState *tfsdk.State, path path.Path) {
	attr := attributes["0"]

	for i, mgcValue := range mgcList {
		attrPath := path.AtListIndex(i)
		c.applyValueToState(mgcValue, attr, ctx, tfState, attrPath)

		if c.diag.HasError() {
			attrSchema, _ := tfState.Schema.AttributeAtPath(ctx, attrPath)
			c.diag.AddAttributeError(attrPath, "unable to convert value", fmt.Sprintf("path: %#v - value: %#v - tfschema: %#v", attrPath, mgcValue, attrSchema))
			return
		}
	}
}

func (c *tfStateConverter) applyValueToState(mgcValue any, attr *attribute, ctx context.Context, tfState *tfsdk.State, path path.Path) {
	rv := reflect.ValueOf(mgcValue)
	t, err := getJsonType(attr.mgcSchema)
	if err != nil {
		c.diag.AddError("Unable to retrieve type", fmt.Sprintf("found an untyped nil attribute `%#v` without valid mgc schema type. Error: %#v", path, err))
	}

	if mgcValue == nil {
		// We must check the nil value type, since SetAttribute method requires a typed nil
		switch t {
		case "string":
			rv = reflect.ValueOf((*string)(nil))
		case "integer":
			rv = reflect.ValueOf((*int64)(nil))
		case "number":
			rv = reflect.ValueOf((*float64)(nil))
		case "boolean":
			rv = reflect.ValueOf((*bool)(nil))
		}
	}

	switch t {
	case "array":
		tflog.Debug(ctx, fmt.Sprintf("populating list in state at path %#v", path))
		c.applyMgcList(mgcValue.([]any), attr.attributes, ctx, tfState, path)

	case "object":
		tflog.Debug(ctx, fmt.Sprintf("populating nested object in state at path %#v", path))
		c.applyMgcMap(mgcValue.(map[string]any), attr.attributes, ctx, tfState, path)

	default:
		c.diag.Append(tfState.SetAttribute(ctx, path, rv.Interface())...)
	}
}
