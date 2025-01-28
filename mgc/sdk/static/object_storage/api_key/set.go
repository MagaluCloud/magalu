package api_key

import (
	"context"

	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/utils"

	"github.com/MagaluCloud/magalu/mgc/core"
)

type selectParams struct {
	UUID string `json:"uuid" jsonschema_description:"UUID of api key to select" mgc:"positional"`
}

var getSetCurrent = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "Change current Object Storage credential to selected",
		},
		setCurrent,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Keys changed successfully\nTenant=\"{{.tenant_name}}\"\nApiKey=\"{{.name}}\"\n{{if .description}}Description=\"{{.description}}\"{{- else}}{{end}}\n"
	})
})

func setCurrent(ctx context.Context, parameter selectParams, _ struct{}) (*apiKeysResult, error) {
	key, err := get(ctx, getKeyParams(parameter), struct{}{})
	if err != nil {
		return nil, err
	}

	err = mgcAuthPkg.FromContext(ctx).SetAccessKey(key.KeyPairID, key.KeyPairSecret)
	if err != nil {
		return nil, err
	}

	return key, nil

}
