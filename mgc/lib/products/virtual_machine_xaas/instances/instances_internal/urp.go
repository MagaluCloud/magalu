/*
Executor: urp

# Summary

# Instance Internal By Urp Detail

# Description

This route is to get a detailed information for a instance but adding the urp id on the response.

### Note
This route is used only for internal proposes.

Version: 1.255.1

import "magalu.cloud/lib/products/virtual_machine_xaas/instances/instances_internal"
*/
package instancesInternal

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UrpParameters struct {
	Id          string  `json:"id"`
	ProjectType *string `json:"project_type,omitempty"`
}

type UrpConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type UrpResult struct {
	AvailabilityZone *string              `json:"availability_zone,omitempty"`
	CreatedAt        string               `json:"created_at"`
	Id               string               `json:"id"`
	Image            UrpResultImage       `json:"image"`
	InstanceId       *string              `json:"instance_id,omitempty"`
	KeyName          *string              `json:"key_name,omitempty"`
	MachineType      UrpResultMachineType `json:"machine_type"`
	Name             *string              `json:"name,omitempty"`
	Network          *UrpResultNetwork    `json:"network,omitempty"`
	State            string               `json:"state"`
	Status           string               `json:"status"`
	UpdatedAt        *string              `json:"updated_at,omitempty"`
	UserData         *string              `json:"user_data,omitempty"`
}

// any of: UrpResultImage
type UrpResultImage struct {
	Id       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	Platform *string `json:"platform,omitempty"`
}

// any of: UrpResultMachineType
type UrpResultMachineType struct {
	Disk  *int    `json:"disk,omitempty"`
	Id    string  `json:"id"`
	Name  *string `json:"name,omitempty"`
	Ram   *int    `json:"ram,omitempty"`
	Vcpus *int    `json:"vcpus,omitempty"`
}

type UrpResultNetwork struct {
	Ports UrpResultNetworkPorts `json:"ports"`
}

type UrpResultNetworkPortsItem struct {
	Id string `json:"id"`
}

type UrpResultNetworkPorts []UrpResultNetworkPortsItem

func (s *service) Urp(
	parameters UrpParameters,
	configs UrpConfigs,
) (
	result UrpResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Urp", mgcCore.RefPath("/virtual-machine-xaas/instances/instances-internal/urp"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UrpParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UrpConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[UrpResult](r)
}

// TODO: links
// TODO: related
