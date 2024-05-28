/*
Executor: update

# Summary

# Update Backup Status

# Description

# Update backup status

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/backups"
*/
package backups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UpdateParameters struct {
	Error            *string `json:"error,omitempty"`
	ExternalBackupId string  `json:"external_backup_id"`
	Id               string  `json:"id"`
	Status           string  `json:"status"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/virtual-machine-xaas/backups/update"), s.client, s.ctx)
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
