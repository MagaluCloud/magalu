package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
	"magalu.cloud/core/schema"
)

type MgcConnectionLink struct {
	resourceName       tfName
	attrTree           resAttrInfoTree
	getPrivateStateKey func(context.Context, string) ([]byte, diag.Diagnostics)
	setPrivateStateKey func(context.Context, string, []byte) diag.Diagnostics
	createOperation    core.Executor
	operationLoader    func(core.Result) MgcOperation
}

func newMgcConnectionLink(
	resourceName tfName,
	attrTree resAttrInfoTree,
	getPrivateStateKey func(context.Context, string) ([]byte, diag.Diagnostics),
	setPrivateStateKey func(context.Context, string, []byte) diag.Diagnostics,
	createOperation core.Executor,
	operationLoader func(core.Result) MgcOperation,
) MgcOperation {
	return &MgcConnectionLink{
		resourceName:       resourceName,
		attrTree:           attrTree,
		getPrivateStateKey: getPrivateStateKey,
		setPrivateStateKey: setPrivateStateKey,
		createOperation:    createOperation,
		operationLoader:    operationLoader,
	}
}

func (o *MgcConnectionLink) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "delete")
	ctx = tflog.SetField(ctx, connectionResourceNameField, o.resourceName)
	return ctx
}

func (o *MgcConnectionLink) CollectParameters(ctx context.Context, _, plan TerraformParams) (core.Parameters, Diagnostics) {
	return nil, nil
}

func (o *MgcConnectionLink) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return nil, nil
}

func (o *MgcConnectionLink) ShouldRun(context.Context, core.Parameters, core.Configs) (run bool, d Diagnostics) {
	return true, d
}

func (o *MgcConnectionLink) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	diagnostics := Diagnostics{}

	createResultData, d := o.getPrivateStateKey(ctx, createResultKey)
	if diagnostics.AppendCheckError(d...) {
		return nil, diagnostics.AppendErrorReturn("unable to read creation result from Terraform state", "")
	}

	tflog.Debug(ctx, "[connection-resource] about to decode creation result", map[string]any{"encoded result": string(createResultData)})
	createResult := o.createOperation.EmptyResult()
	err := createResult.Decode(createResultData)
	if err != nil {
		return nil, diagnostics.AppendErrorReturn("Failed to decode creation result", fmt.Sprintf("%v", err))
	}

	d = o.setPrivateStateKey(ctx, createResultKey, createResultData)
	if diagnostics.AppendCheckError(d...) {
		return nil, diagnostics
	}

	if createResultWithValue, ok := createResult.(core.ResultWithValue); ok {
		return createResultWithValue, diagnostics
	}

	return core.NewSimpleResult(createResult.Source(), schema.NewNullSchema(), nil), diagnostics
}

func (o *MgcConnectionLink) PostRun(_ context.Context, _ core.ResultWithValue, _, _ TerraformParams, _ *tfsdk.State) (runChain bool, diagnostics Diagnostics) {
	return true, nil
}

func (o *MgcConnectionLink) ChainOperations(ctx context.Context, createResult core.ResultWithValue, state, plan TerraformParams) ([]MgcOperation, bool, Diagnostics) {
	if operation := o.operationLoader(createResult); operation != nil {
		return []MgcOperation{operation}, true, nil
	}
	return nil, false, nil
}

var _ MgcOperation = (*MgcConnectionLink)(nil)
