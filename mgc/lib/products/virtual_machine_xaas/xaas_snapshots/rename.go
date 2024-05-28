/*
Executor: rename

# Summary

# Renames a snapshot

# Description

Renames a Snapshot with the id provided in the current tenant which is logged in.

#### Notes
- You can use the snapshots list command to retrieve all snapshots, so you can get the id of
the snapshot that you want to rename.

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/xaas_snapshots"
*/
package xaasSnapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RenameParameters struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ProjectType string `json:"project_type"`
}

type RenameConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Rename(
	parameters RenameParameters,
	configs RenameConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Rename", mgcCore.RefPath("/virtual-machine-xaas/xaas snapshots/rename"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RenameParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[RenameConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
