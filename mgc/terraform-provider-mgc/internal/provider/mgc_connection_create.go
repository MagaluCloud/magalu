package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
)

type MgcConnectionCreate struct {
	resourceName       tfName
	attrTree           resAttrInfoTree
	getPrivateStateKey func(context.Context, string) ([]byte, diag.Diagnostics)
	setPrivateStateKey func(context.Context, string, []byte) diag.Diagnostics
	createConnection   core.Executor
	deleteConnection   core.Linker
}

func newMgcConnectionCreate(
	resourceName tfName,
	attrTree resAttrInfoTree,
	getPrivateStateKey func(context.Context, string) ([]byte, diag.Diagnostics),
	setPrivateStateKey func(context.Context, string, []byte) diag.Diagnostics,
	createConnection core.Executor,
	deleteConnection core.Linker,
) MgcOperation {
	return &MgcConnectionCreate{
		resourceName:       resourceName,
		attrTree:           attrTree,
		getPrivateStateKey: getPrivateStateKey,
		setPrivateStateKey: setPrivateStateKey,
		deleteConnection:   deleteConnection,
		createConnection:   createConnection,
	}
}

func (o *MgcConnectionCreate) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "create")
	ctx = tflog.SetField(ctx, connectionResourceNameField, o.resourceName)
	return ctx
}

func (o *MgcConnectionCreate) CollectParameters(ctx context.Context, _, plan TerraformParams) (core.Parameters, Diagnostics) {
	return loadMgcParamsFromState(ctx, o.createConnection.ParametersSchema(), o.attrTree.createInput, plan)
}

func (o *MgcConnectionCreate) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return getConfigs(ctx, o.createConnection.ConfigsSchema()), nil
}

func (o *MgcConnectionCreate) ShouldRun(context.Context, core.Parameters, core.Configs) (run bool, d Diagnostics) {
	return true, d
}

func (o *MgcConnectionCreate) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	return execute(ctx, o.resourceName, o.createConnection, params, configs)
}

func (o *MgcConnectionCreate) PostRun(ctx context.Context, createResult core.ResultWithValue, state, plan TerraformParams, targetState *tfsdk.State) (runChain bool, diagnostics Diagnostics) {
	tflog.Info(ctx, "connection created")
	diagnostics = Diagnostics{}

	createResultEncoded, err := createResult.Encode()
	if err != nil {
		return true, diagnostics.AppendErrorReturn(
			"failure to encode connection resource creation result",
			"Terraform wasn't able to encode the result of the creation process to save in its state. Creation was successful, but resource will be deleted, try again.",
		)
	}

	tflog.Debug(ctx, "about to store private creation result", map[string]any{"encoded result": createResultEncoded})

	diag := o.setPrivateStateKey(ctx, createResultKey, createResultEncoded)
	if diagnostics.AppendCheckError(diag...) {
		return true, diagnostics
	}

	d := applyStateAfter(ctx, o.resourceName, o.attrTree, createResult, state, targetState)
	if diagnostics.AppendCheckError(d...) {
		return true, diagnostics
	}

	return true, diagnostics
}

func (o *MgcConnectionCreate) ChainOperations(ctx context.Context, createResult core.ResultWithValue, state, plan TerraformParams) ([]MgcOperation, bool, Diagnostics) {
	diagnostics := Diagnostics{}

	createResultData, d := o.getPrivateStateKey(ctx, createResultKey)
	if diagnostics.AppendCheckError(d...) || len(createResultData) == 0 {
		deleteExec, err := o.deleteConnection.CreateExecutor(createResult)
		if err != nil {
			return nil, false, diagnostics.AppendErrorReturn(
				"unable to delete broken connection",
				"unable to delete broken connection, server state will be out of sync with Terraform",
			)
		}
		deleteOperation := newMgcConnectionDelete(
			o.resourceName,
			o.attrTree,
			deleteExec,
		)
		return []MgcOperation{deleteOperation}, true, diagnostics
	}

	return nil, false, diagnostics
}

var _ MgcOperation = (*MgcConnectionCreate)(nil)
