/*
Executor: revoke

# Description

# Revoke credentials used in Object Storage requests

import "magalu.cloud/lib/products/object_storage/api_key"
*/
package apiKey

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RevokeParameters struct {
	Uuid string `json:"uuid"`
}

type RevokeResult struct {
	Uuid string `json:"uuid"`
}

func (s *service) Revoke(
	parameters RevokeParameters,
) (
	result RevokeResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Revoke", mgcCore.RefPath("/object-storage/api-key/revoke"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RevokeParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[RevokeResult](r)
}

// TODO: links
// TODO: related