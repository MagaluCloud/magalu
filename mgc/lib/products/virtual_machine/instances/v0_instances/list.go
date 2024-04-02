/*
Executor: list

# Summary

# List Instances

# Description

Returns a list of instances for a provided tenant_id

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

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Instances ListResultInstances `json:"instances"`
}

type ListResultInstancesItem struct {
	AvailabilityZone  *string                                   `json:"availability_zone,omitempty"`
	CreatedAt         string                                    `json:"created_at"`
	Error             *string                                   `json:"error,omitempty"`
	Id                string                                    `json:"id"`
	Image             ListResultInstancesItemImage              `json:"image"`
	InstanceId        *string                                   `json:"instance_id,omitempty"`
	InstanceType      ListResultInstancesItemInstanceType       `json:"instance_type"`
	KeyName           *string                                   `json:"key_name,omitempty"`
	Memory            *int                                      `json:"memory,omitempty"`
	Name              *string                                   `json:"name,omitempty"`
	NetworkInterfaces *ListResultInstancesItemNetworkInterfaces `json:"network_interfaces,omitempty"`
	Ports             *ListResultInstancesItemPorts             `json:"ports,omitempty"`
	PowerState        *int                                      `json:"power_state,omitempty"`
	PowerStateLabel   *string                                   `json:"power_state_label,omitempty"`
	RootStorage       *int                                      `json:"root_storage,omitempty"`
	SecurityGroups    *ListResultInstancesItemSecurityGroups    `json:"security_groups,omitempty"`
	Status            string                                    `json:"status"`
	UpdatedAt         *string                                   `json:"updated_at,omitempty"`
	Vcpus             *int                                      `json:"vcpus,omitempty"`
	Volumes           *ListResultInstancesItemVolumes           `json:"volumes,omitempty"`
}

type ListResultInstancesItemImage struct {
	Name string `json:"name"`
}

type ListResultInstancesItemInstanceType struct {
	Name string `json:"name"`
}

type ListResultInstancesItemNetworkInterfacesItem struct {
	Addresses  ListResultInstancesItemNetworkInterfacesItemAddresses `json:"addresses"`
	MacAddress *string                                               `json:"mac_address,omitempty"`
	Network    ListResultInstancesItemNetworkInterfacesItemNetwork   `json:"network,omitempty"`
}

type ListResultInstancesItemNetworkInterfacesItemAddressesItem struct {
	IpAddress *string `json:"ip_address,omitempty"`
	Type      *string `json:"type,omitempty"`
	Version   *int    `json:"version,omitempty"`
}

type ListResultInstancesItemNetworkInterfacesItemAddresses []ListResultInstancesItemNetworkInterfacesItemAddressesItem

type ListResultInstancesItemNetworkInterfacesItemNetwork struct {
	Name *string `json:"name,omitempty"`
}

type ListResultInstancesItemNetworkInterfaces []ListResultInstancesItemNetworkInterfacesItem

type ListResultInstancesItemPortsItem struct {
	Id string `json:"id"`
}

type ListResultInstancesItemPorts []ListResultInstancesItemPortsItem

type ListResultInstancesItemSecurityGroupsItem struct {
	Name *string `json:"name,omitempty"`
}

type ListResultInstancesItemSecurityGroups []ListResultInstancesItemSecurityGroupsItem

type ListResultInstancesItemVolumesItem struct {
	Id string `json:"id"`
}

type ListResultInstancesItemVolumes []ListResultInstancesItemVolumesItem

type ListResultInstances []ListResultInstancesItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/instances/v0-instances/list"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ListConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
