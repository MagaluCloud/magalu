/*
Executor: start

# Summary

Starts a running or suspended instance.

# Description

Starts a Virtual Machine instance with the id provided in the current tenant which is logged in.
#### Notes
- You can use the virtual-machine list command to retrieve all instances,
so you can get the id of the instance that you want to start.

#### Rules
- The instance must be in the stopped or suspended states.

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

type StartParameters struct {
	Id string `json:"id"`
}

type StartConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

func Start(
	client *mgcClient.Client,
	ctx context.Context,
	parameters StartParameters,
	configs StartConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Start", mgcCore.RefPath("/virtual-machine/instances/v1-instances/start"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[StartParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[StartConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
