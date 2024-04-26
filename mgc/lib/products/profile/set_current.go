/*
Executor: set-current

# Description

# Sets profile to be used

import "magalu.cloud/lib/products/profile"
*/
package profile

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SetCurrentParameters struct {
	Name string `json:"name"`
}

type SetCurrentResult struct {
	Name string `json:"name"`
}

func (s *service) SetCurrent(
	parameters SetCurrentParameters,
) (
	result SetCurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("SetCurrent", mgcCore.RefPath("/profile/set-current"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[SetCurrentParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[SetCurrentResult](r)
}

// TODO: links
// TODO: related
