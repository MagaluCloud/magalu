/*
Executor: urp

# Summary

# Urp Update Backup Status

# Description

# Update backup status by id on urp

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/backups"
*/
package backups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UrpParameters struct {
	Error            *string `json:"error,omitempty"`
	ExternalBackupId string  `json:"external_backup_id"`
	MinDisk          int     `json:"min_disk"`
	Size             int     `json:"size"`
	State            string  `json:"state"`
	Status           string  `json:"status"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Urp", mgcCore.RefPath("/virtual-machine-xaas/backups/urp"), s.client, s.ctx)
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
