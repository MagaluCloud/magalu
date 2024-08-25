package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"magalu.cloud/core"
	"magalu.cloud/sdk"
)

type MgcOperation interface {
	WrapConext(ctx context.Context) context.Context

	CollectParameters(ctx context.Context, state, plan TerraformParams) (core.Parameters, Diagnostics)
	CollectConfigs(ctx context.Context, state, plan TerraformParams) (core.Configs, Diagnostics)

	// If 'ShouldRun' returns false, 'Run' and 'PostRun' will be skipped. The chain operations will still be called, with
	// 'opResult' passed as nil
	ShouldRun(ctx context.Context, params core.Parameters, configs core.Configs) (run bool, d Diagnostics)
	Run(ctx context.Context, params core.Parameters, configs core.Configs) (core.ResultWithValue, Diagnostics)
	PostRun(ctx context.Context, result core.ResultWithValue, state, plan TerraformParams, targetState *tfsdk.State) (runChain bool, d Diagnostics)

	ChainOperations(ctx context.Context, opResult core.ResultWithValue, state, plan TerraformParams) ([]MgcOperation, bool, Diagnostics)
}

type MgcOperationRunner struct {
	sdk           *sdk.Sdk
	state         tfsdk.State
	plan          tfsdk.Plan
	targetState   *tfsdk.State
	rootOperation MgcOperation
}

func newMgcOperationRunner(
	sdk *sdk.Sdk,
	rootOperation MgcOperation,
	state tfsdk.State,
	plan tfsdk.Plan,
	targetState *tfsdk.State,
) MgcOperationRunner {
	return MgcOperationRunner{
		sdk:           sdk,
		rootOperation: rootOperation,
		state:         state,
		plan:          plan,
		targetState:   targetState,
	}
}

func (r *MgcOperationRunner) Run(ctx context.Context) Diagnostics {
	ctx = r.sdk.WrapContext(ctx)
	diagnostics := Diagnostics{}

	opResult, runChain, d := r.runOperation(ctx, r.rootOperation)
	diagnostics.Append(d...)

	if !runChain {
		return diagnostics
	}

	state, plan, d := r.getCurrentTFParams()
	if diagnostics.AppendCheckError(d...) {
		return diagnostics
	}

	chained, runChain, d := r.rootOperation.ChainOperations(ctx, opResult, state, plan)
	diagnostics.Append(d...)

	if !runChain {
		return diagnostics
	}

	d = r.runChainedOperations(ctx, chained)
	diagnostics.Append(d...)

	return diagnostics
}

func (r *MgcOperationRunner) getCurrentTFParams() (state, plan TerraformParams, d Diagnostics) {
	diagnostics := Diagnostics{}

	state, err := tfStateToParams(r.state)
	if err != nil {
		diagnostics.AddLocalError(
			"error when reading Terraform state",
			fmt.Sprintf("Terraform state wasn't able to be read: %s", err.Error()),
		)
	}

	plan, err = tfStateToParams(tfsdk.State(r.plan))
	if err != nil {
		diagnostics.AddLocalError(
			"error when reading Terraform state",
			fmt.Sprintf("Terraform plan wasn't able to be read: %s", err.Error()),
		)
	}

	return state, plan, diagnostics
}

func (r *MgcOperationRunner) runOperation(
	ctx context.Context,
	operation MgcOperation,
) (runResult core.ResultWithValue, runChain bool, diagnostics Diagnostics) {
	ctx = operation.WrapConext(ctx)
	diagnostics = Diagnostics{}

	state, plan, d := r.getCurrentTFParams()
	if diagnostics.AppendCheckError(d...) {
		return nil, false, diagnostics
	}

	params, d := operation.CollectParameters(ctx, state, plan)
	if diagnostics.AppendCheckError(d...) {
		return nil, false, diagnostics
	}

	configs, d := operation.CollectConfigs(ctx, state, plan)
	if diagnostics.AppendCheckError(d...) {
		return nil, false, diagnostics
	}

	shouldRun, d := operation.ShouldRun(ctx, params, configs)
	if diagnostics.AppendCheckError(d...) || !shouldRun {
		return nil, true, diagnostics
	}

	runResult, d = operation.Run(ctx, params, configs)
	if diagnostics.AppendCheckError(d...) {
		return runResult, false, diagnostics
	}

	d = validateResult(runResult)
	if diagnostics.AppendCheckError(d...) {
		return runResult, false, diagnostics
	}

	runChain, d = operation.PostRun(ctx, runResult, state, plan, r.targetState)
	if diagnostics.AppendCheckError(d...) {
		return runResult, runChain, diagnostics
	}

	r.state = *r.targetState

	return runResult, runChain, diagnostics
}

func (r *MgcOperationRunner) runChainedOperations(ctx context.Context, chained []MgcOperation) Diagnostics {
	diagnostics := Diagnostics{}

	for _, operation := range chained {
		opResult, run, d := r.runOperation(ctx, operation)
		diagnostics.Append(d...)

		if !run {
			continue
		}

		state, plan, d := r.getCurrentTFParams()
		if diagnostics.AppendCheckError(d...) {
			return diagnostics
		}

		chained, run, d := operation.ChainOperations(ctx, opResult, state, plan)
		diagnostics.Append(d...)

		if !run {
			continue
		}

		d = r.runChainedOperations(ctx, chained)
		if diagnostics.AppendCheckError(d...) {
			return diagnostics
		}
	}

	return diagnostics
}
