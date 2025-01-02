/*
Executor: create

# Summary

# Create a new API Key

# Description

# Select the scopes that the new API Key will have access to and set an expiration date

import "github.com/MagaluCloud/magalu/mgc/lib/products/auth/api_key"
*/
package apiKey

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type CreateParameters struct {
	Description *string `json:"description,omitempty"`
	Expiration  *string `json:"expiration,omitempty"`
	Name        string  `json:"name"`
}

type CreateResult struct {
	Used *bool   `json:"used,omitempty"`
	Uuid *string `json:"uuid,omitempty"`
}

func (s *service) Create(
	parameters CreateParameters,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/auth/api-key/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) CreateContext(
	ctx context.Context,
	parameters CreateParameters,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/auth/api-key/create"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
