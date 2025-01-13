package objects

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type publicUrlParams struct {
	Destination mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Path of the object to generate the public url,example=bucket1/file.txt" mgc:"positional"`
}

var getPublicUrl = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "public-url",
			Description: "Get object public url",
		},
		objectPublicUrl,
	)
	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template={{.url}}\n"
	})
})

func objectPublicUrl(ctx context.Context, p publicUrlParams, cfg common.Config) (*common.PublicUrlResult, error) {
	return common.PublicUrl(ctx, cfg, p.Destination)
}
