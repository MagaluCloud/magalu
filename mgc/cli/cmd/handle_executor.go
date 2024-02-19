package cmd

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	"magalu.cloud/cli/ui"
	"magalu.cloud/core"
	"magalu.cloud/core/auth"
	"magalu.cloud/core/progress_report"
	mgcSdk "magalu.cloud/sdk"
)

func handleExecutorResult(ctx context.Context, sdk *mgcSdk.Sdk, cmd *cobra.Command, result core.Result, err error) error {
	if err != nil {
		var failedTerminationError core.FailedTerminationError
		if errors.As(err, &failedTerminationError) {
			_ = formatResult(sdk, cmd, failedTerminationError.Result)
		}
		return err
	}

	return formatResult(sdk, cmd, result)
}

func checkScopes(sdk *mgcSdk.Sdk, exec core.Executor) error {
	a := sdk.Auth()
	if a == nil {
		return fmt.Errorf("programming error: context did not contain SDK Auth information")
	}

	currentScopes, err := a.CurrentScopes()
	if err != nil {
		return fmt.Errorf("unable to get current scopes: %w", err)
	}

	var missing []string
	necessaryScopes := exec.Scopes()
	for _, scope := range necessaryScopes {
		if !slices.Contains(currentScopes, auth.Scope(scope)) {
			missing = append(missing, scope)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("this operation requires the following scopes: %v. To add them, run 'auth scopes add'", missing)
	}

	return nil
}

func handleExecutor(
	ctx context.Context,
	sdk *mgcSdk.Sdk,
	cmd *cobra.Command,
	exec core.Executor,
	parameters core.Parameters,
	configs core.Configs,
) (core.Result, error) {
	if err := checkScopes(sdk, exec); err != nil {
		return nil, err
	}

	if err := exec.ParametersSchema().VisitJSON(parameters); err != nil {
		return nil, core.UsageError{Err: err}
	}

	if err := exec.ConfigsSchema().VisitJSON(configs); err != nil {
		return nil, core.UsageError{Err: err}
	}

	if pb != nil {
		ctx = progress_report.NewContext(ctx, pb.ReportProgress)
	}

	if cExec, ok := core.ExecutorAs[core.ConfirmableExecutor](exec); ok && !getBypassConfirmationFlag(cmd) {
		msg := cExec.ConfirmPrompt(parameters, configs)
		run, err := ui.Confirm(msg)
		if err != nil {
			return nil, err
		}

		if !run {
			return nil, core.UserDeniedConfirmationError{Prompt: msg}
		}
	}
	if pExec, ok := core.ExecutorAs[core.PromptInputExecutor](exec); ok && !getBypassConfirmationFlag(cmd) {
		msg, validate := pExec.PromptInput(parameters, configs)

		input, err := ui.RunPromptInput(msg)
		if err != nil {
			return nil, err
		}

		err = validate(input)
		if err != nil {
			return nil, err
		}
	}

	if t := getTimeoutFlag(cmd); t > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, t)
		defer cancel()
	}

	waitTermination := getWaitTerminationFlag(cmd)
	var cb core.RetryUntilCb
	if tExec, ok := core.ExecutorAs[core.TerminatorExecutor](exec); ok && waitTermination {
		cb = func() (result core.Result, err error) {
			return tExec.ExecuteUntilTermination(ctx, parameters, configs)
		}
	} else {
		cb = func() (result core.Result, err error) {
			return exec.Execute(ctx, parameters, configs)
		}
	}

	retry, err := getRetryUntilFlag(cmd)
	if err != nil {
		return nil, err
	}

	result, err := retry.Run(ctx, cb)

	if pb != nil {
		pb.Flush()
	}

	err = handleExecutorResult(ctx, sdk, cmd, result, err)
	if err != nil {
		return nil, err
	}

	return result, err
}
