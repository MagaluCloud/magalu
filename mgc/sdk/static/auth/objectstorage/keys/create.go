package keys

import (
	"context"
	"fmt"

	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"

	"magalu.cloud/core"
)

type createParams struct {
	NameApiKey string `json:"name" jsonschema_description:"Name of new api key" mgc:"positional"`
}

var getCreate = utils.NewLazyLoader[core.Executor](newCreate)

func newCreate() core.Executor {
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
}

func create(ctx context.Context, parameter createParams, _ struct{}) (*mgcAuthPkg.ApiKeyResult, error) {
	auth := mgcAuthPkg.FromContext(ctx)

	if auth == nil {
		return nil, fmt.Errorf("unable to retrieve authentication configuration")
	}

	result, err := auth.CreateApiKey(ctx, parameter.NameApiKey)
	if err != nil {
		return nil, err
	}

	return result, nil
}
