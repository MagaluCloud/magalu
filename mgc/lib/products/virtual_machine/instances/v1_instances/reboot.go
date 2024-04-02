/*
Executor: reboot

# Summary

Starts a running or suspended instance.

# Description

Reboots a Virtual Machine instance with the id provided in the current tenant which is logged in.
#### Notes
- You can use the virtual-machine list command to retrieve all instances, so you can get the id
of the instance that you want to reboot.

#### Rules
- The instance must be in the running state.

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/instances/v1_instances"
*/
package v1Instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RebootParameters struct {
	Id string `json:"id"`
}

type RebootConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

func Reboot(
	client *mgcClient.Client,
	ctx context.Context,
	parameters RebootParameters,
	configs RebootConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Reboot", mgcCore.RefPath("/virtual-machine/instances/v1-instances/reboot"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RebootParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[RebootConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
