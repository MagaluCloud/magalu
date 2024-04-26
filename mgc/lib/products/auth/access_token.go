/*
Executor: access_token

# Description

# Retrieve the access token used in the APIs

import "magalu.cloud/lib/products/auth"
*/
package auth

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type AccessTokenParameters struct {
	Validate bool `json:"Validate,omitempty"`
}

type AccessTokenResult struct {
	AccessToken string `json:"access_token,omitempty"`
}

func (s *service) AccessToken(
	parameters AccessTokenParameters,
) (
	result AccessTokenResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("AccessToken", mgcCore.RefPath("/auth/access_token"), s.client, s.ctx)
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
