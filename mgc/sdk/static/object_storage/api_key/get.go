package api_key

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"
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
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		err = fmt.Errorf("could not get Auth from context")
		return
	}

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
