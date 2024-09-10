/*
Executor: set

# Description

# Set the active Tenant to be used for all subsequent requests

import "magalu.cloud/lib/products/auth/tenant"
*/
package tenant

import (
	"context"

	mgcCore "magalu.cloud/core"
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

func (s *service) Set(
	parameters SetParameters,
) (
	result SetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Set", mgcCore.RefPath("/auth/tenant/set"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) SetContext(
	ctx context.Context,
	parameters SetParameters,
) (
	result SetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Set", mgcCore.RefPath("/auth/tenant/set"), s.client, ctx)
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
