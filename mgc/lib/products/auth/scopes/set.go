/*
Executor: set

# Summary

Set the scopes for the current scopes in the access token.

# Description

Set the scopes for the current scopes in the access token.
Run 'auth scopes list-all' to see all available scopes

import "magalu.cloud/lib/products/auth/scopes"
*/
package scopes

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SetParameters struct {
	Scopes SetParametersScopes `json:"scopes"`
}

type SetParametersScopes []string

type SetResult []string

func Set(
	client *mgcClient.Client,
	ctx context.Context,
	parameters SetParameters,
) (
	result SetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Set", mgcCore.RefPath("/auth/scopes/set"), client, ctx)
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
