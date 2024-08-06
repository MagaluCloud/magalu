/*
Executor: create

# Summary

# Create a backup of a copy

# Description

Create a backup of a Virtual Machine with the current tenant which is logged in.

A Backup is ready for restore when it's in completed status.

#### Rules
- It's possible to create a maximum of 100 backups per virtual machine.
- In case quota reached, choose a backup to remove.
- You can inform ID or Name from a Virtual Machine if both informed the priority will be ID.
- It's only possible to create a backup of a valid virtual machine.
- Each backup must have a unique name. It's not possible to create backups with the same name.

Version: 1.258.0

import "magalu.cloud/lib/products/virtual_machine_xaas/backups"
*/
package backups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	BackupUrl      string `json:"backup_url"`
	MinDisk        int    `json:"min_disk"`
	Name           string `json:"name"`
	OriginBackupId string `json:"origin_backup_id"`
	OriginRegion   string `json:"origin_region"`
	ProjectType    string `json:"project_type"`
}

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
	Id string `json:"id"`
}

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine-xaas/backups/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
