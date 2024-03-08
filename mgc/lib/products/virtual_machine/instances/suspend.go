/*
Executor: suspend

# Summary

Suspends a running instance.

# Description

Suspends a Virtual Machine instance with the id provided in the current tenant which is logged in.

#### Notes
- You can use the virtual-machine list command to retrieve all instances, so you can get the id of
the instance that you want to suspend.

#### Rules
- The instance must be in the running state.

Version: 0.1.0

import "magalu.cloud/lib/products/virtual_machine/instances"
*/
package instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SuspendParameters struct {
	Id string `json:"id"`
}

type SuspendConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

func Suspend(
	client *mgcClient.Client,
	ctx context.Context,
	parameters SuspendParameters,
	configs SuspendConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Suspend", mgcCore.RefPath("/virtual-machine/instances/suspend"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[SuspendParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[SuspendConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
