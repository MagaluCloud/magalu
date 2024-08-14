package api_key

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"
)

type currentApiKeyResult struct {
	ApiKey string `json:"api_key,omitempty"`
}

var getCurrentApiKey = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "current",
			Description: "Retrieve the api-key used in the APIs",
		},
		func(ctx context.Context, _ struct{}, _ struct{}) (output *currentApiKeyResult, err error) {
			auth := mgcAuthPkg.FromContext(ctx)
			if auth == nil {
				return nil, fmt.Errorf("unable to retrieve authentication configuration")
			}

			apiKey, err := auth.ApiKey(ctx)
			if err != nil {
				return nil, err
			}

			return &currentApiKeyResult{ApiKey: apiKey}, nil
		},
	)
})
