/*
Executor: get

# Summary

Retrieve the details of a backup.

# Description

Get a backup details for the current tenant which is logged in.

#### Notes
- You can use the backup list command to retrieve all backups,
so you can get the id of the backup that you want to get details.

- You can use the **expand** argument to get more details from the object
like instance.

Version: 1.249.0

import "magalu.cloud/lib/products/virtual_machine_xaas/xaas_backups"
*/
package xaasBackups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	Expand      *GetParametersExpand `json:"expand,omitempty"`
	Id          string               `json:"id"`
	ProjectType string               `json:"project_type"`
}

type GetParametersExpand []string

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	BackupType string            `json:"backup_type"`
	CreatedAt  string            `json:"created_at"`
	Id         string            `json:"id"`
	Instance   GetResultInstance `json:"instance"`
	MinDisk    *int              `json:"min_disk,omitempty"`
	Name       string            `json:"name"`
	Size       *int              `json:"size,omitempty"`
	State      string            `json:"state"`
	Status     string            `json:"status"`
	UpdatedAt  *string           `json:"updated_at,omitempty"`
}

// any of: GetResultInstance
type GetResultInstance struct {
	Id          string                       `json:"id"`
	Image       GetResultInstanceImage       `json:"image"`
	MachineType GetResultInstanceMachineType `json:"machine_type"`
	Name        string                       `json:"name"`
	State       string                       `json:"state"`
	Status      string                       `json:"status"`
}

// any of: GetResultInstanceImage
type GetResultInstanceImage struct {
	Id       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	Platform *string `json:"platform,omitempty"`
}

// any of: GetResultInstanceMachineType
type GetResultInstanceMachineType struct {
	Disk  *int    `json:"disk,omitempty"`
	Id    string  `json:"id"`
	Name  *string `json:"name,omitempty"`
	Ram   *int    `json:"ram,omitempty"`
	Vcpus *int    `json:"vcpus,omitempty"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine-xaas/xaas backups/get"), s.client, s.ctx)
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
