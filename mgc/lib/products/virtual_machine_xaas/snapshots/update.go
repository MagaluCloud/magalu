/*
Executor: update

# Summary

# Update snapshot

# Description

# Internal route for update a snapshot using its internal id

Version: 1.249.0

import "magalu.cloud/lib/products/virtual_machine_xaas/snapshots"
*/
package snapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UpdateParameters struct {
	Error         *string `json:"error,omitempty"`
	Id            string  `json:"id"`
	Size          *int    `json:"size,omitempty"`
	State         *string `json:"state,omitempty"`
	Status        *string `json:"status,omitempty"`
	UrpSnapshotId *string `json:"urp_snapshot_id"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/virtual-machine-xaas/snapshots/update"), s.client, s.ctx)
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
