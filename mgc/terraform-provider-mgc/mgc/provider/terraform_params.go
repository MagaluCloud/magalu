package provider

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
)

type TerraformParam tftypes.Value
type TerraformParams map[tfName]tftypes.Value

func tfStateToParams(state tfsdk.State) (TerraformParams, error) {
	var m map[string]tftypes.Value
	err := state.Raw.As(&m)
	if err != nil {
		return nil, err
	}

	result := make(TerraformParams, len(m))
	for k, v := range m {
		result[tfName(k)] = v
	}

	return result, nil
}

func loadMgcParamsFromState(
	ctx context.Context,
	paramsSchema *mgcSchemaPkg.Schema,
	attrMap resAttrInfoMap,
	state TerraformParams,
) (core.Parameters, Diagnostics) {
	params := core.Parameters{}
	diagnostics := Diagnostics{}

	keys := make([]string, 0, len(paramsSchema.Properties))
	for paramName := range paramsSchema.Properties {
		keys = append(keys, paramName)
	}

	tflog.Debug(
		ctx,
		"[loader] loading parameters from parameters schema",
		map[string]any{"parametersToLoad": keys},
	)

	for paramName := range paramsSchema.Properties {
		required := slices.Contains(paramsSchema.Required, paramName)

		attr, ok := attrMap[mgcName(paramName)]
		if !ok {
			if !required {
				continue
			}

			return params, diagnostics.AppendLocalErrorReturn(
				fmt.Sprintf("[loader] Tried to load parameter %q from state, but couldn't", paramName),
				fmt.Sprintf(
					"Tried to load parameter %q from state, but couldn't, as there was no attribute in the schema that matches it",
					paramName,
				),
			)
		}

		tfStateVal, ok := state[attr.tfName]
		if !ok {
			if !required {
				continue
			}

			return params, diagnostics.AppendLocalErrorReturn(
				fmt.Sprintf("[loader] Tried to load parameter %q from state, but couldn't", paramName),
				fmt.Sprintf(
					"Tried to load parameter %q from state, but couldn't, as there was no value in the state. This probably means that a default value isn't being sent by Terraform",
					paramName,
				),
			)
		}

		param, ok, d := loadMgcSchemaValue(ctx, attr, tfStateVal, true, !required)
		if diagnostics.AppendCheckError(d...) {
			return params, diagnostics
		}

		if !ok {
			if required {
				return params, diagnostics.AppendErrorReturn(
					fmt.Sprintf("[loader] Tried to load parameter %v from state, but couldn't", paramName),
					"Tried to load parameter %s from state, but couldn't, as there was no value in the state. This probably means that a default value isn't being sent by Terraform",
				)
			}
			continue
		}

		params[paramName] = param
	}
	return params, diagnostics

}
