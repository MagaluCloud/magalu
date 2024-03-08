/*
Executor: delete

# Description

# Deletes the profile with the specified name

import "magalu.cloud/lib/products/profile"
*/
package profile

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteParameters struct {
	Name string `json:"name"`
}

type DeleteResult struct {
	Name string `json:"name"`
}

func Delete(
	client *mgcClient.Client,
	ctx context.Context,
	parameters DeleteParameters,
) (
	result DeleteResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Delete", mgcCore.RefPath("/profile/delete"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DeleteResult](r)
}

// TODO: links
// TODO: related
