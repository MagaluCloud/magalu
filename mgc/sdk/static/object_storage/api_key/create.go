package api_key

import (
	"context"

	"magalu.cloud/core/utils"

	"magalu.cloud/core"
)

type createParams struct {
	ApiKeyName        string  `json:"name" jsonschema:"description=Name of new api key" mgc:"positional"`
	ApiKeyDescription *string `json:"description,omitempty" jsonschema:"description=Description of new api key" mgc:"positional"`
	ApiKeyExpiration  *string `json:"expiration,omitempty" jsonschema:"description=Date to expire new api,example=2024/11/07" mgc:"positional"`
}

var getCreate = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "create",
			Description: "Create new credentials used for Object Storage requests",
		},
		create,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Key created successfully\nUuid={{.uuid}}\n"
	})
})

func create(ctx context.Context, parameter createParams, _ struct{}) (*ApiKeyResult, error) {
	result, err := NewApiKey(ctx, parameter.ApiKeyName, parameter.ApiKeyDescription, parameter.ApiKeyExpiration)
	if err != nil {
		return nil, err
	}

	return result, nil
}
