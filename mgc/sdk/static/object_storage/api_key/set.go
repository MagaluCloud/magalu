package api_key

import (
	"context"
	"fmt"

	mgcAuthPkg "magalu.cloud/core/auth"
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
		return "template=Keys changed successfully\nTenant=\"{{.tenant_name}}\"\nApiKey=\"{{.name}}\"\n{{if .description}}Description=\"{{.description}}\"{{- else}}{{end}}\n"
	})
})

func selectKey(ctx context.Context, parameter selectParams, _ struct{}) (*apiKeysResult, error) {

	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("could not get Auth from context")
	}

	apiList, err := list(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range apiList {
		if v.UUID == parameter.UUID {
			if err = mgcAuthPkg.FromContext(ctx).SetAccessKey(v.KeyPairID, v.KeyPairSecret); err != nil {
				return nil, err
			}
			return v, nil
		}
	}

	return nil, fmt.Errorf("the  API key (%s) is no longer valid", parameter.UUID)

}
