/*
Executor: current

# Description

# Retrieve the api-key used in the APIs

import "magalu.cloud/lib/products/auth/api_key"
*/
package apiKey

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CurrentResult struct {
	ApiKey *string `json:"api_key,omitempty"`
}

func (s *service) Current() (
	result CurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Current", mgcCore.RefPath("/auth/api-key/current"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CurrentResult](r)
}

// TODO: links
// TODO: related
