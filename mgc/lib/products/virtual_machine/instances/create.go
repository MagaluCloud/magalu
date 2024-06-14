/*
Executor: create

# Summary

Create an instance asynchronously.

# Description

Creates a Virtual Machine instance in the current tenant which is logged in.

An instance is ready for you to use when it's in the running state.

#### Notes
  - For the image data, you can use the virtual-machine images list command
    to list all available images.
  - For the machine type data, you can use the virtual-machine machine-types
    list command to list all available machine types.
  - You can verify the state of your instance using the virtual-machine get

command.

#### Rules

- If you don't specify a VPC, the default VPC will be used. When the
default VPC is not available, the command will fail.
- If you don't specify an network interface, an default network interface
will be created.
- You can either specify an image id or an image name. If you specify
both, the image id will be used.
- You can either specify a machine type id or a machine type name. If
you specify both, the machine type id will be used.
- You can either specify an VPC id or an VPC name. If you specify both,
the VPC id will be used.
- The user data must be a Base64 encoded string.

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

// any of: CreateParametersImage0, CreateParametersImage1
type CreateParametersImage struct {
	CreateParametersImage0 `json:",squash"` // nolint
	CreateParametersImage1 `json:",squash"` // nolint
}

type CreateParametersImage0 struct {
	Id string `json:"id"`
}

type CreateParametersImage1 struct {
	Name string `json:"name"`
}

// any of: CreateParametersImage0, CreateParametersImage1
type CreateParametersMachineType struct {
	CreateParametersImage0 `json:",squash"` // nolint
	CreateParametersImage1 `json:",squash"` // nolint
}

type CreateParametersNetwork struct {
	AssociatePublicIp *bool                       `json:"associate_public_ip,omitempty"`
	Nic               *CreateParametersNetworkNic `json:"nic,omitempty"`
	Vpc               *CreateParametersNetworkVpc `json:"vpc,omitempty"`
}

// any of: CreateParametersImage0, CreateParametersNetworkNic1
type CreateParametersNetworkNic struct {
	CreateParametersImage0      `json:",squash"` // nolint
	CreateParametersNetworkNic1 `json:",squash"` // nolint
}

type CreateParametersNetworkNic1 struct {
	SecurityGroups *CreateParametersNetworkNic1SecurityGroups `json:"security_groups,omitempty"`
}

type CreateParametersNetworkNic1SecurityGroupsItem struct {
	Id string `json:"id"`
}

type CreateParametersNetworkNic1SecurityGroups []CreateParametersNetworkNic1SecurityGroupsItem

// any of: CreateParametersImage0, CreateParametersImage1
type CreateParametersNetworkVpc struct {
	CreateParametersImage0 `json:",squash"` // nolint
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
