package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
)

type MgcResourceDelete struct {
	resourceName   tfName
	attrTree       resAttrInfoTree
	deleteResource core.Executor
}

func newMgcResourceDelete(
	resourceName tfName,
	attrTree resAttrInfoTree,
	deleteResource core.Executor,
) MgcOperation {
	return &MgcResourceDelete{
		resourceName:   resourceName,
		attrTree:       attrTree,
		deleteResource: deleteResource,
	}
}

func (o *MgcResourceDelete) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "delete")
	ctx = tflog.SetField(ctx, resourceNameField, o.resourceName)
	return ctx
}

func (o *MgcResourceDelete) CollectParameters(ctx context.Context, _, plan TerraformParams) (core.Parameters, Diagnostics) {
	return loadMgcParamsFromState(ctx, o.deleteResource.ParametersSchema(), o.attrTree.deleteInput, plan)
}

func (o *MgcResourceDelete) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return getConfigs(ctx, o.deleteResource.ConfigsSchema()), nil
}

func (o *MgcResourceDelete) ShouldRun(context.Context, core.Parameters, core.Configs) (run bool, d Diagnostics) {
	return true, d
}

func (o *MgcResourceDelete) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	result, d := execute(ctx, o.resourceName, o.deleteResource, params, configs)
	return result, d
}

func (o *MgcResourceDelete) PostRun(ctx context.Context, result core.ResultWithValue, state, plan TerraformParams, targetState *tfsdk.State) (runChain bool, diagnostics Diagnostics) {
	tflog.Info(ctx, "resource deleted")
	return false, nil
}

func (o *MgcResourceDelete) ChainOperations(_ context.Context, _ core.ResultWithValue, _, _ TerraformParams) ([]MgcOperation, bool, Diagnostics) {
	return nil, false, nil
}

var _ MgcOperation = (*MgcResourceDelete)(nil)
