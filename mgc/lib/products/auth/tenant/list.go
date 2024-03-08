/*
Executor: list

# Description

# List all available tenants for current login

import "magalu.cloud/lib/products/auth/tenant"
*/
package tenant

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListResultItem struct {
	Email       string `json:"email"`
	IsDelegated bool   `json:"is_delegated"`
	IsManaged   bool   `json:"is_managed"`
	LegalName   string `json:"legal_name"`
	Uuid        string `json:"uuid"`
}

type ListResult []ListResultItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/auth/tenant/list"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
