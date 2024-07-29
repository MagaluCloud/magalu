/*
Executor: copy

# Summary

# Copy backup cross region

# Description

Copy a backup cross region for the currently authenticated tenant.

#### Rules
- The copy only be accepted when the destiny region is different from origin region.
- The copy only be accepted if the backup's name in destiny region is different from input name.
- The copy only be accepted if the user has access to destiny region.

#### Notes
  - Utilize the **block-storage backups list** command to retrieve a list of
    all Backups and obtain the ID of the Backup you wish to copy across different region.

Version: v1

import "magalu.cloud/lib/products/block_storage/backups"
*/
package backups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CopyParameters struct {
	Backup        CopyParametersBackup `json:"backup"`
	DestinyRegion string               `json:"destiny_region"`
}

type CopyParametersBackup struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type CopyConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Copy(
	parameters CopyParameters,
	configs CopyConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Copy", mgcCore.RefPath("/block-storage/backups/copy"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CopyParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CopyConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
