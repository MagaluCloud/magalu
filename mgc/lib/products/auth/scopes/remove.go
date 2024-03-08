/*
Executor: remove

# Summary

Remove scopes from the current scopes in the access token.

# Description

Remove scopes from the current scopes in the access token.
Run 'auth scopes list-current' to see current scopes

import "magalu.cloud/lib/products/auth/scopes"
*/
package scopes

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RemoveParameters struct {
	Scopes RemoveParametersScopes `json:"scopes"`
}

type RemoveParametersScopes []string

type RemoveResult []string

func Remove(
	client *mgcClient.Client,
	ctx context.Context,
	parameters RemoveParameters,
) (
	result RemoveResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Remove", mgcCore.RefPath("/auth/scopes/remove"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RemoveParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[RemoveResult](r)
}

// TODO: links
// TODO: related
