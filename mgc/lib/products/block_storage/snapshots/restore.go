/*
Executor: restore

# Summary

# Restore a snapshot to a new volume

# Description

Restore a Snapshot on a new Volume to the currently

	authenticated tenant.

The restored Volume can be used when it reaches the "available" state and the

	"completed" status.

#### Notes
  - To obtain the ID of the Snapshot you wish to restore, you can use the
    **block-storage snapshots list** command to list all Snapshots.
  - Check the state and status of your Volume using the
    **block-storage volume get --id [uuid]** command.

Version: v1

import "magalu.cloud/lib/products/block_storage/snapshots"
*/
package snapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RestoreParameters struct {
	Name       string `json:"name"`
	SnapshotId string `json:"snapshot_id"`
}

type RestoreConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type RestoreResult struct {
	Id string `json:"id"`
}

func (s *service) Restore(
	parameters RestoreParameters,
	configs RestoreConfigs,
) (
	result RestoreResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Restore", mgcCore.RefPath("/block-storage/snapshots/restore"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RestoreParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[RestoreConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[RestoreResult](r)
}

// TODO: links
// TODO: related
