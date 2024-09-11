/*
Executor: list

# Summary

Lists all backups in the current tenant.

# Description

List Virtual Machine backups in the current tenant which is logged in.

#### Notes
- You can use the **extend** argument to get more details from the inner objects
like image or type.

Version: v1

import "magalu.cloud/lib/products/virtual_machine/backups"
*/
package backups

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
	Backups ListResultBackups `json:"backups"`
}

type ListResultBackupsItem struct {
	BackupType string                        `json:"backup_type"`
	CreatedAt  string                        `json:"created_at"`
	Id         string                        `json:"id"`
	Instance   ListResultBackupsItemInstance `json:"instance"`
	MinDisk    *int                          `json:"min_disk,omitempty"`
	Name       string                        `json:"name"`
	Size       *int                          `json:"size,omitempty"`
	State      string                        `json:"state"`
	Status     string                        `json:"status"`
	UpdatedAt  *string                       `json:"updated_at,omitempty"`
}

// any of: ListResultBackupsItemInstance
type ListResultBackupsItemInstance struct {
	Id          *string                                  `json:"id"`
	Image       ListResultBackupsItemInstanceImage       `json:"image"`
	MachineType ListResultBackupsItemInstanceMachineType `json:"machine_type"`
	Name        string                                   `json:"name"`
	State       string                                   `json:"state"`
	Status      string                                   `json:"status"`
}

// any of: ListResultBackupsItemInstanceImage
type ListResultBackupsItemInstanceImage struct {
	Id       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	Platform *string `json:"platform,omitempty"`
}

// any of: ListResultBackupsItemInstanceMachineType
type ListResultBackupsItemInstanceMachineType struct {
	Disk  *int    `json:"disk,omitempty"`
	Id    string  `json:"id"`
	Name  *string `json:"name,omitempty"`
	Ram   *int    `json:"ram,omitempty"`
	Vcpus *int    `json:"vcpus,omitempty"`
}

type ListResultBackups []ListResultBackupsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/backups/list"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ListContext(
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/backups/list"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
