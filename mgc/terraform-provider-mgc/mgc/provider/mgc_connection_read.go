package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
)

type MgcConnectionRead struct {
	resourceName   tfName
	attrTree       resAttrInfoTree
	readConnection core.Executor
}

func newMgcConnectionRead(
	resourceName tfName,
	attrTree resAttrInfoTree,
	readConnection core.Executor,
) MgcOperation {
	return &MgcConnectionRead{
		resourceName:   resourceName,
		attrTree:       attrTree,
		readConnection: readConnection,
	}
}

func (o *MgcConnectionRead) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "read")
	ctx = tflog.SetField(ctx, connectionResourceNameField, o.resourceName)
	return ctx
}

func (o *MgcConnectionRead) CollectParameters(ctx context.Context, _, plan TerraformParams) (core.Parameters, Diagnostics) {
	return loadMgcParamsFromState(ctx, o.readConnection.ParametersSchema(), o.attrTree.input, plan)
}

func (o *MgcConnectionRead) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return getConfigs(ctx, o.readConnection.ConfigsSchema()), nil
}

func (o *MgcConnectionRead) ShouldRun(context.Context, core.Parameters, core.Configs) (run bool, d Diagnostics) {
	return true, d
}

func (o *MgcConnectionRead) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	return execute(ctx, o.resourceName, o.readConnection, params, configs)
}

func (o *MgcConnectionRead) PostRun(ctx context.Context, readResult core.ResultWithValue, state, plan TerraformParams, targetState *tfsdk.State) (runChain bool, diagnostics Diagnostics) {
	tflog.Info(ctx, "connection read")
	diagnostics = Diagnostics{}

	d := applyStateAfter(ctx, o.resourceName, o.attrTree, readResult, state, targetState)
	if diagnostics.AppendCheckError(d...) {
		return false, diagnostics
	}

	return true, diagnostics
}

func (o *MgcConnectionRead) ChainOperations(_ context.Context, _ core.ResultWithValue, _, _ TerraformParams) ([]MgcOperation, bool, Diagnostics) {
	return nil, false, nil
}

var _ MgcOperation = (*MgcConnectionRead)(nil)
