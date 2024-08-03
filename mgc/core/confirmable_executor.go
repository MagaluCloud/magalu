package core

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"magalu.cloud/core/utils"
)

const defaultPromptMessage = "This action cannot be undone. Proceed?"

type UserDeniedConfirmationError struct {
	Prompt string
}

func (e UserDeniedConfirmationError) Error() string {
	return fmt.Sprintf("%s: user denied", e.Prompt)
}

type ConfirmableExecutor interface {
	Executor
	ConfirmPrompt(parameters Parameters, configs Configs) (message string)
}

func NewConfirmableExecutor(
	exec Executor,
	confirmPrompt func(parameters Parameters, configs Configs) (message string),
) ConfirmableExecutor {
	return &confirmableExecutor{exec, confirmPrompt}
}

type confirmableExecutor struct {
	Executor
	confirmPrompt func(parameters Parameters, configs Configs) (message string)
}

func (o *confirmableExecutor) Execute(ctx context.Context, parameters Parameters, configs Configs) (result Result, err error) {
	result, err = o.Executor.Execute(ctx, parameters, configs)
	return ExecutorWrapResult(o, result, err)
}

func (o *confirmableExecutor) Unwrap() Executor {
	return o.Executor
}

func (o *confirmableExecutor) ConfirmPrompt(parameters Parameters, configs Configs) (message string) {
	return o.confirmPrompt(parameters, configs)
}

func DefaultConfirmPrompt(parameters Parameters, configs Configs) string {
	return defaultPromptMessage
}

// ConfirmPromptWithTemplate parses and executes the template provided in the message.
// If the template fails to parse or execute, a default generic message is returned
func ConfirmPromptWithTemplate(msg string) func(parameters Parameters, configs Configs) string {
	if msg == "" {
		return DefaultConfirmPrompt
	}

	tmpl, err := utils.NewTemplate(msg)
	if err != nil {
		return DefaultConfirmPrompt
	}

	return func(parameters Parameters, configs Configs) string {
		value := map[string]any{"parameters": parameters, "configs": configs}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, value)
		if err != nil {
			return defaultPromptMessage
		}
		s := buf.String()
		s = strings.Trim(s, " \t\n\r")
		return s
	}
}

var _ Executor = (*confirmableExecutor)(nil)
var _ ExecutorWrapper = (*confirmableExecutor)(nil)
var _ ConfirmableExecutor = (*confirmableExecutor)(nil)

var _ error = UserDeniedConfirmationError{}
