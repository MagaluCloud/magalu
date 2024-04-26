/*
Executor: current

# Description

Shows current selected profile. Any changes to auth or config values will only affect this profile

import "magalu.cloud/lib/products/profile"
*/
package profile

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CurrentResult struct {
	Name string `json:"name"`
}

func (s *service) Current() (
	result CurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Current", mgcCore.RefPath("/profile/current"), s.client, s.ctx)
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
