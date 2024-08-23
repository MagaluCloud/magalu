/*
Executor: attach

# Summary

Attach the volume to an instance.

# Description

Attach a Volume to a Virtual Machine instance for the currently

	authenticated tenant.

The Volume attachment will be completed when the Volume status returns to
"completed", and the state becomes "in-use".

#### Rules
- The Volume and the Virtual Machine must belong to the same tenant.
- Both the Volume and Virtual Machine must have the status "completed".
- The Volume's state must be "available".
- The Virtual Machine's state must be "stopped" or "running".

#### Notes
  - Verify the state and status of your Volume using the
    **block-storage volume get --id [uuid]** command.
  - Verify the state and status of your Virtual Machine using the

**virtual-machine instances get --id [uuid]** command".

Version: v1

import "magalu.cloud/lib/products/block_storage/volumes"
*/
package volumes

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type AttachParameters struct {
	Id               string `json:"id"`
	VirtualMachineId string `json:"virtual_machine_id"`
}

type AttachConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Attach(
	parameters AttachParameters,
	configs AttachConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Attach", mgcCore.RefPath("/block-storage/volumes/attach"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[AttachParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[AttachConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
