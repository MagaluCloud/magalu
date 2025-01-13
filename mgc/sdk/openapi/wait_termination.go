package openapi

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

func getDocument(result core.ResultWithValue) any {
	return map[string]any{
		"result":     result.Value(),
		"parameters": result.Source().Parameters,
		"configs":    result.Source().Configs,
	}
}

func addDocumentOwner(ownerResult core.ResultWithValue, doc any) any {
	data := doc.(map[string]any)
	data["owner"] = getDocument(ownerResult)
	return doc
}

func createGetDocument(ownerResult core.Result) func(result core.ResultWithValue) any {
	if ownerResult, ok := core.ResultAs[core.ResultWithValue](ownerResult); ok {
		return func(result core.ResultWithValue) any {
			return addDocumentOwner(ownerResult, getDocument(result))
		}
	}
	return getDocument
}

func wrapInTerminatorExecutor(exec core.Executor, wtExt map[string]any) (core.TerminatorExecutor, error) {
	return wrapInTerminatorExecutorWithOwnerResult(exec, wtExt, nil)
}

func wrapInTerminatorExecutorWithOwnerResult(exec core.Executor, wtExt map[string]any, ownerResult core.Result) (core.TerminatorExecutor, error) {
	cfg := core.WaitTerminationConfig{}
	if err := utils.DecodeValue(wtExt, &cfg); err != nil {
		return nil, fmt.Errorf("invalid wait-termination: %w", err)
	}
	return cfg.Build(exec, createGetDocument(ownerResult))
}
