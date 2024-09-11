package provider

import (
	"context"
	"reflect"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
)

type MgcResourceUpdate struct {
	*MgcResourceWithPropSetterChain
	updateResource core.Executor
}

func newMgcResourceUpdate(
	resourceName tfName,
	attrTree resAttrInfoTree,
	updateResource core.Executor,
	readResource core.Executor,
	propertySetters map[mgcName]propertySetter,
) MgcOperation {
	return &MgcResourceUpdate{
		MgcResourceWithPropSetterChain: &MgcResourceWithPropSetterChain{
			resourceName:     resourceName,
			attrTree:         attrTree,
			readResource:     readResource,
			remainingSetters: propertySetters,
		},
		updateResource: updateResource,
	}
}

func (o *MgcResourceUpdate) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "update")
	ctx = tflog.SetField(ctx, resourceNameField, o.resourceName)
	return ctx
}

func (o *MgcResourceUpdate) CollectParameters(ctx context.Context, state, plan TerraformParams) (core.Parameters, Diagnostics) {
	diagnostics := Diagnostics{}
	params := map[string]any{}
	paramsSchema := o.updateResource.ParametersSchema()
	readParamsSchema := o.readResource.ParametersSchema()

	plannedParams, d := loadMgcParamsFromState(ctx, paramsSchema, o.attrTree.updateInput, plan)
	if diagnostics.AppendCheckError(d...) {
		return nil, d
	}

	stateParams, d := loadMgcParamsFromState(ctx, paramsSchema, o.attrTree.updateInput, state)
	if diagnostics.AppendCheckError(d...) {
		return nil, d
	}

	for prop := range paramsSchema.Properties {
		if _, ok := plannedParams[prop]; !ok {
			continue
		}

		// Prioritize property setters in case the property can be set both via Update and via Property Setters
		if _, ok := o.remainingSetters[mgcName(prop)]; ok {
			continue
		}

		if _, ok := readParamsSchema.Properties[prop]; ok || !reflect.DeepEqual(plannedParams[prop], stateParams[prop]) {
			params[prop] = plannedParams[prop]
			continue
		}

		if !slices.Contains(paramsSchema.Required, prop) {
			continue
		}

		tflog.Warn(
			ctx,
			"update request has unchanged parameters values from current state. This may cause an error",
			map[string]any{
				"parameter": prop,
				"value":     plannedParams[prop],
			},
		)

		params[prop] = plannedParams[prop]
	}
	return params, diagnostics
}

func (o *MgcResourceUpdate) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return getConfigs(ctx, o.updateResource.ConfigsSchema()), nil
}

func (o *MgcResourceUpdate) ShouldRun(_ context.Context, params core.Parameters, _ core.Configs) (run bool, d Diagnostics) {
	readParamsSchema := o.readResource.ParametersSchema()
	for paramName := range params {
		if _, ok := readParamsSchema.Properties[paramName]; !ok {
			return true, nil
		}
	}
	// Don't perform update request if there are only "read" parameters
	return false, d
}

func (o *MgcResourceUpdate) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	return execute(ctx, o.resourceName, o.updateResource, params, configs)
}

func (o *MgcResourceUpdate) PostRun(ctx context.Context, updateResult core.ResultWithValue, state, plan TerraformParams, targetState *tfsdk.State) (runChain bool, diagnostics Diagnostics) {
	tflog.Info(ctx, "resource updated")
	diagnostics = Diagnostics{}

	d := applyStateAfter(ctx, o.resourceName, o.attrTree, updateResult, state, targetState)
	if diagnostics.AppendCheckError(d...) {
		return false, diagnostics
	}

	return true, diagnostics
}

var _ MgcOperation = (*MgcResourceUpdate)(nil)
