/*
Executor: create

# Summary

Create a snapshot of a virtual machine asynchronously.

# Description

Create a snapshot of a Virtual Machine in the current tenant which is logged in. </br>
A Snapshot is ready for restore when it's in available state.

#### Notes
- You can verify the state of snapshot using the snapshot get command,
- To create a snapshot it's mandatory inform a valid and unique name.

#### Rules
- It's only possible to create a snapshot of a valid virtual machine.
- It's not possible to create 2 snapshots with the same name.
- You can inform ID or Name from a Virtual Machine if both informed the priority will be ID.

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/snapshots/v1_snapshots"
*/
package v1Snapshots

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Name           string                         `json:"name"`
	VirtualMachine CreateParametersVirtualMachine `json:"virtual_machine"`
}

// any of: CreateParametersVirtualMachine0, CreateParametersVirtualMachine1
type CreateParametersVirtualMachine struct {
	CreateParametersVirtualMachine0 `json:",squash"` // nolint
	CreateParametersVirtualMachine1 `json:",squash"` // nolint
}

type CreateParametersVirtualMachine0 struct {
	Id string `json:"id"`
}

type CreateParametersVirtualMachine1 struct {
	Name string `json:"name"`
}

type CreateConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
	Id string `json:"id"`
}

func Create(
	client *mgcClient.Client,
	ctx context.Context,
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine/snapshots/v1-snapshots/create"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
