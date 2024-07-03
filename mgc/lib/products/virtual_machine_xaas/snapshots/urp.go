/*
Executor: urp

# Summary

# Worker Update Snapshot Urp

# Description

Internal route for update status of a snapshot when receive a update from URP.

### Note
This route is used only for internal proposes.

Version: 1.249.0

import "magalu.cloud/lib/products/virtual_machine_xaas/snapshots"
*/
package snapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UrpParameters struct {
	Error  *string `json:"error,omitempty"`
	Size   *int    `json:"size,omitempty"`
	State  *string `json:"state,omitempty"`
	Status *string `json:"status,omitempty"`
	UrpId  string  `json:"urp_id"`
}

type UrpConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Urp(
	parameters UrpParameters,
	configs UrpConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Urp", mgcCore.RefPath("/virtual-machine-xaas/snapshots/urp"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UrpParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UrpConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
