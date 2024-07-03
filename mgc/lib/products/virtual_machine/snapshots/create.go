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

Version: v1

import "magalu.cloud/lib/products/virtual_machine/snapshots"
*/
package snapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Name           string                         `json:"name"`
	VirtualMachine CreateParametersVirtualMachine `json:"virtual_machine"`
}

// any of: CreateParametersVirtualMachine
type CreateParametersVirtualMachine struct {
	Id             string                                        `json:"id"`
	Name           *string                                       `json:"name,omitempty"`
	SecurityGroups *CreateParametersVirtualMachineSecurityGroups `json:"security_groups,omitempty"`
}

type CreateParametersVirtualMachineSecurityGroupsItem struct {
	Id string `json:"id"`
}

type CreateParametersVirtualMachineSecurityGroups []CreateParametersVirtualMachineSecurityGroupsItem

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
	Id string `json:"id"`
}

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine/snapshots/create"), s.client, s.ctx)
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
