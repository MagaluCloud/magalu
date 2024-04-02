/*
Executor: retype

# Summary

Changes a running or stopped instance machine type for another one.

# Description

Changes a Virtual Machine instance machine type with the id provided in the current tenant
which is logged in.

#### Notes
- You can use the virtual-machine list command to retrieve all instances, so you can get
the id of the instance that you want to change the machine type.

#### Rules
- The instance must be in the running or stopped state.
- The new machine type must be compatible with the current machine type.
- The new machine type must be available in the same region as the current machine type.
- You must provide either the machine type id or the machine type name, if you provide both,
the machine type id will be used.

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

type RetypeParameters struct {
	Id          string                      `json:"id"`
	MachineType RetypeParametersMachineType `json:"machine_type"`
}

// any of: RetypeParametersMachineType0, RetypeParametersMachineType1
type RetypeParametersMachineType struct {
	RetypeParametersMachineType0 `json:",squash"` // nolint
	RetypeParametersMachineType1 `json:",squash"` // nolint
}

type RetypeParametersMachineType0 struct {
	Id string `json:"id"`
}

type RetypeParametersMachineType1 struct {
	Name string `json:"name"`
}

type RetypeConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

func Retype(
	client *mgcClient.Client,
	ctx context.Context,
	parameters RetypeParameters,
	configs RetypeConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Retype", mgcCore.RefPath("/virtual-machine/instances/latest-instances/retype"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RetypeParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[RetypeConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
