/*
Executor: set

# Description

# Set the active Tenant to be used for all subsequential requests

import "magalu.cloud/lib/products/auth/tenant"
*/
package tenant

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SetParameters struct {
	Uuid string `json:"uuid"`
}

type SetResult struct {
	AccessToken  string             `json:"access_token"`
	CreatedAt    SetResultCreatedAt `json:"created_at"`
	RefreshToken string             `json:"refresh_token"`
	Scope        SetResultScope     `json:"scope"`
	Uuid         string             `json:"uuid"`
}

type SetResultCreatedAt struct {
}

type SetResultScope []string

func Set(
	client *mgcClient.Client,
	ctx context.Context,
	parameters SetParameters,
) (
	result SetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Set", mgcCore.RefPath("/auth/tenant/set"), client, ctx)
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
