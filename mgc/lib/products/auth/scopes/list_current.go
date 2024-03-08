/*
Executor: list-current

# Description

# List scopes present in the current access token

import "magalu.cloud/lib/products/auth/scopes"
*/
package scopes

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListCurrentResult []string

func ListCurrent(
	client *mgcClient.Client,
	ctx context.Context,
) (
	result ListCurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("ListCurrent", mgcCore.RefPath("/auth/scopes/list-current"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListCurrentResult](r)
}

// TODO: links
// TODO: related
