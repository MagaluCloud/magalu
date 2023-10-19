package auth

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
)

type accessTokenParameters struct {
	Validate bool `json:",omitempty" jsonschema_description:"Validate the token, refreshing if needed"`
}

type accessTokenResult struct {
	AccessToken string `json:"access_token,omitempty"`
}

func newAccessToken() *core.StaticExecute {
	return core.NewStaticExecute(
		"access_token",
		"",
		"Retrieve the access token to use the APIs",
		func(ctx context.Context, parameters accessTokenParameters, _ struct{}) (output *accessTokenResult, err error) {
			auth := mgcAuthPkg.FromContext(ctx)
			if auth == nil {
				return nil, fmt.Errorf("unable to retrieve authentication configuration")
			}

			if parameters.Validate {
				err := auth.ValidateAccessToken(ctx)
				if err != nil {
					return nil, fmt.Errorf("Could not validate the Access Token: %w", err)
				}
			}

			token, err := auth.AccessToken(ctx)
			if err != nil {
				return nil, err
			}

			return &accessTokenResult{AccessToken: token}, nil
		},
	)
}
