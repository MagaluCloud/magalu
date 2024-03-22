/*
Executor: create

# Summary

Create an instance asynchronously.

# Description

# Create a Virtual Machine instance

Version: 0.1.0

import "magalu.cloud/lib/products/virtual_machine/instances"
*/
package instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	AvailabilityZone *string                     `json:"availability_zone,omitempty"`
	Image            CreateParametersImage       `json:"image"`
	MachineType      CreateParametersMachineType `json:"machine_type"`
	Name             string                      `json:"name"`
	Network          CreateParametersNetwork     `json:"network,omitempty"`
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
	AssociatePublicIp bool                       `json:"associate_public_ip,omitempty"`
	Vpc               CreateParametersNetworkVpc `json:"vpc,omitempty"`
}

// any of: CreateParametersImage0, CreateParametersImage1
type CreateParametersNetworkVpc struct {
	CreateParametersImage0 `json:",squash"` // nolint
	CreateParametersImage1 `json:",squash"` // nolint
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine/instances/create"), client, ctx)
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