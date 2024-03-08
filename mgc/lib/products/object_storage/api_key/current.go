/*
Executor: current

# Description

# Get the current Object Storage credentials

import "magalu.cloud/lib/products/object_storage/api_key"
*/
package apiKey

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CurrentResult struct {
	Description   string `json:"description"`
	EndValidity   string `json:"end_validity,omitempty"`
	KeyPairId     string `json:"key_pair_id"`
	KeyPairSecret string `json:"key_pair_secret"`
	Name          string `json:"name"`
	RevokedAt     string `json:"revoked_at,omitempty"`
	StartValidity string `json:"start_validity"`
	TenantName    string `json:"tenant_name,omitempty"`
	Uuid          string `json:"uuid"`
}

func Current(
	client *mgcClient.Client,
	ctx context.Context,
) (
	result CurrentResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Current", mgcCore.RefPath("/object-storage/api-key/current"), client, ctx)
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
