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
	"context"

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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) AttachContext(
	ctx context.Context,
	parameters AttachParameters,
	configs AttachConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Attach", mgcCore.RefPath("/block-storage/volumes/attach"), s.client, ctx)
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

	sdkConfig := s.client.Sdk().Config().TempConfig()
	if c["serverUrl"] == nil && sdkConfig["serverUrl"] != nil {
		c["serverUrl"] = sdkConfig["serverUrl"]
	}

	if c["env"] == nil && sdkConfig["env"] != nil {
		c["env"] = sdkConfig["env"]
	}

	if c["region"] == nil && sdkConfig["region"] != nil {
		c["region"] = sdkConfig["region"]
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
