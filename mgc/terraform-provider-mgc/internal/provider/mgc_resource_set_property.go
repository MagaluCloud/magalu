package provider

import (
	"context"
	"fmt"
	"maps"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
)

type MgcResourceWithPropSetterChain struct {
	resourceName     tfName
	attrTree         resAttrInfoTree
	remainingSetters map[mgcName]propertySetter
	readResource     core.Executor
}

func (o *MgcResourceWithPropSetterChain) collectChainSetPropertyOperation(
	ctx context.Context,
	readResult core.ResultWithValue,
	state, plan TerraformParams,
) (MgcOperation, bool, Diagnostics) {
	diagnostics := Diagnostics{}

	for stateValKey, currentStateVal := range state {
		attrInfo, ok := o.attrTree.getTFInputFirst(stateValKey)
		if !ok {
			continue
		}

		propertySetter, ok := o.remainingSetters[attrInfo.mgcName]
		if !ok {
			tflog.Debug(ctx, fmt.Sprintf(
				"[property-setter] property %q does not have a property setter. If the result value doesn't match the plan, inconsistent results will appear",
				stateValKey,
			))
			continue
		}

		tflog.Debug(ctx, fmt.Sprintf(
			"[property-setter] checking if property %q is already in planned state",
			stateValKey,
		))

		currentVal, _, d := loadMgcSchemaValue(ctx, attrInfo, currentStateVal, true, true)
		if diagnostics.AppendCheckError(d...) {
			return nil, false, diagnostics
		}

		plannedStateVal, ok := plan[stateValKey]
		if !ok {
			tflog.Debug(ctx,
				"[property-setter] planned state is unknown, so property setter won't be called",
				map[string]any{"propTFName": stateValKey},
			)
			continue
		}

		plannedVal, ok, d := loadMgcSchemaValue(ctx, attrInfo, plannedStateVal, true, true)
		if diagnostics.AppendCheckError(d...) {
			return nil, false, diagnostics
		}
		if !ok {
			tflog.Debug(ctx,
				"[property-setter] planned state is unknown, so property setter won't be called",
				map[string]any{"propTFName": stateValKey},
			)
			continue
		}

		if reflect.DeepEqual(currentVal, plannedVal) {
			tflog.Debug(ctx, fmt.Sprintf(
				"[property-setter] property %q is already in desired state, skipping: %#v",
				stateValKey,
				plannedVal,
			))
			continue
		}

		tflog.Debug(ctx,
			fmt.Sprintf("[property-setter] will chain property update for %q", stateValKey),
			map[string]any{"currentVal": currentVal, "plannedVal": plannedVal},
		)

		remainingSetters := maps.Clone(o.remainingSetters)
		delete(remainingSetters, attrInfo.mgcName)

		operation, d := newMgcResourceSetProperty(
			o.resourceName,
			o.attrTree,
			readResult,
			propertySetter,
			currentVal, plannedVal,
			o.readResource,
			remainingSetters,
		)
		if diagnostics.AppendCheckError(d...) {
			return nil, false, diagnostics
		}

		return operation, true, diagnostics
	}
	return nil, false, diagnostics
}

type MgcResourceSetProperty struct {
	*MgcResourceWithPropSetterChain
	setter         propertySetter
	previousResult core.ResultWithValue
	operationLink  core.Linker
}

func newMgcResourceSetProperty(
	resourceName tfName,
	attrTree resAttrInfoTree,
	previousResult core.ResultWithValue,
	setter propertySetter,
	currentVal, plannedVal any,
	readResource core.Executor,
	remainingSetters map[mgcName]propertySetter,
) (MgcOperation, Diagnostics) {
	operationLink, err := setter.getTarget(currentVal, plannedVal)
	if err != nil {
		return nil, NewLocalErrorDiagnostics(
			"unable to call property setter",
			fmt.Sprintf("unable to create executor: %v", err),
		)
	}
	return &MgcResourceSetProperty{
		MgcResourceWithPropSetterChain: &MgcResourceWithPropSetterChain{
			resourceName:     resourceName,
			attrTree:         attrTree,
			readResource:     readResource,
			remainingSetters: remainingSetters,
		},
		setter:         setter,
		previousResult: previousResult,
		operationLink:  operationLink,
	}, nil
}

