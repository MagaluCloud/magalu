/*
Executor: set

# Description

# Change current Object Storage credential to selected

import "magalu.cloud/lib/products/object_storage/api_key"
*/
package apiKey

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SetParameters struct {
	Uuid string `json:"uuid"`
}

type SetResult struct {
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

func Set(
	client *mgcClient.Client,
	ctx context.Context,
	parameters SetParameters,
) (
	result SetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Set", mgcCore.RefPath("/object-storage/api-key/set"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[SetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[SetResult](r)
}

// TODO: links
// TODO: related
