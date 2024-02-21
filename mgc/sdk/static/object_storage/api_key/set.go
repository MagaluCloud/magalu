package api_key

import (
	"context"

	"magalu.cloud/core/utils"

	"magalu.cloud/core"
)

type selectParams struct {
	UUID string `json:"uuid" jsonschema_description:"UUID of api key to select" mgc:"positional"`
}

var getSet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "Change current Object Storage credential to selected",
		},
		selectKey,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Keys changed successfully\nTenant={{.tenant_name}}\nApiKey Name={{.name}}\nDescription={{.description}}\n"
	})
})

func selectKey(ctx context.Context, parameter selectParams, _ struct{}) (*apiKeysResult, error) {
	result, err := SelectApiKey(ctx, parameter.UUID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
