/*
Executor: revoke

# Description

# Revoke an API key by its ID

import "magalu.cloud/lib/products/auth/api_key"
*/
package apiKey

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RevokeParameters struct {
	Id string `json:"id"`
}

type RevokeResult struct {
	Id string `json:"id"`
}

func (s *service) Revoke(
	parameters RevokeParameters,
) (
	result RevokeResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Revoke", mgcCore.RefPath("/auth/api-key/revoke"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RevokeParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[RevokeResult](r)
}

// TODO: links
// TODO: related
