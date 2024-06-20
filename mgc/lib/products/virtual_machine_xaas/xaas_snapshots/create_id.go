/*
Executor: create-id

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

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/xaas_snapshots"
*/
package xaasSnapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateIdParameters struct {
	AvailabilityZone *string                       `json:"availability_zone,omitempty"`
	Id               string                        `json:"id"`
	MachineType      CreateIdParametersMachineType `json:"machine_type"`
	Name             string                        `json:"name"`
	Network          CreateIdParametersNetwork     `json:"network"`
	ProjectType      string                        `json:"project_type"`
	SshKeyName       string                        `json:"ssh_key_name"`
	UserData         *string                       `json:"user_data,omitempty"`
}

// any of: CreateIdParametersMachineType0, CreateIdParametersMachineType1
type CreateIdParametersMachineType struct {
	CreateIdParametersMachineType0 `json:",squash"` // nolint
	CreateIdParametersMachineType1 `json:",squash"` // nolint
}

type CreateIdParametersMachineType0 struct {
	Id string `json:"id"`
}

type CreateIdParametersMachineType1 struct {
	Name string `json:"name"`
}

type CreateIdParametersNetwork struct {
	AssociatePublicIp *bool                         `json:"associate_public_ip,omitempty"`
	Nic               *CreateIdParametersNetworkNic `json:"nic,omitempty"`
	Vpc               *CreateIdParametersNetworkVpc `json:"vpc,omitempty"`
}

// any of: CreateIdParametersMachineType0, CreateIdParametersNetworkNic1
type CreateIdParametersNetworkNic struct {
	CreateIdParametersMachineType0 `json:",squash"` // nolint
	CreateIdParametersNetworkNic1  `json:",squash"` // nolint
}

type CreateIdParametersNetworkNic1 struct {
	SecurityGroups *CreateIdParametersNetworkNic1SecurityGroups `json:"security_groups,omitempty"`
}

type CreateIdParametersNetworkNic1SecurityGroups []CreateIdParametersMachineType0

// any of: CreateIdParametersMachineType0, CreateIdParametersMachineType1
type CreateIdParametersNetworkVpc struct {
	CreateIdParametersMachineType0 `json:",squash"` // nolint
	CreateIdParametersMachineType1 `json:",squash"` // nolint
}

type CreateIdConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type CreateIdResult struct {
	Id string `json:"id"`
}

func (s *service) CreateId(
	parameters CreateIdParameters,
	configs CreateIdConfigs,
) (
	result CreateIdResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("CreateId", mgcCore.RefPath("/virtual-machine-xaas/xaas snapshots/create-id"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateIdParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateIdConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateIdResult](r)
}

// TODO: links
// TODO: related
