package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
)

type MgcResourceCreate struct {
	*MgcResourceWithPropSetterChain
	createResource core.Executor
}

func newMgcResourceCreate(
	resourceName tfName,
	attrTree resAttrInfoTree,
	createResource core.Executor,
	readResource core.Executor,
	propertySetters map[mgcName]propertySetter,
) MgcOperation {
	return &MgcResourceCreate{
		MgcResourceWithPropSetterChain: &MgcResourceWithPropSetterChain{
			resourceName:     resourceName,
			attrTree:         attrTree,
			readResource:     readResource,
			remainingSetters: propertySetters,
		},
		createResource: createResource,
	}
}

func (o *MgcResourceCreate) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "create")
	ctx = tflog.SetField(ctx, resourceNameField, o.resourceName)
	return ctx
}

func (o *MgcResourceCreate) CollectParameters(ctx context.Context, _, plan TerraformParams) (core.Parameters, Diagnostics) {
	return loadMgcParamsFromState(ctx, o.createResource.ParametersSchema(), o.attrTree.createInput, plan)
}

func (o *MgcResourceCreate) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return getConfigs(ctx, o.createResource.ConfigsSchema()), nil
}

func (o *MgcResourceCreate) ShouldRun(context.Context, core.Parameters, core.Configs) (run bool, d Diagnostics) {
	return true, d
}

func (o *MgcResourceCreate) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	return execute(ctx, o.resourceName, o.createResource, params, configs)
}

func (o *MgcResourceCreate) PostRun(ctx context.Context, createResult core.ResultWithValue, state, plan TerraformParams, targetState *tfsdk.State) (runChain bool, diagnostics Diagnostics) {
	tflog.Info(ctx, "resource created")
	diagnostics = Diagnostics{}

	d := applyStateAfter(ctx, o.resourceName, o.attrTree, createResult, targetState)
	if diagnostics.AppendCheckError(d...) {
		return false, diagnostics
	}

	return true, diagnostics
}

var _ MgcOperation = (*MgcResourceCreate)(nil)
