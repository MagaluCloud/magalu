/*
Executor: set-current

# Description

# Sets profile to be used

import "magalu.cloud/lib/products/profile"
*/
package profile

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SetCurrentParameters struct {
	Name string `json:"name"`
}

type SetCurrentResult struct {
	Name string `json:"name"`
}

func SetCurrent(
	client *mgcClient.Client,
	ctx context.Context,
	parameters SetCurrentParameters,
) (
	result SetCurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("SetCurrent", mgcCore.RefPath("/profile/set-current"), client, ctx)
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
