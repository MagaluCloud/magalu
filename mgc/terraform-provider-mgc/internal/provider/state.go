package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"magalu.cloud/core"
	mgcSdk "magalu.cloud/sdk"
)

type tfStateHandler interface {
	Name() string

	TFSchema() *schema.Schema
	ReadResultSchema() *mgcSdk.Schema

	InputAttributes() mgcAttributes
	OutputAttributes() mgcAttributes
	SplitAttributes() []splitMgcAttribute
}

func readMgcMap(handler tfStateHandler, mgcSchema *mgcSdk.Schema, ctx context.Context, tfState tfsdk.State, diag *diag.Diagnostics) map[string]any {
	conv := newTFStateConverter(ctx, diag, handler.TFSchema())
	return conv.readMgcMap(mgcSchema, handler.InputAttributes(), tfState)
}

func applyMgcInputMap(handler tfStateHandler, mgcMap map[string]any, ctx context.Context, tfState *tfsdk.State, diag *diag.Diagnostics) {
	conv := newTFStateConverter(ctx, diag, handler.TFSchema())
	conv.applyMgcMap(mgcMap, handler.InputAttributes(), ctx, tfState, path.Empty())
}

func applyMgcOutputMap(handler tfStateHandler, mgcMap map[string]any, ctx context.Context, tfState *tfsdk.State, diag *diag.Diagnostics) {
	conv := newTFStateConverter(ctx, diag, handler.TFSchema())
	conv.applyMgcMap(mgcMap, handler.OutputAttributes(), ctx, tfState, path.Empty())
}

func verifyCurrentDesiredMismatch(handler tfStateHandler, inputMgcMap map[string]any, outputMgcMap map[string]any, diag *diag.Diagnostics) {
	for _, splitAttr := range handler.SplitAttributes() {
		input, ok := inputMgcMap[string(splitAttr.desired.mgcName)]
		if !ok {
			continue
		}

		output, ok := outputMgcMap[string(splitAttr.current.mgcName)]
		if !ok {
			continue
		}

		if !reflect.DeepEqual(input, output) {
			diag.AddWarning(
				"current/desired attribute mismatch",
				fmt.Sprintf(
					"Terraform isn't able to verify the equality between %q (%v) and %q (%v) because their structures are different. Assuming success.",
					splitAttr.current.tfName,
					output,
					splitAttr.desired.tfName,
					input,
				),
			)
		}
	}
}

// Does not return error, check for 'diag.HasError' to see if operation was successful
func castToMap(result core.ResultWithValue, diag *diag.Diagnostics) map[string]any {
	resultMap, ok := result.Value().(map[string]any)
	if !ok {
		diag.AddError(
			"Operation output mismatch",
			fmt.Sprintf("Unable to convert %v to map.", result),
		)
	}
	return resultMap
}

func applyStateAfter(
	handler tfStateHandler,
	result core.ResultWithValue,
	ctx context.Context,
	tfState *tfsdk.State,
	diag *diag.Diagnostics,
) {
	var resultMap map[string]any
	resultSchema := result.Schema()

	if checkSimilarJsonSchemas(ctx, resultSchema, handler.ReadResultSchema()) {
		if resultMap = castToMap(result, diag); diag.HasError() {
			return
		}
	} else {
		readLink, ok := result.Source().Executor.Links()["read"]
		if !ok {
			diag.AddError(
				"Read link failed",
				fmt.Sprintf("Unable to resolve Read link for applying new state on resource %q. Available links: %v", handler.Name(), result.Source().Executor.Links()),
			)
			return
		}

		additionalParametersSchema := readLink.AdditionalParametersSchema()
		if len(additionalParametersSchema.Required) > 0 {
			diag.AddError("Read link failed", fmt.Sprintf("Unable to resolve parameters on Read link for applying new state on resource %q", handler.Name()))
			return
		}

		exec, err := readLink.CreateExecutor(result)
		if err != nil {
			diag.AddError("Read link failed", fmt.Sprintf("Unable to create Read link executor for applying new state on resource %q: %s", handler.Name(), err))
			return
		}

		result := execute(handler.Name(), ctx, exec, core.Parameters{}, core.Configs{}, diag)
		if diag.HasError() {
			return
		}

		if resultMap = castToMap(result, diag); diag.HasError() {
			return
		}
	}

	// We must apply the input parameters in the state, considering that the request went successfully.
	// BE CAREFUL: Don't apply Plan.Raw values into the State they might be Unknown! State only handles Known/Null values.
	// Also, this must come BEFORE applying the result to the state, as that should override these values when valid.
	applyMgcInputMap(handler, result.Source().Parameters, ctx, tfState, diag)

	applyMgcOutputMap(handler, resultMap, ctx, tfState, diag)
	verifyCurrentDesiredMismatch(handler, result.Source().Parameters, resultMap, diag)
}
