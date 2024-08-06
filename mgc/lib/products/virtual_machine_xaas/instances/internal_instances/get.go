/*
Executor: get

# Summary

# Instance Internal Detail

# Description

This route is to get a detailed information for a instance but adding the urp id on the response.

### Note
This route is used only for internal proposes.

Version: 1.258.0

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
