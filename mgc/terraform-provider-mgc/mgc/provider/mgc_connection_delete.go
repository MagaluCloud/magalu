package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
)

type MgcConnectionDelete struct {
	resourceName     tfName
	attrTree         resAttrInfoTree
	deleteConnection core.Executor
}

func newMgcConnectionDelete(
	resourceName tfName,
	attrTree resAttrInfoTree,
	deleteConnection core.Executor,
) MgcOperation {
	return &MgcConnectionDelete{
		resourceName:     resourceName,
		attrTree:         attrTree,
		deleteConnection: deleteConnection,
	}
}

func (o *MgcConnectionDelete) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "delete")
	ctx = tflog.SetField(ctx, connectionResourceNameField, o.resourceName)
	return ctx
}

func (o *MgcConnectionDelete) CollectParameters(ctx context.Context, _, plan TerraformParams) (core.Parameters, Diagnostics) {
	return loadMgcParamsFromState(ctx, o.deleteConnection.ParametersSchema(), o.attrTree.deleteInput, plan)
}

func (o *MgcConnectionDelete) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return getConfigs(ctx, o.deleteConnection.ConfigsSchema()), nil
}

func (o *MgcConnectionDelete) ShouldRun(context.Context, core.Parameters, core.Configs) (run bool, d Diagnostics) {
	return true, d
}

func (o *MgcConnectionDelete) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	return execute(ctx, o.resourceName, o.deleteConnection, params, configs)
}

func (o *MgcConnectionDelete) PostRun(ctx context.Context, _ core.ResultWithValue, _, _ TerraformParams, _ *tfsdk.State) (runChain bool, diagnostics Diagnostics) {
	tflog.Info(ctx, "connection deleted")
	return false, nil
}

func (o *MgcConnectionDelete) ChainOperations(_ context.Context, _ core.ResultWithValue, _, _ TerraformParams) ([]MgcOperation, bool, Diagnostics) {
	return nil, false, nil
}

var _ MgcOperation = (*MgcConnectionDelete)(nil)
