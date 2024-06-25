/*
Executor: create

# Summary

Create an instance asynchronously.

# Description

# Create a Virtual Machine instance

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances"
*/
package instances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	AvailabilityZone *string                     `json:"availability_zone,omitempty"`
	Image            CreateParametersImage       `json:"image"`
	MachineType      CreateParametersMachineType `json:"machine_type"`
	Name             string                      `json:"name"`
	Network          *CreateParametersNetwork    `json:"network,omitempty"`
	SshKeyName       string                      `json:"ssh_key_name"`
	UserData         *string                     `json:"user_data,omitempty"`
}

// any of: , CreateParametersImage1
type CreateParametersImage struct {
	CreateParametersImage1 `json:",squash"` // nolint
}

type CreateParametersImage1 struct {
	Name string `json:"name"`
}

// any of: , CreateParametersImage1
type CreateParametersMachineType struct {
	CreateParametersImage1 `json:",squash"` // nolint
}

type CreateParametersNetwork struct {
	AssociatePublicIp *bool                             `json:"associate_public_ip,omitempty"`
	Interface         *CreateParametersNetworkInterface `json:"interface,omitempty"`
	Vpc               *CreateParametersNetworkVpc       `json:"vpc,omitempty"`
}

// any of: , CreateParametersNetworkInterface1
type CreateParametersNetworkInterface struct {
	CreateParametersNetworkInterface1 `json:",squash"` // nolint
}

type CreateParametersNetworkInterface1 struct {
	SecurityGroups *CreateParametersNetworkInterface1SecurityGroups `json:"security_groups,omitempty"`
}

type CreateParametersNetworkInterface1SecurityGroupsItem struct {
	Id string `json:"id"`
}

type CreateParametersNetworkInterface1SecurityGroups []CreateParametersNetworkInterface1SecurityGroupsItem

// any of: , CreateParametersImage1
type CreateParametersNetworkVpc struct {
	CreateParametersImage1 `json:",squash"` // nolint
}

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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine/instances/create"), s.client, s.ctx)
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
