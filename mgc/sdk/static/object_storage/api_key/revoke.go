package api_key

import (
	"context"
	"fmt"

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

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Revoked api-key %q", result.Source().Parameters["uuid"])
	})

	exec = core.NewPromptInputExecutor(
		exec,
		core.NewPromptInput(
			"This command will revoke the api-key {{.parameters.uuid}}, and its result is NOT reversible.\nPlease confirm by retyping: {{.confirmationValue}}",
			"yes",
		),
	)

	return exec
})

func revoke(ctx context.Context, parameter revokeParams, _ struct{}) (bool, error) {

	if err := RevokeApiKey(ctx, parameter.UUID); err != nil {
		return false, err
	}

	return true, nil
}
