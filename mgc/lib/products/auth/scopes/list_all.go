/*
Executor: list-all

# Description

# List all available scopes for all commands

import "magalu.cloud/lib/products/auth/scopes"
*/
package scopes

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListAllParameters struct {
	Target ListAllParametersTarget `json:"target,omitempty"`
}

type ListAllParametersTarget []string

type ListAllResult []string

func ListAll(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListAllParameters,
) (
	result ListAllResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("ListAll", mgcCore.RefPath("/auth/scopes/list-all"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListAllParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListAllResult](r)
}

// TODO: links
// TODO: related
