/*
Executor: get

# Summary

# Instance Details

# Description

# Returns a instance details

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

type GetParameters struct {
	Id string `json:"id"`
}

type GetConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	AvailabilityZone  *string                     `json:"availability_zone,omitempty"`
	CreatedAt         string                      `json:"created_at"`
	Error             *string                     `json:"error,omitempty"`
	Id                string                      `json:"id"`
	Image             GetResultImage              `json:"image"`
	InstanceId        *string                     `json:"instance_id,omitempty"`
	InstanceType      GetResultInstanceType       `json:"instance_type"`
	KeyName           *string                     `json:"key_name,omitempty"`
	Memory            *int                        `json:"memory,omitempty"`
	Name              *string                     `json:"name,omitempty"`
	NetworkInterfaces *GetResultNetworkInterfaces `json:"network_interfaces,omitempty"`
	Ports             *GetResultPorts             `json:"ports,omitempty"`
	PowerState        *int                        `json:"power_state,omitempty"`
	PowerStateLabel   *string                     `json:"power_state_label,omitempty"`
	RootStorage       *int                        `json:"root_storage,omitempty"`
	SecurityGroups    *GetResultSecurityGroups    `json:"security_groups,omitempty"`
	Status            string                      `json:"status"`
	UpdatedAt         *string                     `json:"updated_at,omitempty"`
	Vcpus             *int                        `json:"vcpus,omitempty"`
	Volumes           *GetResultVolumes           `json:"volumes,omitempty"`
}

type GetResultImage struct {
	Name string `json:"name"`
}

type GetResultInstanceType struct {
	Name string `json:"name"`
}

type GetResultNetworkInterfacesItem struct {
	Addresses  GetResultNetworkInterfacesItemAddresses `json:"addresses"`
	MacAddress *string                                 `json:"mac_address,omitempty"`
	Network    GetResultNetworkInterfacesItemNetwork   `json:"network,omitempty"`
}

type GetResultNetworkInterfacesItemAddressesItem struct {
	IpAddress *string `json:"ip_address,omitempty"`
	Type      *string `json:"type,omitempty"`
	Version   *int    `json:"version,omitempty"`
}

type GetResultNetworkInterfacesItemAddresses []GetResultNetworkInterfacesItemAddressesItem

type GetResultNetworkInterfacesItemNetwork struct {
	Name *string `json:"name,omitempty"`
}

type GetResultNetworkInterfaces []GetResultNetworkInterfacesItem

type GetResultPortsItem struct {
	Id string `json:"id"`
}

type GetResultPorts []GetResultPortsItem

type GetResultSecurityGroupsItem struct {
	Name *string `json:"name,omitempty"`
}

type GetResultSecurityGroups []GetResultSecurityGroupsItem

type GetResultVolumesItem struct {
	Id string `json:"id"`
}

type GetResultVolumes []GetResultVolumesItem

func Get(
	client *mgcClient.Client,
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine/instances/v0-instances/get"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[GetConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
