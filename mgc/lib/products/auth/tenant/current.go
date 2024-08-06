/*
Executor: current

# Summary

# Get the currently active Tenant

# Description

# The current Tenant is used for all Magalu HTTP requests

import "magalu.cloud/lib/products/auth/tenant"
*/
package tenant

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CurrentResult struct {
	Email       string `json:"email"`
	IsDelegated bool   `json:"is_delegated"`
	IsManaged   bool   `json:"is_managed"`
	LegalName   string `json:"legal_name"`
	Uuid        string `json:"uuid"`
}

func (s *service) Current() (
	result CurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Current", mgcCore.RefPath("/auth/tenant/current"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CurrentResult](r)
}

// TODO: links
// TODO: related
