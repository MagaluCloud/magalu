/*
Executor: list

# Summary

# List your account API keys

# Description

# This APIs Keys are from your account and can be used to authenticate in the Magalu Cloud

import "github.com/MagaluCloud/magalu/mgc/lib/products/auth/api_key"
*/
package apiKey

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type ListParameters struct {
	InvalidKeys bool `json:"invalid-keys"`
}

type ListResultItem struct {
	Description   *string `json:"description,omitempty"`
	EndValidity   *string `json:"end_validity,omitempty"`
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	StartValidity string  `json:"start_validity"`
	TenantName    *string `json:"tenant_name,omitempty"`
}

type ListResult []ListResultItem

func (s *service) List(
	parameters ListParameters,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/auth/api-key/list"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ListContext(
	ctx context.Context,
	parameters ListParameters,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/auth/api-key/list"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
