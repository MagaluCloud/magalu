/*
Executor: create

# Summary

Create an instance.

# Description

# Create a Virtual Machine instance

Version: v1

import "github.com/MagaluCloud/magalu/mgc/lib/products/virtual_machine/instances"
*/
package instances

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type CreateParameters struct {
	AvailabilityZone *string                     `json:"availability_zone,omitempty"`
	Image            CreateParametersImage       `json:"image"`
	Labels           *CreateParametersLabels     `json:"labels,omitempty"`
	MachineType      CreateParametersMachineType `json:"machine_type"`
	Name             string                      `json:"name"`
	Network          *CreateParametersNetwork    `json:"network,omitempty"`
	SshKeyName       *string                     `json:"ssh_key_name,omitempty"`
	UserData         *string                     `json:"user_data,omitempty"`
}

// any of: CreateParametersImage
type CreateParametersImage struct {
	Id   string  `json:"id"`
	Name *string `json:"name,omitempty"`
}

type CreateParametersLabels []string

// any of: CreateParametersMachineType
type CreateParametersMachineType struct {
	Id   string  `json:"id"`
	Name *string `json:"name,omitempty"`
}

type CreateParametersNetwork struct {
	AssociatePublicIp *bool                             `json:"associate_public_ip,omitempty"`
	Interface         *CreateParametersNetworkInterface `json:"interface,omitempty"`
	Vpc               *CreateParametersNetworkVpc       `json:"vpc,omitempty"`
}

// any of: CreateParametersNetworkInterface
type CreateParametersNetworkInterface struct {
	Id             string                                          `json:"id"`
	Name           *string                                         `json:"name,omitempty"`
	SecurityGroups *CreateParametersNetworkInterfaceSecurityGroups `json:"security_groups,omitempty"`
}

type CreateParametersNetworkInterfaceSecurityGroupsItem struct {
	Id string `json:"id"`
}

type CreateParametersNetworkInterfaceSecurityGroups []CreateParametersNetworkInterfaceSecurityGroupsItem

// any of: CreateParametersNetworkVpc
type CreateParametersNetworkVpc struct {
	Id             string                                          `json:"id"`
	Name           *string                                         `json:"name,omitempty"`
	SecurityGroups *CreateParametersNetworkInterfaceSecurityGroups `json:"security_groups,omitempty"`
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) CreateContext(
	ctx context.Context,
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine/instances/create"), s.client, ctx)
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

	sdkConfig := s.client.Sdk().Config().TempConfig()
	if c["serverUrl"] == nil && sdkConfig["serverUrl"] != nil {
		c["serverUrl"] = sdkConfig["serverUrl"]
	}

	if c["env"] == nil && sdkConfig["env"] != nil {
		c["env"] = sdkConfig["env"]
	}

	if c["region"] == nil && sdkConfig["region"] != nil {
		c["region"] = sdkConfig["region"]
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
