/*
Executor: get

# Summary

# Get Instance Diagnostic

# Description

Get more internal information about a instance informing a
instance id and the instance`s tenant.

This route will get DB info and URP info merged in a response.

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/diagnostics"
*/
package diagnostics

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	Id string `json:"id"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	Instance GetResultInstance `json:"instance"`
}

type GetResultInstance struct {
	Database GetResultInstanceDatabase `json:"database"`
	Urp      *GetResultInstanceUrp     `json:"urp,omitempty"`
}

type GetResultInstanceDatabase struct {
	CreatedAt   string                         `json:"created_at"`
	Error       *string                        `json:"error,omitempty"`
	Id          string                         `json:"id"`
	ImageName   string                         `json:"image_name"`
	KeyPairName string                         `json:"key_pair_name"`
	Name        string                         `json:"name"`
	Ports       GetResultInstanceDatabasePorts `json:"ports"`
	ProjectType string                         `json:"project_type"`
	Retries     int                            `json:"retries"`
	SnapshotId  *string                        `json:"snapshot_id,omitempty"`
	Status      string                         `json:"status"`
	Step        int                            `json:"step"`
	UpdatedAt   *string                        `json:"updated_at,omitempty"`
	UserData    *string                        `json:"user_data,omitempty"`
}

type GetResultInstanceDatabasePortsItem struct {
	Port string `json:"port"`
}

type GetResultInstanceDatabasePorts []GetResultInstanceDatabasePortsItem

type GetResultInstanceUrp struct {
	Addresses        GetResultInstanceUrpAddresses      `json:"addresses"`
	AvailabilityZone string                             `json:"availability_zone"`
	CreatedAt        string                             `json:"created_at"`
	Flavor           GetResultInstanceUrpFlavor         `json:"flavor"`
	HostId           string                             `json:"host_id"`
	Id               string                             `json:"id"`
	Image            GetResultInstanceUrpImage          `json:"image"`
	KeyName          string                             `json:"key_name"`
	Name             string                             `json:"name"`
	PowerState       int                                `json:"power_state"`
	SecurityGroups   GetResultInstanceUrpSecurityGroups `json:"security_groups"`
	Status           string                             `json:"status"`
	UpdatedAt        string                             `json:"updated_at"`
	VmState          string                             `json:"vm_state"`
}

type GetResultInstanceUrpAddresses struct {
}

type GetResultInstanceUrpFlavor struct {
	Disk  int    `json:"disk"`
	Id    string `json:"id"`
	Name  string `json:"name"`
	Ram   int    `json:"ram"`
	Swap  int    `json:"swap"`
	Vcpus int    `json:"vcpus"`
}

type GetResultInstanceUrpImage struct {
	Id string `json:"id"`
}

type GetResultInstanceUrpSecurityGroupsItem struct {
	Name string `json:"name"`
}

type GetResultInstanceUrpSecurityGroups []GetResultInstanceUrpSecurityGroupsItem

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine-xaas/diagnostics/get"), s.client, s.ctx)
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
