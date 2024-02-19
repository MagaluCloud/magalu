package objects

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

var getDelete = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "delete",
			Description: "Delete an object from a bucket",
		},
		deleteObject,
	)
	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Deleted object %q", result.Source().Parameters["dst"])
	})
	exec = core.NewPromptInputExecutor(
		exec,
		core.NewPromptInput(
			"This command will delete the object at {{.confirmationValue}}, and its result is NOT reversible. Please confirm by retyping: {{.confirmationValue}}",
			"{{.parameters.dst}}",
		),
	)

	return exec
})

func deleteObject(ctx context.Context, params common.DeleteObjectParams, cfg common.Config) (core.Value, error) {
	err := common.Delete(ctx, params, cfg)
	if err != nil {
		return nil, err
	}
	return nil, err
}
