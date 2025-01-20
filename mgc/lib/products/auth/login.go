/*
Executor: login

# Summary

# Authenticate with Magalu Cloud

# Description

Log in to your Magalu Cloud account. When you login with this command,
the current Tenant will always be set to the default one. To see more details
about a successful login, use the '--show' flag when logging in

import "github.com/MagaluCloud/magalu/mgc/lib/products/auth"
*/
package auth

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type LoginParameters struct {
	Headless *bool                  `json:"headless,omitempty"`
	Qrcode   *bool                  `json:"qrcode,omitempty"`
	Scopes   *LoginParametersScopes `json:"scopes,omitempty"`
	Show     *bool                  `json:"show,omitempty"`
}

type LoginParametersScopes []string

type LoginResult struct {
	AccessToken    *string                    `json:"access_token,omitempty"`
	SelectedTenant *LoginResultSelectedTenant `json:"selected_tenant,omitempty"`
}

type LoginResultSelectedTenant struct {
	Email       string `json:"email"`
	IsDelegated bool   `json:"is_delegated"`
	IsManaged   bool   `json:"is_managed"`
	LegalName   string `json:"legal_name"`
	Uuid        string `json:"uuid"`
}

func (s *service) Login(
	parameters LoginParameters,
) (
	result LoginResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Login", mgcCore.RefPath("/auth/login"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[LoginParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[LoginResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) LoginContext(
	ctx context.Context,
	parameters LoginParameters,
) (
	result LoginResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Login", mgcCore.RefPath("/auth/login"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[LoginParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[LoginResult](r)
}

// TODO: links
// TODO: related
