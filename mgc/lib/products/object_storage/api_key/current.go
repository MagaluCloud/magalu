/*
Executor: current

# Description

# Get the current Object Storage credentials

import "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/api_key"
*/
package apiKey

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type CurrentResult struct {
	Description   string  `json:"description"`
	EndValidity   *string `json:"end_validity,omitempty"`
	KeyPairId     string  `json:"key_pair_id"`
	KeyPairSecret string  `json:"key_pair_secret"`
	Name          string  `json:"name"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	StartValidity string  `json:"start_validity"`
	TenantName    *string `json:"tenant_name,omitempty"`
	Uuid          string  `json:"uuid"`
}

func (s *service) Current() (
	result CurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Current", mgcCore.RefPath("/object-storage/api-key/current"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) CurrentContext(
	ctx context.Context,
) (
	result CurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Current", mgcCore.RefPath("/object-storage/api-key/current"), s.client, ctx)
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
