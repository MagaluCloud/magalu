package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core"
)

type MgcResourceRead struct {
	resourceName tfName
	attrTree     resAttrInfoTree
	operation    core.Executor
}

func newMgcResourceRead(resourceName tfName, attrTree resAttrInfoTree, operation core.Executor) *MgcResourceRead {
	return &MgcResourceRead{resourceName: resourceName, attrTree: attrTree, operation: operation}
}

func (o *MgcResourceRead) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "read")
	ctx = tflog.SetField(ctx, resourceNameField, o.resourceName)
	return ctx
}

func (o *MgcResourceRead) CollectParameters(ctx context.Context, state, _ TerraformParams) (core.Parameters, Diagnostics) {
	return loadMgcParamsFromState(ctx, o.operation.ParametersSchema(), o.attrTree.input, state)
}

func (o *MgcResourceRead) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return getConfigs(ctx, o.operation.ConfigsSchema()), nil
}

func (o *MgcResourceRead) ShouldRun(context.Context, core.Parameters, core.Configs) (run bool, d Diagnostics) {
	return true, d
}

func (o *MgcResourceRead) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	return execute(ctx, o.resourceName, o.operation, params, configs)
}

func (o *MgcResourceRead) PostRun(ctx context.Context, result core.ResultWithValue, state, plan TerraformParams, targetState *tfsdk.State) (runChain bool, diagnostics Diagnostics) {
	tflog.Info(ctx, "resource read")
	diagnostics = Diagnostics{}

	d := applyStateAfter(ctx, o.resourceName, o.attrTree, result, state, targetState)
	if diagnostics.AppendCheckError(d...) {
		return false, diagnostics
	}

	return true, diagnostics
}

func (o *MgcResourceRead) ReadResultSchema() *mgcSchemaPkg.Schema {
	return o.operation.ResultSchema()
}

func (o *MgcResourceRead) ChainOperations(_ context.Context, readResult core.ResultWithValue, _, _ TerraformParams) ([]MgcOperation, bool, Diagnostics) {
	readResultKeys := []tfName{}
	if readMap, ok := readResult.Value().(map[string]any); ok {
		for k := range readMap {
			attr, ok := o.attrTree.output[mgcName(k)]
			if !ok {
				continue
			}
			readResultKeys = append(readResultKeys, attr.tfName)
		}
	}
	return []MgcOperation{newMgcPopulateUnknownState(o.resourceName, readResultKeys)}, true, nil
}

var _ MgcOperation = (*MgcResourceRead)(nil)
