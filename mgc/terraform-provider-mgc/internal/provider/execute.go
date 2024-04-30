package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
	"magalu.cloud/core/http"
)

func validateResult(result core.ResultWithValue) Diagnostics {
	err := result.ValidateSchema()
	if err != nil {
		// TODO: Return errors instead of warnings
		return NewWarningDiagnostics(
			"Operation output mismatch",
			fmt.Sprintf("Result has invalid structure: %v", err),
		)
	}
	return nil
}

func executeTerminator(
	ctx context.Context,
	resName tfName,
	exec core.TerminatorExecutor,
	params core.Parameters,
	configs core.Configs,
) (result core.Result, err error) {
	for i := 0; i < 20; i++ {
		result, err = exec.ExecuteUntilTermination(ctx, params, configs)
		if err == nil {
			return
		}

		tflog.Debug(
			ctx,
			"[resource] operation returned error, checking if it's HTTP error",
			map[string]any{"err": err},
		)
		httpError := new(http.HttpError)
		ok := errors.As(err, &httpError)
		if !ok {
			tflog.Debug(
				ctx,
				"[resource] error is not HTTP error, continuing normal flow",
				map[string]any{"err": err},
			)
			return
		}

		if httpError.Code < 500 || httpError.Code > 599 {
			tflog.Debug(
				ctx,
				"[resource] error is HTTP error but not server error, continuing normal flow",
				map[string]any{"err": err},
			)
			return
		}

		// We ignore internal server errors because some Resources take a long time to be
		// created, and when we poll them with more than 50 requests and 1 of them fails
		// due to a server instability, we don't want to fail completely
		tflog.Debug(
			ctx,
			"[resource] error is HTTP server error. Retrying operation.",
			map[string]any{"err": err, "maxRetries": 20, "currentIteration": i},
		)
	}
	return
}

func execute(
	ctx context.Context,
	resName tfName,
	exec core.Executor,
	params core.Parameters,
	configs core.Configs,
) (core.ResultWithValue, Diagnostics) {
	var diagnostics = Diagnostics{}
	var result core.Result
	var err error

	tflog.Debug(ctx, fmt.Sprintf("[resource] will %s new %s resource - request info with params: %#v and configs: %#v", exec.Name(), resName, params, configs))
	if tExec, ok := core.ExecutorAs[core.TerminatorExecutor](exec); ok {
		tflog.Debug(ctx, "[resource] running as TerminatorExecutor")
		result, err = executeTerminator(ctx, resName, tExec, params, configs)
	} else {
		tflog.Debug(ctx, "[resource] running as Executor")
		result, err = exec.Execute(ctx, params, configs)
	}
	if err != nil {
		return nil, diagnostics.AppendErrorReturn(
			fmt.Sprintf("Unable to %s %s", exec.Name(), resName),
			fmt.Sprintf("Service returned with error: %v", err),
		)
	}

	resultWithValue, ok := core.ResultAs[core.ResultWithValue](result)
	if !ok {
		if resultSchema := exec.ResultSchema(); resultSchema.Nullable || resultSchema.IsEmpty() {
			resultWithValue = core.NewSimpleResult(result.Source(), exec.ResultSchema(), nil)
		} else {
			// Should this really be an error? Don't really know. Why not let 'validateResult' handle this?
			// This would probably further state updates so it's probably better NOT to error here
			return nil, diagnostics.AppendErrorReturn(
				"Operation output mismatch",
				fmt.Sprintf("result has no value %#v", result),
			)
		}
	}

	return resultWithValue, diagnostics
}
