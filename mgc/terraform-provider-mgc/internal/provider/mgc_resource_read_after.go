package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
)

type MgcResourceReadAfter struct {
	resourceName    tfName
	attrTree        resAttrInfoTree
	operation       core.Executor
	chainOperations func(ctx context.Context, readResult core.ResultWithValue, state, plan TerraformParams) ([]MgcOperation, bool, Diagnostics)
}

func newMgcResourceReadAfterWithLinkOrFallback(
	resourceName tfName,
	attrTree resAttrInfoTree,
	previousResult core.Result,
	readLinkName string,
	fallbackRead core.Executor,
	chainOperations func(ctx context.Context, readResult core.ResultWithValue, state, plan TerraformParams) ([]MgcOperation, bool, Diagnostics),
) *MgcResourceReadAfter {
	if readAfter, ok := newMgcResourceReadAfterWithLink(resourceName, attrTree, previousResult, readLinkName, chainOperations); ok {
		return readAfter
	}
	return newMgcResourceReadAfter(resourceName, attrTree, fallbackRead, chainOperations)
}

func newMgcResourceReadAfterWithLink(
	resourceName tfName,
	attrTree resAttrInfoTree,
	previousResult core.Result,
	readLinkName string,
	chainOperations func(ctx context.Context, readResult core.ResultWithValue, state, plan TerraformParams) ([]MgcOperation, bool, Diagnostics),
) (*MgcResourceReadAfter, bool) {
	if previousResult == nil {
		return nil, false
	}

	readLink, ok := previousResult.Source().Executor.Links()[readLinkName]
	if !ok {
		return nil, false
	}

	readOp, err := readLink.CreateExecutor(previousResult)
	if err != nil {
		return nil, false
	}

	return newMgcResourceReadAfter(resourceName, attrTree, readOp, chainOperations), true
}

func newMgcResourceReadAfter(
	resourceName tfName,
	attrTree resAttrInfoTree,
	readOperation core.Executor,
	chainOperations func(_ context.Context, _ core.ResultWithValue, _, _ TerraformParams) ([]MgcOperation, bool, Diagnostics),
) *MgcResourceReadAfter {
	return &MgcResourceReadAfter{
		resourceName:    resourceName,
		attrTree:        attrTree,
		operation:       readOperation,
		chainOperations: chainOperations,
	}
}

func (o *MgcResourceReadAfter) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "read-after")
	ctx = tflog.SetField(ctx, resourceNameField, o.resourceName)
	return ctx
}

func (o *MgcResourceReadAfter) CollectParameters(ctx context.Context, state, _ TerraformParams) (core.Parameters, Diagnostics) {
	return loadMgcParamsFromState(ctx, o.operation.ParametersSchema(), o.attrTree.input, state)
}

func (o *MgcResourceReadAfter) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return getConfigs(ctx, o.operation.ConfigsSchema()), nil
}

func (o *MgcResourceReadAfter) ShouldRun(context.Context, core.Parameters, core.Configs) (run bool, d Diagnostics) {
	return true, d
}

func (o *MgcResourceReadAfter) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	return execute(ctx, o.resourceName, o.operation, params, configs)
}

func (o *MgcResourceReadAfter) PostRun(ctx context.Context, result core.ResultWithValue, state, plan TerraformParams, targetState *tfsdk.State) (runChain bool, diagnostics Diagnostics) {
	tflog.Info(ctx, "resource read after operation")
	diagnostics = Diagnostics{}

	d := applyStateAfter(ctx, o.resourceName, o.attrTree, result, targetState)
	if diagnostics.AppendCheckError(d...) {
		return false, diagnostics
	}

	return true, diagnostics
}

func (o *MgcResourceReadAfter) ChainOperations(ctx context.Context, readResult core.ResultWithValue, state, plan TerraformParams) ([]MgcOperation, bool, Diagnostics) {
	return o.chainOperations(ctx, readResult, state, plan)
}

var _ MgcOperation = (*MgcResourceReadAfter)(nil)
