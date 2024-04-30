package provider

import (
	"context"
	"fmt"

	mgcSchemaPkg "magalu.cloud/core/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func applyMgcMapToTFState(ctx context.Context, mgcMap map[string]any, schema *mgcSchemaPkg.Schema, attrInfoMap resAttrInfoMap, tfState *tfsdk.State) Diagnostics {
	resInfo := &resAttrInfo{
		tfName:          "tfState",
		mgcName:         "tfState",
		mgcSchema:       schema,
		childAttributes: attrInfoMap,
	}
	return applyMgcObject(ctx, mgcMap, resInfo, tfState, path.Empty())
}

func applyMgcMap(ctx context.Context, mgcMap map[string]any, attr *resAttrInfo, tfState *tfsdk.State, path path.Path) Diagnostics {
	diagnostics := Diagnostics{}
	tflog.Debug(
		ctx,
		"[applier] will apply as map",
		map[string]any{"mgcName": attr.mgcName, "tfName": attr.tfName, "value": mgcMap},
	)

	d := tfState.SetAttribute(ctx, path, mgcMap)
	if diagnostics.AppendCheckError(d...) {
		return diagnostics
	}

	return diagnostics
}

func applyMgcObject(ctx context.Context, mgcValue any, attr *resAttrInfo, tfState *tfsdk.State, path path.Path) Diagnostics {
	diagnostics := Diagnostics{}
	tflog.Debug(
		ctx,
		"[applier] will apply as object",
		map[string]any{"mgcName": attr.mgcName, "tfName": attr.tfName, "value": mgcValue, "mgcSchema": attr.mgcSchema},
	)

	if isSchemaXOfObject(attr.mgcSchema) {
		return applyMgcXOfObject(ctx, mgcValue, attr, tfState, path)
	}

	mgcMap, ok := mgcValue.(map[string]any)
	if !ok {
		return diagnostics.AppendLocalErrorReturn(
			"value is not a map",
			fmt.Sprintf("value is not a map: %#v", mgcValue),
		)
	}

	if attr.mgcSchema.AdditionalProperties.Has != nil && *attr.mgcSchema.AdditionalProperties.Has {
		return diagnostics.AppendLocalErrorReturn(
			"Unable to apply map to TF State",
			fmt.Sprintf(
				"Schema for %q with AdditionalProperties must have type information, and not just a boolean. Value: %#v Schema: %#v",
				attr.tfName,
				mgcMap,
				attr.mgcSchema,
			),
		)
	}

	if attr.mgcSchema.AdditionalProperties.Schema != nil {
		return applyMgcMap(ctx, mgcMap, attr, tfState, path)
	}

	for childName, childAttr := range attr.childAttributes {
		tflog.Debug(
			ctx,
			fmt.Sprintf("[applier] looping through object properties, current: %q", childName),
			map[string]any{"mgcName": attr.mgcName, "tfName": attr.tfName, "value": mgcValue, "mgcSchema": attr.mgcSchema},
		)

		childValue, ok := mgcMap[string(childName)]
		if !ok {
			continue
		}

		tflog.Debug(
			ctx,
			"[applier] will try to apply map property",
			map[string]any{"propMgcName": childName, "propMgcValue": childValue},
		)

		tflog.Debug(ctx, fmt.Sprintf("applying %q attribute in state", childName), map[string]any{"value": childValue})

		childPath := path.AtName(string(childAttr.tfName))
		d := applyValueToState(ctx, childValue, childAttr, tfState, childPath)

		if diagnostics.AppendCheckError(d...) {
			childAttrSchema, _ := tfState.Schema.AttributeAtPath(ctx, childPath)
			diagnostics.AddLocalAttributeError(
				childPath,
				"unable to load value",
				fmt.Sprintf("path: %#v - value: %#v - tfschema: %#v", childPath, mgcValue, childAttrSchema),
			)
			return diagnostics
		}
	}
	return diagnostics
}

