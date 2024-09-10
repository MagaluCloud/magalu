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

Version: v1

import "magalu.cloud/lib/products/virtual_machine/backups"
*/
package backups

import (
	"context"

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
	Id          *string                      `json:"id"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine/backups/get"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) GetContext(
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine/backups/get"), s.client, ctx)
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

	sdkConfig := s.client.Sdk().Config().TempConfig()
	if c["serverUrl"] == nil && sdkConfig["serverUrl"] != nil {
		c["serverUrl"] = sdkConfig["serverUrl"]
	}

	if c["env"] == nil && sdkConfig["env"] != nil {
		c["env"] = sdkConfig["env"]
	}

	if c["region"] == nil && sdkConfig["region"] != nil {
		c["region"] = sdkConfig["region"]
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
