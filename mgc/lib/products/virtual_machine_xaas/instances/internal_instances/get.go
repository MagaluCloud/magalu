/*
Executor: get

# Summary

# Instance Internal Detail

# Description

This route is to get a detailed information for a instance but adding the urp id on the response.

### Note
This route is used only for internal proposes.

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/instances/internal_instances"
*/
package internalInstances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	Id          string  `json:"id"`
	ProjectType *string `json:"project_type,omitempty"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	AvailabilityZone *string              `json:"availability_zone,omitempty"`
	CreatedAt        string               `json:"created_at"`
	Id               string               `json:"id"`
	Image            GetResultImage       `json:"image"`
	InstanceId       *string              `json:"instance_id,omitempty"`
	KeyName          *string              `json:"key_name,omitempty"`
	MachineType      GetResultMachineType `json:"machine_type"`
	Name             *string              `json:"name,omitempty"`
	Network          *GetResultNetwork    `json:"network,omitempty"`
	State            string               `json:"state"`
	Status           string               `json:"status"`
	UpdatedAt        *string              `json:"updated_at,omitempty"`
	UserData         *string              `json:"user_data,omitempty"`
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

type GetResultNetwork struct {
	Ports GetResultNetworkPorts `json:"ports"`
}

type GetResultNetworkPortsItem struct {
	Id string `json:"id"`
}

type GetResultNetworkPorts []GetResultNetworkPortsItem

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine-xaas/instances/internal-instances/get"), s.client, s.ctx)
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
