/*
Executor: set

# Description

# Sets profile to be used

import "magalu.cloud/lib/products/profile"
*/
package profile

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SetParameters struct {
	Name string `json:"name"`
}

type SetResult struct {
	Name string `json:"name"`
}

func (s *service) Set(
	parameters SetParameters,
) (
	result SetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Set", mgcCore.RefPath("/profile/set"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[SetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[SetResult](r)
}

// TODO: links
// TODO: related
