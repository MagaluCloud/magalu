/*
Executor: access_token

# Description

# Retrieve the access token used in the APIs

import "magalu.cloud/lib/products/auth"
*/
package auth

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type AccessTokenParameters struct {
	Validate bool `json:"Validate,omitempty"`
}

type AccessTokenResult struct {
	AccessToken string `json:"access_token,omitempty"`
}

func AccessToken(
	client *mgcClient.Client,
	ctx context.Context,
	parameters AccessTokenParameters,
) (
	result AccessTokenResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("AccessToken", mgcCore.RefPath("/auth/access_token"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[AccessTokenParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[AccessTokenResult](r)
}

// TODO: links
// TODO: related
