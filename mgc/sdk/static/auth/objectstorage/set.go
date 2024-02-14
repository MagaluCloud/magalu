package objectstorage

import (
	"context"
	"fmt"

	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"

	"magalu.cloud/core"
)

type authSetParams struct {
	AccessKeyId     string `json:"access_key_id" jsonschema_description:"Access key id value" mgc:"positional"`
	SecretAccessKey string `json:"secret_access_key" jsonschema_description:"Secret access key value" mgc:"positional"`
}

var getSet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "Set the credentials values used for Object Storage requests",
		},
		set,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Keys saved successfully\nAccessKeyId={{.access_key_id}}\nSecretAccessKey={{.secret_access_key}}\n"
	})
})

func set(ctx context.Context, parameter authSetParams, _ struct{}) (*authSetParams, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("unable to retrieve authentication configuration")
	}

	if err := auth.SetAccessKey(parameter.AccessKeyId, parameter.SecretAccessKey); err != nil {
		return nil, err
	}

	return &parameter, nil
}
