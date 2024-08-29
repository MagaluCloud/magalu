/*
Executor: rename

# Summary

# Renames a backup

# Description

Renames a backup with the id provided in the current tenant which is logged in.

#### Rules
- You can use the backup list command to retrieve all backups, so you can get the id of
the backup that you want to rename.
- A Backup can only renamed when it's in completed status.
- Each backup must have a unique name.

Version: v1

import "magalu.cloud/lib/products/virtual_machine/backups"
*/
package backups

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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Rename", mgcCore.RefPath("/virtual-machine/backups/rename"), s.client, s.ctx)
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
