/*
Executor: current

# Description

Shows current selected profile. Any changes to auth or config values will only affect this profile

import "magalu.cloud/lib/products/profile"
*/
package profile

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CurrentResult struct {
	Name string `json:"name"`
}

func Current(
	client *mgcClient.Client,
	ctx context.Context,
) (
	result CurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Current", mgcCore.RefPath("/profile/current"), client, ctx)
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
