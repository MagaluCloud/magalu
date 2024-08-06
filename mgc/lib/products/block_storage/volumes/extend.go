/*
Executor: extend

# Summary

# Extend a Volume

# Description

Extend the size of an existing Volume for the currently

	authenticated tenant.

The Volume extension will be completed when the Volume status returns to

	"completed".

#### Rules
- The Volume state must be "available".
- The Volume status must be "completed" or "extend_error".
- The new Volume size must be larger than the current size.

#### Notes
  - Utilize the block-storage volume list command to retrieve a list of all
    Volumes and obtain the ID of the Volume you want to extend.
  - Storage quotas are managed internally. If the operation fails due to quota
    restrictions, please contact our support team for assistance.

Version: v1

import "magalu.cloud/lib/products/block_storage/volumes"
*/
package volumes

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ExtendParameters struct {
	Id   string `json:"id"`
	Size int    `json:"size"`
}

type ExtendConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Extend(
	parameters ExtendParameters,
	configs ExtendConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Extend", mgcCore.RefPath("/block-storage/volumes/extend"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ExtendParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ExtendConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
