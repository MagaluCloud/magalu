package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
)

func verifyCurrentDesiredMismatch(inputAttr resAttrInfoMap, inputMgcMap map[string]any, outputMgcMap map[string]any) Diagnostics {
	diagnostics := Diagnostics{}

	for _, desired := range inputAttr {
		current := desired.currentCounterpart
		if current == nil {
			continue
		}

		input, ok := inputMgcMap[string(desired.mgcName)]
		if !ok {
			continue
		}

		output, ok := outputMgcMap[string(current.mgcName)]
		if !ok {
			continue
		}

		if !reflect.DeepEqual(input, output) {
			diagnostics.AddWarning(
				"current/desired attribute mismatch",
				fmt.Sprintf(
					"Terraform isn't able to verify the equality between %q (%v) and %q (%v) because their structures are different. Assuming success.",
					current.tfName,
					output,
					desired.tfName,
					input,
				),
			)
		}
	}
	return diagnostics
}

// Does not return error, check for 'diag.HasError' to see if operation was successful
func resultAsMap(result core.ResultWithValue) (map[string]any, Diagnostics) {
	if result == nil {
		return nil, nil
	}
	resultMap, ok := result.Value().(map[string]any)
	if !ok {
		return nil, NewLocalErrorDiagnostics(
			"Operation output mismatch",
			fmt.Sprintf("Unable to convert %v to map.", result),
		)
	}

	return resultMap, nil
}

// 'read' parameter may be nil ONLY IF 'result.Source().Executor' already has a Link named "read"
func applyStateAfter(
	ctx context.Context,
	resourceName tfName,
	attrTree resAttrInfoTree,
	result core.ResultWithValue,
	state TerraformParams,
	tfState *tfsdk.State,
) Diagnostics {
	diagnostics := Diagnostics{}
	tflog.Debug(ctx, fmt.Sprintf("[resource] applying state after for %q", resourceName))

	tflog.Debug(ctx, "[resource] applying request parameters in state")
	// First, apply the values that the user passed as Parameters to the state (assuming success)
	d := applyMgcMapToTFState(ctx, result.Source().Parameters, result.Source().Executor.ParametersSchema(), attrTree.input, nil, tfState)
	if diagnostics.AppendCheckError(d...) {
		return diagnostics
	}

	resultMap, d := resultAsMap(result)
	if !d.HasError() {
		tflog.Debug(ctx, "[resource] applying request result in state")
		// Then, apply the result values of the request that was performed
		d := applyMgcMapToTFState(ctx, resultMap, result.Schema(), attrTree.output, state, tfState)
		if diagnostics.AppendCheckError(d...) {
			return diagnostics
		}

		tflog.Debug(ctx, "[resource] checking result state current/desired mismatches")
		d = verifyCurrentDesiredMismatch(attrTree.input, result.Source().Parameters, resultMap)
		if diagnostics.AppendCheckError(d...) {
			return diagnostics
		}
	} else {
		tflog.Debug(ctx, "[resource] request result was not a map, not applying")
	}

	return diagnostics
}
