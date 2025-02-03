package api_key

import (
	"context"

	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/utils"

	"github.com/MagaluCloud/magalu/mgc/core"
)

type addParams struct {
	KeyPairID     string `json:"keyId" jsonschema_description:"ID of api key to use" mgc:"positional"`
	KeyPairSecret string `json:"keySecret" jsonschema_description:"Secret of api key to use" mgc:"positional"`
}

var getAdd = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "add",
			Description: "Change current Object Storage credential",
			IsInternal:  utils.BoolPtr(true),
		},
		addKey,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Keys changed successfully\n"
	})
})

func addKey(ctx context.Context, parameter addParams, _ struct{}) (*apiKeysResult, error) {
	if err := mgcAuthPkg.FromContext(ctx).SetAccessKey(parameter.KeyPairID, parameter.KeyPairSecret); err != nil {
		return nil, err
	}

	return &apiKeysResult{KeyPairID: parameter.KeyPairID, KeyPairSecret: parameter.KeyPairSecret}, nil
}