func applyMgcXOfObject(ctx context.Context, mgcValue any, attr *resAttrInfo, tfState *tfsdk.State, path path.Path) Diagnostics {
	diagnostics := Diagnostics{}
	tflog.Debug(
		ctx,
		"[applier] object is xof schema",
		map[string]any{"mgcName": attr.mgcName, "tfName": attr.tfName, "value": mgcValue, "mgcSchema": attr.mgcSchema},
	)

	propThatFits := ""

	for propName, propSchemaRef := range attr.mgcSchema.Properties {
		propSchema := propSchemaRef.Value
		if propSchema.VisitJSON(mgcValue) != nil {
			continue
		}

		propAttr, ok := attr.childAttributes[mgcName(propName)]
		if !ok {
			return diagnostics.AppendLocalErrorReturn(
				fmt.Sprintf("Unable to set property %q to state", propName),
				"Couldn't find schema definition for property. This probably means that the property isn't expected at all in this resource",
			)
		}

		attrPath := path.AtName(string(propName))
		d := applyValueToState(ctx, mgcValue, propAttr, tfState, attrPath)
		if diagnostics.AppendCheckError(d...) {
			return diagnostics
		}

		propThatFits = propName
	}

	if propThatFits == "" {
		return diagnostics.AppendLocalErrorReturn(
			"value does not fit any xof alternatives",
			fmt.Sprintf("value does not fit any xof alternatives: %#v", mgcValue),
		)
	}

	for propName := range attr.mgcSchema.Properties {
		if propName == propThatFits {
			continue
		}

		tfState.SetAttribute(ctx, path.AtName(propName), nil)
	}

	return diagnostics
}

func applyMgcList(ctx context.Context, mgcValue any, attr *resAttrInfo, tfState *tfsdk.State, path path.Path) Diagnostics {
	diagnostics := Diagnostics{}
	attr = attr.childAttributes["0"]

	// This shouldn't happen, probably, but sometimes the Services return null values for non-nullable values
	if mgcValue == nil {
		d := tfState.SetAttribute(ctx, path, []any{})
		return diagnostics.AppendReturn(Diagnostics(d).DemoteErrorsToWarnings()...)
	}

	mgcList, ok := mgcValue.([]any)
	if !ok {
		diagnostics.AppendReturn(NewLocalErrorDiagnostic(
			fmt.Sprintf("Unable to apply list property %q to State, value is not list", attr.tfName),
			fmt.Sprintf("Property value received from service was not a list: %#v", mgcValue),
		))
	}

	// First overwrite the current list values completely, so empty list
	d := tfState.SetAttribute(ctx, path, []any{})
	if diagnostics.AppendCheckError(d...) {
		return diagnostics
	}

	if len(mgcList) == 0 {
		return diagnostics
	}

	for i, mgcValue := range mgcList {
		attrPath := path.AtListIndex(i)
		d := applyValueToState(ctx, mgcValue, attr, tfState, attrPath)
		if diagnostics.AppendCheckError(d...) {
			attrSchema, _ := tfState.Schema.AttributeAtPath(ctx, attrPath)
			diagnostics.AddLocalAttributeError(attrPath, "unable to load value", fmt.Sprintf("path: %#v - value: %#v - tfschema: %#v", attrPath, mgcValue, attrSchema))
			return diagnostics
		}
	}

	return diagnostics
}

func applyValueToState(ctx context.Context, mgcValue any, attr *resAttrInfo, tfState *tfsdk.State, path path.Path) Diagnostics {
	tflog.Debug(
		ctx,
		"[applier] starting applying mgc value to TF state",
		map[string]any{"mgcName": attr.mgcName, "tfName": attr.tfName, "value": mgcValue},
	)

	switch attr.mgcSchema.Type {
	case "array":
		tflog.Debug(ctx, fmt.Sprintf("populating list in state at path %#v", path))
		return applyMgcList(ctx, mgcValue, attr, tfState, path)

	case "object":
		tflog.Debug(ctx, fmt.Sprintf("populating nested object in state at path %#v", path))
		return applyMgcObject(ctx, mgcValue, attr, tfState, path)

	default:
		if mgcValue == nil {
			// We must check the nil value type, since SetAttribute method requires a typed nil
			switch attr.mgcSchema.Type {
			case "string":
				mgcValue = (*string)(nil)
			case "integer":
				mgcValue = (*int64)(nil)
			case "number":
				mgcValue = (*float64)(nil)
			case "boolean":
				mgcValue = (*bool)(nil)
			}
		}

		// Should this be a local error? Does TF know it already, since it's their function?
		d := tfState.SetAttribute(ctx, path, mgcValue)
		return Diagnostics(d)
	}
}
