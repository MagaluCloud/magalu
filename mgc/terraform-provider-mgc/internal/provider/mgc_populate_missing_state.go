package provider

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
)

type MgcPopulateUnknownState struct {
	resourceName tfName
	ignoreKeys   []tfName
}

func newMgcPopulateUnknownState(
	resourceName tfName,
	ignoreKeys []tfName,
) MgcOperation {
	return &MgcPopulateUnknownState{
		resourceName: resourceName,
		ignoreKeys:   ignoreKeys,
	}
}

func (o *MgcPopulateUnknownState) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "populate unknown state")
	ctx = tflog.SetField(ctx, connectionResourceNameField, o.resourceName)
	return ctx
}

func (o *MgcPopulateUnknownState) CollectParameters(ctx context.Context, state, plan TerraformParams) (core.Parameters, Diagnostics) {
	parameters := core.Parameters{}
	for plannedKey, plannedVal := range state {
		if _, ok := plan[plannedKey]; ok || slices.Contains(o.ignoreKeys, plannedKey) {
			continue
		}
		parameters[string(plannedKey)] = plannedVal
	}
	return parameters, nil
}

func (o *MgcPopulateUnknownState) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return nil, nil
}

func (o *MgcPopulateUnknownState) ShouldRun(context.Context, core.Parameters, core.Configs) (run bool, d Diagnostics) {
	return true, d
}

func (o *MgcPopulateUnknownState) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	return core.NewSimpleResult(
		core.ResultSource{Parameters: params},
		mgcSchemaPkg.NewAnySchema(),
		params,
	), nil
}

func (o *MgcPopulateUnknownState) PostRun(
	ctx context.Context,
	result core.ResultWithValue,
	state, plan TerraformParams,
	targetState *tfsdk.State,
) (runChain bool, diagnostics Diagnostics) {
	diagnostics = Diagnostics{}

	tflog.Info(ctx, "populating unknown state values")

	for paramTFName, paramValue := range result.Source().Parameters {
		paramTFValue := paramValue.(tftypes.Value)

		if !paramTFValue.IsKnown() {
			continue
		}

		attrPath := path.Empty().AtName(string(paramTFName))
		attr, d := targetState.Schema.AttributeAtPath(ctx, attrPath)
		if diagnostics.AppendCheckError(d...) {
			return false, diagnostics
		}

		attrValue, err := attr.GetType().ValueFromTerraform(ctx, paramTFValue)

		if err != nil {
			return false, diagnostics.AppendLocalErrorReturn(
				"Unable to pre-populate Response State with Plan",
				fmt.Sprintf("Attribute %q returned error %v", paramTFName, err),
			)
		}

		tflog.Debug(ctx, fmt.Sprintf("populating attribute %q in response state with known value", paramTFName))

		if !attrValue.IsUnknown() {
			tflog.Debug(ctx, fmt.Sprintf("attribute %q is already known in response state, ignoring", paramTFName))
			continue
		}

		tflog.Debug(ctx, fmt.Sprintf("populating attribute %q in response state with known value", paramTFName))

		d = targetState.SetAttribute(ctx, attrPath, attrValue)
		if diagnostics.AppendCheckError(d...) {
			return false, diagnostics
		}
		tflog.Debug(ctx, fmt.Sprintf("populated attribute %q in response state with known value %s", paramTFName, attrValue.String()))
	}

	tflog.Info(ctx, "populated unknown state values")
	return false, diagnostics
}

func (o *MgcPopulateUnknownState) ChainOperations(ctx context.Context, _ core.ResultWithValue, state, plan TerraformParams) ([]MgcOperation, bool, Diagnostics) {
	return nil, false, nil
}

var _ MgcOperation = (*MgcPopulateUnknownState)(nil)
