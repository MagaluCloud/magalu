/*
Executor: get

# Summary

Retrieve the details of an instance.

# Description

# Get a Virtual Machine instance details

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances"
*/
package instances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	Expand *GetParametersExpand `json:"expand,omitempty"`
	Id     string               `json:"id"`
}

type GetParametersExpand []string

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	AvailabilityZone *string              `json:"availability_zone,omitempty"`
	CreatedAt        string               `json:"created_at"`
	Error            *GetResultError      `json:"error,omitempty"`
	Id               string               `json:"id"`
	Image            GetResultImage       `json:"image"`
	MachineType      GetResultMachineType `json:"machine_type"`
	Name             *string              `json:"name,omitempty"`
	Network          *GetResultNetwork    `json:"network,omitempty"`
	SshKeyName       *string              `json:"ssh_key_name,omitempty"`
	State            string               `json:"state"`
	Status           string               `json:"status"`
	UpdatedAt        *string              `json:"updated_at,omitempty"`
	UserData         *string              `json:"user_data,omitempty"`
}

type GetResultError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// any of: GetResultImage
type GetResultImage struct {
	Id       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	Platform *string `json:"platform,omitempty"`
}

// any of: GetResultMachineType
type GetResultMachineType struct {
	Disk  *int    `json:"disk,omitempty"`
	Id    string  `json:"id"`
	Name  *string `json:"name,omitempty"`
	Ram   *int    `json:"ram,omitempty"`
	Vcpus *int    `json:"vcpus,omitempty"`
}

// any of: GetResultNetwork
type GetResultNetwork struct {
	Ports *GetResultNetworkPorts `json:"ports"`
	Vpc   *GetResultNetworkVpc   `json:"vpc,omitempty"`
}

type GetResultNetworkPortsItem struct {
	Id          string                               `json:"id"`
	IpAddresses GetResultNetworkPortsItemIpAddresses `json:"ipAddresses"`
	Name        string                               `json:"name"`
}

type GetResultNetworkPortsItemIpAddresses struct {
	IpV6address      *string `json:"ipV6Address,omitempty"`
	PrivateIpAddress string  `json:"privateIpAddress"`
	PublicIpAddress  *string `json:"publicIpAddress,omitempty"`
}

type GetResultNetworkPorts []GetResultNetworkPortsItem

type GetResultNetworkVpc struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine/instances/get"), s.client, s.ctx)
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

func (s *service) GetUntilTermination(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	e, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine/instances/get"), s.client, s.ctx)
	if err != nil {
		return
	}

	exec, ok := e.(mgcCore.TerminatorExecutor)
	if !ok {
		// Not expected, but let's fallback
		return s.Get(
			parameters,
			configs,
		)
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[GetConfigs](configs); err != nil {
		return
	}

	r, err := exec.ExecuteUntilTermination(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
