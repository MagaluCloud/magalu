package api_key

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type getKeyParams struct {
	UUID string `json:"uuid" mgc:"positional"`
}

var getGet = utils.NewLazyLoader(func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get details about a specific key",
		},
		get,
	)
})

func get(ctx context.Context, params getKeyParams, _ struct{}) (result *apiKeysResult, err error) {
	apiList, err := list(ctx)
	if err != nil {
		return
	}

	for _, v := range apiList {
		if v.UUID == params.UUID {
			result = v
			return
		}
	}

	err = fmt.Errorf("unable to find key with UUID %q", params.UUID)
	return
}
