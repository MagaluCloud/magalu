/*
Executor: delete

# Summary

Delete/unset a Config value that had been previously set

# Description

Delete/unset a Config value that had been previously set. This does not
affect the environment variables

import "magalu.cloud/lib/products/config"
*/
package config

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteParameters struct {
	Key string `json:"key"`
}

type DeleteResult any

func Delete(
	client *mgcClient.Client,
	ctx context.Context,
	parameters DeleteParameters,
) (
	result DeleteResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Delete", mgcCore.RefPath("/config/delete"), client, ctx)
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
