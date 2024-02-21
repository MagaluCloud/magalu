package api_key

import (
	"context"
	"fmt"

	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"

	"magalu.cloud/core"
)

var getGet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get the current Object Storage credentials",
		},
		get,
	)
})

func get(ctx context.Context) (*authSetParams, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("unable to retrieve authentication configuration")
	}

	id, secretKey := auth.AccessKeyPair()
	return &authSetParams{AccessKeyId: id, SecretAccessKey: secretKey}, nil
}
