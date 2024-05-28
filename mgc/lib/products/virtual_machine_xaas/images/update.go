/*
Executor: update

# Summary

# Update Image Status

# Description

# Update image status

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/images"
*/
package images

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UpdateParameters struct {
	Id         string  `json:"id"`
	IsApproved *bool   `json:"is_approved,omitempty"`
	Sku        *string `json:"sku,omitempty"`
	Status     *string `json:"status,omitempty"`
}

type UpdateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Update(
	parameters UpdateParameters,
	configs UpdateConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/virtual-machine-xaas/images/update"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UpdateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UpdateConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
