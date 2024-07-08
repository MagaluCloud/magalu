package api_key

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"

	"magalu.cloud/core"
)

var getGetCurrent = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "current",
			Description: "Get the current Object Storage credentials",
		},
		getCurrent,
	)
})

var currentLogger = utils.NewLazyLoader(func() *zap.SugaredLogger {
	return logger().Named("current")
})

func getCurrent(ctx context.Context) (*apiKeysResult, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("unable to retrieve authentication configuration")
	}

	id, secretKey := auth.AccessKeyPair()
	if id == "" && secretKey == "" {
		fmt.Print("🔓 No current key pair set! \n\n")
		return &apiKeysResult{}, nil
	}

	keys, err := list(ctx)
	if err != nil {
		currentLogger().Warnw("Failed to get detailed info about current key, returning only KeyPairID and SecretKey", "err", err)
		return &apiKeysResult{KeyPairID: id, KeyPairSecret: secretKey}, nil
	}

	for _, key := range keys {
		if key.KeyPairID == id && key.KeyPairSecret == secretKey {
			return key, nil
		}
	}

	currentLogger().Warnw("unable to find a key in key pair list that matches the current KeyPairID", "keyPairId", id)
	return &apiKeysResult{KeyPairID: id, KeyPairSecret: secretKey}, nil
}
