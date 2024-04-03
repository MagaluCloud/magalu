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
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	Expand GetParametersExpand `json:"expand,omitempty"`
	Id     string              `json:"id"`
}

type GetParametersExpand []string

type GetConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	AvailabilityZone *string              `json:"availability_zone,omitempty"`
	CreatedAt        string               `json:"created_at"`
	Error            GetResultError       `json:"error,omitempty"`
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

// any of: GetResultImage0, GetResultImage1
type GetResultImage struct {
	GetResultImage0 `json:",squash"` // nolint
	GetResultImage1 `json:",squash"` // nolint
}

type GetResultImage0 struct {
	Id string `json:"id"`
}

type GetResultImage1 struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Platform *string `json:"platform,omitempty"`
}

// any of: GetResultMachineType0, GetResultMachineType1
type GetResultMachineType struct {
	GetResultMachineType0 `json:",squash"` // nolint
	GetResultMachineType1 `json:",squash"` // nolint
}

type GetResultMachineType0 struct {
	Id string `json:"id"`
}

type GetResultMachineType1 struct {
	Disk  int    `json:"disk"`
	Id    string `json:"id"`
	Name  string `json:"name"`
	Ram   int    `json:"ram"`
	Vcpus int    `json:"vcpus"`
}

// any of: GetResultNetwork0, GetResultNetwork1
type GetResultNetwork struct {
	GetResultNetwork0 `json:",squash"` // nolint
	GetResultNetwork1 `json:",squash"` // nolint
}

type GetResultNetwork0 struct {
	Ports GetResultNetwork0Ports `json:"ports"`
}

type GetResultNetwork0PortsItem struct {
	Id string `json:"id"`
}

type GetResultNetwork0Ports []GetResultNetwork0PortsItem

type GetResultNetwork1 struct {
	Ports *GetResultNetwork1Ports `json:"ports,omitempty"`
	Vpc   GetResultNetwork1Vpc    `json:"vpc,omitempty"`
}

type GetResultNetwork1PortsItem struct {
	Id          string                                `json:"id"`
	IpAddresses GetResultNetwork1PortsItemIpAddresses `json:"ipAddresses"`
	Name        string                                `json:"name"`
}

type GetResultNetwork1PortsItemIpAddresses struct {
	PrivateIpAddress string `json:"privateIpAddress"`
	PublicIpAddress  string `json:"publicIpAddress"`
}

type GetResultNetwork1Ports []GetResultNetwork1PortsItem

type GetResultNetwork1Vpc struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func Get(
	client *mgcClient.Client,
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine/instances/get"), client, ctx)
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

func GetUntilTermination(
	client *mgcClient.Client,
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	e, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine/instances/get"), client, ctx)
	if err != nil {
		return
	}

	exec, ok := e.(mgcCore.TerminatorExecutor)
	if !ok {
		// Not expected, but let's fallback
		return Get(
			client,
			ctx,
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
