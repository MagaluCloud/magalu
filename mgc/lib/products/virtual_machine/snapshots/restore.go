/*
Executor: restore

# Summary

# Restore a snapshot to a virtual machine

# Description

Restore a snapshot of a Virtual Machine with the current tenant which is logged in. </br>
A Snapshot is ready for restore when it's in available state.

#### Notes
- You can verify the state of snapshot using the snapshot list command.
- Use machine-types list to see all machine types available.

#### Rules
- To restore a snapshot  you have to inform the new virtual machine information.
- You can choose a machine-type that has a disk equal or larger
than the original virtual machine type.

Version: v1

import "magalu.cloud/lib/products/virtual_machine/snapshots"
*/
package snapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RestoreParameters struct {
	AvailabilityZone *string                      `json:"availability_zone,omitempty"`
	Id               string                       `json:"id"`
	MachineType      RestoreParametersMachineType `json:"machine_type"`
	Name             string                       `json:"name"`
	Network          RestoreParametersNetwork     `json:"network"`
	SshKeyName       string                       `json:"ssh_key_name"`
	UserData         *string                      `json:"user_data,omitempty"`
}

// any of: RestoreParametersMachineType0, RestoreParametersMachineType1
type RestoreParametersMachineType struct {
	RestoreParametersMachineType0 `json:",squash"` // nolint
	RestoreParametersMachineType1 `json:",squash"` // nolint
}

type RestoreParametersMachineType0 struct {
	Id string `json:"id"`
}

type RestoreParametersMachineType1 struct {
	Name string `json:"name"`
}

type RestoreParametersNetwork struct {
	AssociatePublicIp *bool                              `json:"associate_public_ip,omitempty"`
	Interface         *RestoreParametersNetworkInterface `json:"interface,omitempty"`
	Vpc               *RestoreParametersNetworkVpc       `json:"vpc,omitempty"`
}

// any of: RestoreParametersMachineType0, RestoreParametersNetworkInterface1
type RestoreParametersNetworkInterface struct {
	RestoreParametersMachineType0      `json:",squash"` // nolint
	RestoreParametersNetworkInterface1 `json:",squash"` // nolint
}

type RestoreParametersNetworkInterface1 struct {
	SecurityGroups *RestoreParametersNetworkInterface1SecurityGroups `json:"security_groups,omitempty"`
}

type RestoreParametersNetworkInterface1SecurityGroupsItem struct {
	Id string `json:"id"`
}

type RestoreParametersNetworkInterface1SecurityGroups []RestoreParametersNetworkInterface1SecurityGroupsItem

// any of: RestoreParametersMachineType0, RestoreParametersMachineType1
type RestoreParametersNetworkVpc struct {
	RestoreParametersMachineType0 `json:",squash"` // nolint
	RestoreParametersMachineType1 `json:",squash"` // nolint
}

type RestoreConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type RestoreResult struct {
	Id string `json:"id"`
}

func (s *service) Restore(
	parameters RestoreParameters,
	configs RestoreConfigs,
) (
	result RestoreResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Restore", mgcCore.RefPath("/virtual-machine/snapshots/restore"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RestoreParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[RestoreConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[RestoreResult](r)
}

// TODO: links
// TODO: related
