package auth

import (
	"context"
	"fmt"

	mgcAuthPkg "magalu.cloud/core/auth"

	"magalu.cloud/core"
)

type authSetParams struct {
	AccessKeyId     string `jsonschema_description:"Access key id value"`
	SecretAccessKey string `jsonschema_description:"Secret access key value"`
}

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

func newSet() core.Executor {
	executor := core.NewStaticExecute(
		"set",
		"",
		"Sets auth values",
		set,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Keys saved successfully\nAccessKeyId={{.AccessKeyId}}\nSecretAccessKey={{.SecretAccessKey}}\n"
	})
}
