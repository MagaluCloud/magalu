/*
Executor: stop

# Summary

Stops a running instance.

# Description

Stops a Virtual Machine instance with the id provided in the current tenant which is logged in.
#### Notes
- You can use the virtual-machine list command to retrieve all instances, so you can get the id of
the instance that you want to stop.

#### Rules
- The instance must be in the running state.

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/instances/latest_instances"
*/
package latestInstances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type StopParameters struct {
	Id string `json:"id"`
}

type StopConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

func Stop(
	client *mgcClient.Client,
	ctx context.Context,
	parameters StopParameters,
	configs StopConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Stop", mgcCore.RefPath("/virtual-machine/instances/latest-instances/stop"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[StopParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[StopConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
