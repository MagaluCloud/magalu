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
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CurrentResult struct {
	Email       string `json:"email"`
	IsDelegated bool   `json:"is_delegated"`
	IsManaged   bool   `json:"is_managed"`
	LegalName   string `json:"legal_name"`
	Uuid        string `json:"uuid"`
}

func Current(
	client *mgcClient.Client,
	ctx context.Context,
) (
	result CurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Current", mgcCore.RefPath("/auth/tenant/current"), client, ctx)
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