func (o *MgcResourceSetProperty) WrapConext(ctx context.Context) context.Context {
	ctx = tflog.SetField(ctx, rpcField, "set property")
	ctx = tflog.SetField(ctx, resourceNameField, o.resourceName)
	return ctx
}

func (o *MgcResourceSetProperty) CollectParameters(ctx context.Context, _, plan TerraformParams) (core.Parameters, Diagnostics) {
	propAttrMaps, ok := o.attrTree.propertySetterInputs[o.setter.propertyName()]
	if !ok {
		return nil, nil
	}

	attrMap, ok := propAttrMaps[o.operationLink.AdditionalParametersSchema()]
	if !ok {
		return nil, nil
	}

	return loadMgcParamsFromState(ctx, o.operationLink.AdditionalParametersSchema(), attrMap, plan)
}

func (o *MgcResourceSetProperty) CollectConfigs(ctx context.Context, _, _ TerraformParams) (core.Configs, Diagnostics) {
	return getConfigs(ctx, o.operationLink.AdditionalConfigsSchema()), nil
}

func (o *MgcResourceSetProperty) ShouldRun(context.Context, core.Parameters, core.Configs) (run bool, d Diagnostics) {
	return true, d
}

func (o *MgcResourceSetProperty) Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics) {
	operation, err := o.operationLink.CreateExecutor(o.previousResult)
	if err != nil {
		return nil, NewErrorDiagnostics("unable o call property setter", "unable o call property setter")
	}

	return execute(ctx, o.resourceName, operation, params, configs)
}

func (o *MgcResourceSetProperty) PostRun(ctx context.Context, result core.ResultWithValue, state, plan TerraformParams, targetState *tfsdk.State) (runChain bool, diagnostics Diagnostics) {
	diagnostics = Diagnostics{}

	d := applyStateAfter(ctx, o.resourceName, o.attrTree, result, targetState)
	if diagnostics.AppendCheckError(d...) {
		return false, diagnostics
	}

	// If prop has a current counterpart (was split due to different schemas), 'applyStateAfter'
	// won't set the value of the prop in the state because it will only apply the result to the
	// output attributes (which will only have the current counterpart). Normally, applying the
	// parameters to the state into the input attributes would suffice, but since property setters
	// are defined via links, sometimes the parameter will be built into the executor, and it won't
	// be an "AdditionalParameter", so it won't be applied and we need to apply it manually
	if prop, ok := o.attrTree.input[o.setter.propertyName()]; ok && prop.currentCounterpart != nil {
		propVal, isKnown, d := loadMgcSchemaValue(ctx, prop, plan[prop.tfName], false, false)
		if diagnostics.AppendCheckError(d...) {
			return false, diagnostics
		}

		if isKnown {
			d := applyValueToState(ctx, propVal, prop, targetState, path.Empty().AtName(string(o.setter.propertyName())))
			if diagnostics.AppendCheckError(d...) {
				return false, diagnostics
			}
		}
	}

	return true, diagnostics
}

func (o *MgcResourceWithPropSetterChain) ChainOperations(ctx context.Context, propSetResult core.ResultWithValue, state, plan TerraformParams) ([]MgcOperation, bool, Diagnostics) {
	readAfter := newMgcResourceReadAfterWithLinkOrFallback(
		o.resourceName,
		o.attrTree,
		propSetResult,
		"get",
		o.readResource,
		func(ctx context.Context, readResult core.ResultWithValue, state, plan TerraformParams) ([]MgcOperation, bool, Diagnostics) {
			diagnostics := Diagnostics{}

			setPropOp, ok, d := o.collectChainSetPropertyOperation(ctx, readResult, state, plan)
			if diagnostics.AppendCheckError(d...) || !ok {
				return nil, false, diagnostics
			}

			return []MgcOperation{setPropOp}, true, diagnostics
		},
	)

	return []MgcOperation{readAfter}, true, nil
}

var _ MgcOperation = (*MgcResourceSetProperty)(nil)
