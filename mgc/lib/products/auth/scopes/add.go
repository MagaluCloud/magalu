/*
Executor: add

# Summary

# Add new scopes to the current access token

# Description

Add new scopes to the current access token. Run 'auth scopes list-all'
to see all available scopes to be added

import "magalu.cloud/lib/products/auth/scopes"
*/
package scopes

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type AddParameters struct {
	Scopes AddParametersScopes `json:"scopes"`
}

type AddParametersScopes []string

type AddResult []string

func Add(
	client *mgcClient.Client,
	ctx context.Context,
	parameters AddParameters,
) (
	result AddResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Add", mgcCore.RefPath("/auth/scopes/add"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[AddParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[AddResult](r)
}

// TODO: links
// TODO: related
