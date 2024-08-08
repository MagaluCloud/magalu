package api_key

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

type getKeyParams struct {
	UUID string `json:"id" mgc:"positional"`
}

var getGet = utils.NewLazyLoader(func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get a specific API key by its ID",
		},
		get,
	)
})

func get(ctx context.Context, params getKeyParams, _ struct{}) (result getApiKeyResult, err error) {
	apiList, err := listFull(ctx, true)
	if err != nil {
		return
	}

	for _, v := range apiList {
		if v.UUID == params.UUID {
			result = *v.ToResultGet()
			return
		}
	}

	err = fmt.Errorf("unable to find key with ID %q", params.UUID)
	return
}
