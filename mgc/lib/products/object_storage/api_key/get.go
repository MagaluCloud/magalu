/*
Executor: get

# Description

# Get details about a specific key

import "magalu.cloud/lib/products/object_storage/api_key"
*/
package apiKey

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	Uuid string `json:"uuid"`
}

type GetResult struct {
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

func Get(
	client *mgcClient.Client,
	ctx context.Context,
	parameters GetParameters,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/object-storage/api-key/get"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
