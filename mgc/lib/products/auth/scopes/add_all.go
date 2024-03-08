/*
Executor: add-all

# Summary

# Add all scopes to the current access token

# Description

Add all scopes from all operations to the current access token.

import "magalu.cloud/lib/products/auth/scopes"
*/
package scopes

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type AddAllResult []string

func AddAll(
	client *mgcClient.Client,
	ctx context.Context,
) (
	result AddAllResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("AddAll", mgcCore.RefPath("/auth/scopes/add-all"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[AddAllResult](r)
}

// TODO: links
// TODO: related
