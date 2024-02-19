package core

import (
	"context"
	"fmt"

	"magalu.cloud/core/utils"
)

type PromptInput func(parameters Parameters, configs Configs) (message string, validate func(input string) error)

type PromptInputExecutor interface {
	Executor
	PromptInput(parameters Parameters, configs Configs) (message string, validate func(input string) error)
}

type promptInputExecutor struct {
	Executor
	promptInput PromptInput
}

func (o *promptInputExecutor) PromptInput(parameters Parameters, configs Configs) (message string, validate func(input string) error) {
	return o.promptInput(parameters, configs)
}

func (o *promptInputExecutor) Unwrap() Executor {
	return o.Executor
}

func (o *promptInputExecutor) Execute(ctx context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	result, err = o.Executor.Execute(ctx, parameters, configs)
	return ExecutorWrapResult(o, result, err)
}

func NewPromptInputExecutor(
	exec Executor,
	promptInput PromptInput,
) PromptInputExecutor {
	return &promptInputExecutor{
		exec,
		promptInput,
	}
}

// The msgTemplate represents a template used to generate the confirmation message for the user
// The confirmationValueTemplate define what the user need to type in order to confirm the operation
func NewPromptInput(msgTemplate string, confirmationTemplate string) PromptInput {
	if msgTemplate == "" {
		msgTemplate = "Please type \"{{ .confirmationValue }}\" to confirm"
	}

	messageTemplate, err := utils.NewTemplate(msgTemplate)
	if err != nil {
		panic(fmt.Sprintf("\"messageTemplate\": %s", err.Error()))
	}

	confirmationValueTemplate, err := utils.NewTemplate(confirmationTemplate)
	if err != nil {
		panic(fmt.Sprintf("\"confirmationValueTemplate\": %s", err.Error()))
	}

	return func(parameters Parameters, configs Configs) (string, func(input string) error) {
		value := map[string]any{"parameters": parameters, "configs": configs}

		confirmValue, err := utils.ExecuteTemplateTrimmed(confirmationValueTemplate, value)
		if err != nil {
			panic(fmt.Errorf("could not render \"confirmationValueTemplate\": %w", err))
		}
		value["confirmationValue"] = confirmValue

		msg, err := utils.ExecuteTemplateTrimmed(messageTemplate, value)
		if err != nil {
			panic(fmt.Errorf("could not render \"messageTemplate\": %w", err))
		}

		validate := func(input string) error {
			if input != confirmValue {
				return fmt.Errorf("aborted, input didn't match %q", confirmValue)
			}
			return nil
		}
		return msg, validate
	}
}

var _ Executor = (*promptInputExecutor)(nil)
var _ ExecutorWrapper = (*promptInputExecutor)(nil)
var _ PromptInputExecutor = (*promptInputExecutor)(nil)
