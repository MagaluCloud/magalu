package keys

import (
	"context"
	"fmt"

	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"

	"magalu.cloud/core"
)

type revokeParams struct {
	UUID string `json:"uuid" jsonschema_description:"UUID of api key to revoke" mgc:"positional"`
}

var getRevoke = utils.NewLazyLoader[core.Executor](newRevoke)

func newRevoke() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "revoke",
			Description: "Revoke credentials used in Object Storage requests",
		},
		revoke,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=API Key revoked successfully\n"
	})
}

func revoke(ctx context.Context, parameter revokeParams, _ struct{}) (bool, error) {
	auth := mgcAuthPkg.FromContext(ctx)

	if auth == nil {
		return false, fmt.Errorf("unable to retrieve authentication configuration")
	}

	if err := auth.RevokeApiKey(ctx, parameter.UUID); err != nil {
		return false, err
	}

	return true, nil
}
