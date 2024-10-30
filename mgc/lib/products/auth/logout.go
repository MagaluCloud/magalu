/*
Executor: logout

# Description

# Run logout

import "magalu.cloud/lib/products/auth"
*/
package auth

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type LogoutParameters struct {
	Validate *bool `json:"Validate,omitempty"`
}

type LogoutResult string

func (s *service) Logout(
	parameters LogoutParameters,
) (
	result LogoutResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Logout", mgcCore.RefPath("/auth/logout"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[LogoutParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[LogoutResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) LogoutContext(
	ctx context.Context,
	parameters LogoutParameters,
) (
	result LogoutResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Logout", mgcCore.RefPath("/auth/logout"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[LogoutParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[LogoutResult](r)
}

// TODO: links
// TODO: related
