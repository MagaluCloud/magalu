/*
Executor: list

# Summary

List all instances.

# Description

# List Virtual Machine instances

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances"
*/
package instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit  *int                  `json:"_limit,omitempty"`
	Offset *int                  `json:"_offset,omitempty"`
	Sort   *string               `json:"_sort,omitempty"`
	Expand *ListParametersExpand `json:"expand,omitempty"`
}

type ListParametersExpand []string

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Instances ListResultInstances `json:"instances"`
}

type ListResultInstancesItem struct {
	AvailabilityZone *string                            `json:"availability_zone,omitempty"`
	CreatedAt        string                             `json:"created_at"`
	Error            *ListResultInstancesItemError      `json:"error,omitempty"`
	Id               string                             `json:"id"`
	Image            ListResultInstancesItemImage       `json:"image"`
	MachineType      ListResultInstancesItemMachineType `json:"machine_type"`
	Name             *string                            `json:"name,omitempty"`
	Network          *ListResultInstancesItemNetwork    `json:"network,omitempty"`
	SshKeyName       *string                            `json:"ssh_key_name,omitempty"`
	State            string                             `json:"state"`
	Status           string                             `json:"status"`
	UpdatedAt        *string                            `json:"updated_at,omitempty"`
	UserData         *string                            `json:"user_data,omitempty"`
}

type ListResultInstancesItemError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// any of: ListResultInstancesItemImage
type ListResultInstancesItemImage struct {
	Id       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	Platform *string `json:"platform,omitempty"`
}

// any of: ListResultInstancesItemMachineType
type ListResultInstancesItemMachineType struct {
	Disk  *int    `json:"disk,omitempty"`
	Id    string  `json:"id"`
	Name  *string `json:"name,omitempty"`
	Ram   *int    `json:"ram,omitempty"`
	Vcpus *int    `json:"vcpus,omitempty"`
}

// any of: ListResultInstancesItemNetwork
type ListResultInstancesItemNetwork struct {
	Ports *ListResultInstancesItemNetworkPorts `json:"ports"`
	Vpc   *ListResultInstancesItemNetworkVpc   `json:"vpc,omitempty"`
}

type ListResultInstancesItemNetworkPortsItem struct {
	Id          string                                             `json:"id"`
	IpAddresses ListResultInstancesItemNetworkPortsItemIpAddresses `json:"ipAddresses"`
	Name        string                                             `json:"name"`
}

type ListResultInstancesItemNetworkPortsItemIpAddresses struct {
	IpV6address      *string `json:"ipV6Address,omitempty"`
	PrivateIpAddress string  `json:"privateIpAddress"`
	PublicIpAddress  *string `json:"publicIpAddress,omitempty"`
}

type ListResultInstancesItemNetworkPorts []ListResultInstancesItemNetworkPortsItem

type ListResultInstancesItemNetworkVpc struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ListResultInstances []ListResultInstancesItem

/*func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/instances/list"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ListConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}*/

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ListContext(
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/instances/list"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListParameters](parameters); err != nil {
		return
	}

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
