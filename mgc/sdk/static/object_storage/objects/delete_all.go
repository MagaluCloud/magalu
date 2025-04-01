package objects

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

var getDeleteAll = utils.NewLazyLoader[core.Executor](func() core.Executor {
	exec := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "delete-all",
			Description: "Delete all objects from a bucket",
		},
		deleteAll,
	)

	return core.NewPromptInputExecutor(
		exec,
		core.NewPromptInput(
			`This command will delete all objects at {{.confirmationValue}}, and its result is NOT reversible.
Please confirm by retyping: {{.confirmationValue}}`,
			"{{.parameters.bucket}}",
		),
	)
})

func deleteAll(ctx context.Context, params common.DeleteAllObjectsInBucketParams, cfg common.Config) (result core.Value, err error) {
	err = common.DeleteAllObjectsInBucket(ctx, params, cfg)
	return
}
