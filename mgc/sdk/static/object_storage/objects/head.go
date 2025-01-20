package objects

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type headObjectParams struct {
	Destination mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Path of the object to be get metadata from,example=bucket1/file.txt" mgc:"positional"`
	Version     string           `json:"objVersion,omitempty" jsonschema:"description=Version of the object to be get metadata from"`
}

var getHead = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "head",
			Description: "Get object metadata",
		},
		headObject,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "table"
	})
	return exec
})

func headObject(ctx context.Context, p headObjectParams, cfg common.Config) (common.HeadObjectResponse, error) {
	return common.HeadFile(ctx, cfg, p.Destination, p.Version)
}
