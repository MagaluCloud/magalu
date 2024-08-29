/*
Executor: rename

# Summary

Rename a snapshot.

# Description

Rename a Snapshot for the currently authenticated tenant.

#### Rules
- The Snapshot name must be unique; otherwise, renaming will not be allowed.

#### Notes
  - Utilize the **block-storage snapshots list** command to retrieve a list of
    all Snapshots and obtain the ID of the Snapshot you wish to rename.

Version: v1

import "magalu.cloud/lib/products/block_storage/snapshots"
*/
package snapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RenameParameters struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Rename", mgcCore.RefPath("/block-storage/snapshots/rename"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RenameParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[map[string]interface{}](s.client.Sdk().Config().TempConfig()); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
