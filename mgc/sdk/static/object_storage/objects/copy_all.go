package objects

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

var getCopyAll = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "copy-all",
			Description: "Copy all objects from a bucket to another bucket",
		},
		copyAll,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Copied from {{.src}} to {{.dst}}\n"
	})
})

func copyAll(ctx context.Context, params common.CopyAllObjectsParams, cfg common.Config) (common.CopyAllObjectsParams, error) {
	err := common.CopyMultipleFiles(ctx, cfg, params)
	return params, err
}
