package api_key

import (
	"context"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

type revokeParams struct {
	UUID string `json:"uuid" jsonschema_description:"UUID of api key to revoke" mgc:"positional"`
}

var getRevoke = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Scopes:      core.Scopes{"pa:api-keys:revoke"},
			Name:        "revoke",
			Description: "Revoke credentials used in Object Storage requests",
		},
		revoke,
	)

	msg := "This operation will permanently revoke the api-key {{.parameters.uuid}}. Do you wish to continue?"

	cExecutor := core.NewConfirmableExecutor(
		exec,
		core.ConfirmPromptWithTemplate(msg),
	)

	return core.NewExecuteResultOutputOptions(cExecutor, func(exec core.Executor, result core.Result) string {
		return "template=Revoked!\n"
	})
})

func revoke(ctx context.Context, parameter revokeParams, _ struct{}) (bool, error) {

	if err := RevokeApiKey(ctx, parameter.UUID); err != nil {
		return false, err
	}

	return true, nil
}
