/*
Executor: create

# Summary

# Instance Create

# Description

# Create a instance asynchronously

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/instances/v0_instances"
*/
package v0Instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	AllocateFip       bool                               `json:"allocate_fip,omitempty"`
	AvailabilityZone  *string                            `json:"availability_zone,omitempty"`
	Image             *string                            `json:"image,omitempty"`
	KeyName           string                             `json:"key_name"`
	Name              string                             `json:"name"`
	NetworkInterfaces *CreateParametersNetworkInterfaces `json:"network_interfaces,omitempty"`
	Type              string                             `json:"type"`
	UserData          *string                            `json:"user_data,omitempty"`
}

type CreateParametersNetworkInterfacesItem struct {
	Id string `json:"id"`
}

type CreateParametersNetworkInterfaces []CreateParametersNetworkInterfacesItem

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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine/instances/v0-instances/create"), client, ctx)
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
