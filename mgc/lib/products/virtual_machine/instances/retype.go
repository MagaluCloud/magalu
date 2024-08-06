/*
Executor: retype

# Summary

Changes an instance machine-type.

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

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances"
*/
package instances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RetypeParameters struct {
	Id          string                      `json:"id"`
	MachineType RetypeParametersMachineType `json:"machine_type"`
}

// any of: RetypeParametersMachineType
type RetypeParametersMachineType struct {
	Id             string                                     `json:"id"`
	Name           *string                                    `json:"name,omitempty"`
	SecurityGroups *RetypeParametersMachineTypeSecurityGroups `json:"security_groups,omitempty"`
}

type RetypeParametersMachineTypeSecurityGroupsItem struct {
	Id string `json:"id"`
}

type RetypeParametersMachineTypeSecurityGroups []RetypeParametersMachineTypeSecurityGroupsItem

type RetypeConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Retype(
	parameters RetypeParameters,
	configs RetypeConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Retype", mgcCore.RefPath("/virtual-machine/instances/retype"), s.client, s.ctx)
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
