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

		return "", fmt.Errorf("unable to find schema %+v type", v)
	}
	return v.Type, nil
}

func (c *tfStateConverter) toMgcSchemaValue(atinfo *attribute, tfValue tftypes.Value, ignoreUnknown bool, filterUnset bool) (mgcValue any, isKnown bool) {
	mgcSchema := atinfo.mgcSchema
	if mgcSchema == nil {
		c.diag.AddError("Invalid schema", "null schema provided to convert state to go values")
		return nil, false
	}

	if !tfValue.IsKnown() {
		if !ignoreUnknown {
			c.diag.AddError("Unable to convert unknown value", fmt.Sprintf("[convert] unable to convert since value is unknown: value %+v - schema: %+v", tfValue, mgcSchema))
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
			c.diag.AddError("Unable to convert non nullable value", fmt.Sprintf("[convert] unable to convert since value is null and not nullable by the schema: value %+v - schema: %+v", tfValue, mgcSchema))
			return nil, true
		}
		return nil, true
	}

	t, err := getJsonType(mgcSchema)
	if err != nil {
		c.diag.AddError("Unable to get schema type", err.Error())
		return nil, false
	}

	switch t {
	case "string":
		var state string
		err := tfValue.As(&state)
		if err != nil {
			c.diag.AddError("Unable to convert value to string", fmt.Sprintf("[convert] unable to convert value %+v to schema %+v - error: %s", tfValue, mgcSchema, err.Error()))
			return nil, true
		}
		return state, true
	case "number":
		var state big.Float
		err := tfValue.As(&state)
		if err != nil {
			c.diag.AddError("Unable to convert value to number", fmt.Sprintf("[convert] unable to convert value %+v to schema %+v - error: %s", tfValue, mgcSchema, err.Error()))
			return nil, true
		}

		result, accuracy := state.Float64()
		if accuracy != big.Exact {
			c.diag.AddError("Unable to convert value to float", fmt.Sprintf("[convert] value %+v lost accuracy in conversion to %+v", state, result))
			return nil, true
		}
		return result, true
	case "integer":
		var state big.Float
		err := tfValue.As(&state)
		if err != nil {
			c.diag.AddError("Unable to convert value to integer", fmt.Sprintf("[convert] unable to convert value %+v to schema %+v - error: %s", tfValue, mgcSchema, err.Error()))
			return nil, true
		}

		result, accuracy := state.Int64()
		if accuracy != big.Exact {
			c.diag.AddError("Unable to convert value to integer", fmt.Sprintf("[convert] value %+v lost accuracy in conversion to %+v", state, result))
			return nil, true
		}
		return result, true
	case "boolean":
		var state bool
		err := tfValue.As(&state)
		if err != nil {
			c.diag.AddError("Unable to convert value to boolean", fmt.Sprintf("[convert] unable to convert value %+v to schema %+v - error: %s", tfValue, mgcSchema, err.Error()))
			return nil, true
		}
		return state, true
	case "array":
		return c.toMgcSchemaArray(atinfo, tfValue, ignoreUnknown, filterUnset)
	case "object":
		return c.toMgcSchemaMap(atinfo, tfValue, ignoreUnknown, filterUnset)
	default:
		c.diag.AddError("Unknown value", fmt.Sprintf("[convert] unable to convert value %+v to schema %+v", tfValue, mgcSchema))
		return nil, false
	}
}

func (c *tfStateConverter) toMgcSchemaArray(atinfo *attribute, tfValue tftypes.Value, ignoreUnknown bool, filterUnset bool) (mgcArray []any, isKnown bool) {
	mgcSchema := atinfo.mgcSchema
	var tfArray []tftypes.Value
	err := tfValue.As(&tfArray)
	if err != nil {
		c.diag.AddError("Unable to convert value to list", fmt.Sprintf("[convert] unable to convert value %+v to schema %+v - error: %s", tfValue, mgcSchema, err.Error()))
		return nil, false
	}

	// TODO: Handle attribute information in list - it should be mapped to "0" key
	itemAttr := atinfo.attributes["0"]
	mgcArray = make([]any, len(tfArray))
	isKnown = true
	for i, tfItem := range tfArray {
		mgcItem, isItemKnown := c.toMgcSchemaValue(itemAttr, tfItem, ignoreUnknown, filterUnset)
		if c.diag.HasError() {
			c.diag.AddError("Unable to convert array", fmt.Sprintf("unknown value inside array at %v", i))
			return nil, isItemKnown
		}
		if !isItemKnown {
			// TODO: confirm this logic, should we just keep going?
			c.diag.AddWarning("Unknown list item", fmt.Sprintf("Item %d is unknown: %+v", i, tfItem))
			isKnown = false
			return
		}
		mgcArray[i] = mgcItem
	}
	return
}

func (c *tfStateConverter) toMgcSchemaMap(atinfo *attribute, tfValue tftypes.Value, ignoreUnknown bool, filterUnset bool) (mgcMap map[string]any, isKnown bool) {
	mgcSchema := atinfo.mgcSchema
	var tfMap map[string]tftypes.Value
	err := tfValue.As(&tfMap)
	if err != nil {
		c.diag.AddError("Unable to convert value to map", fmt.Sprintf("[convert] unable to convert value %+v to schema %+v - error: %s", tfValue, mgcSchema, err.Error()))
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
				fmt.Sprintf("[convert] schema attribute `%s` doesn't have attribute information", mgcName),
			)
			continue
		}

		tfName := itemAttr.tfName
		tfItem, ok := tfMap[string(tfName)]
		if !ok {
			title := "Schema attribute missing from state value"
			msg := fmt.Sprintf("[convert] schema attribute `%s` with info `%+v` missing from state %+v", mgcName, atinfo, tfMap)
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

		attrPath := path.AtName(string(attr.tfName))
		c.applyValueToState(mgcValue, attr, ctx, tfState, attrPath)

		if c.diag.HasError() {
			attrSchema, _ := tfState.Schema.AttributeAtPath(ctx, attrPath)
			c.diag.AddAttributeError(attrPath, "unable to convert value", fmt.Sprintf("path: %#v - value: %#v - tfschema: %#v", attrPath, mgcValue, attrSchema))
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
