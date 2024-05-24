/*
Executor: list

# Summary

Lists all instances in the current tenant.

# Description

# List Virtual Machine instances

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances"
*/
package instances

import (
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

// any of: ListResultInstancesItemImage0, ListResultInstancesItemImage1
type ListResultInstancesItemImage struct {
	ListResultInstancesItemImage0 `json:",squash"` // nolint
	ListResultInstancesItemImage1 `json:",squash"` // nolint
}

type ListResultInstancesItemImage0 struct {
	Id string `json:"id"`
}

type ListResultInstancesItemImage1 struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Platform *string `json:"platform,omitempty"`
}

// any of: ListResultInstancesItemMachineType0, ListResultInstancesItemMachineType1
type ListResultInstancesItemMachineType struct {
	ListResultInstancesItemMachineType0 `json:",squash"` // nolint
	ListResultInstancesItemMachineType1 `json:",squash"` // nolint
}

type ListResultInstancesItemMachineType0 struct {
	Id string `json:"id"`
}

type ListResultInstancesItemMachineType1 struct {
	Disk  int    `json:"disk"`
	Id    string `json:"id"`
	Name  string `json:"name"`
	Ram   int    `json:"ram"`
	Vcpus int    `json:"vcpus"`
}

// any of: ListResultInstancesItemNetwork0, ListResultInstancesItemNetwork1
type ListResultInstancesItemNetwork struct {
	ListResultInstancesItemNetwork0 `json:",squash"` // nolint
	ListResultInstancesItemNetwork1 `json:",squash"` // nolint
}

type ListResultInstancesItemNetwork0 struct {
	Ports ListResultInstancesItemNetwork0Ports `json:"ports"`
}

type ListResultInstancesItemNetwork0PortsItem struct {
	Id string `json:"id"`
}

type ListResultInstancesItemNetwork0Ports []ListResultInstancesItemNetwork0PortsItem

type ListResultInstancesItemNetwork1 struct {
	Ports *ListResultInstancesItemNetwork1Ports `json:"ports,omitempty"`
	Vpc   *ListResultInstancesItemNetwork1Vpc   `json:"vpc,omitempty"`
}

type ListResultInstancesItemNetwork1PortsItem struct {
	Id          string                                              `json:"id"`
	IpAddresses ListResultInstancesItemNetwork1PortsItemIpAddresses `json:"ipAddresses"`
	Name        string                                              `json:"name"`
}

type ListResultInstancesItemNetwork1PortsItemIpAddresses struct {
	IpV6address      *string `json:"ipV6Address,omitempty"`
	PrivateIpAddress string  `json:"privateIpAddress"`
	PublicIpAddress  *string `json:"publicIpAddress,omitempty"`
}

type ListResultInstancesItemNetwork1Ports []ListResultInstancesItemNetwork1PortsItem

type ListResultInstancesItemNetwork1Vpc struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ListResultInstances []ListResultInstancesItem

func (s *service) List(
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
}

// TODO: links
// TODO: related
