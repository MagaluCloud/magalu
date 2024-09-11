/*
Executor: get

# Summary

Retrieve the details of a specific backup.

# Description

Retrieve details of a Backup for the currently authenticated tenant.

#### Notes
  - Use the **expand** argument to obtain additional details about the Volume
    used to create the Backup.
  - Utilize the **block-storage backups list** command to retrieve a list of
    all Backups and obtain the ID of the Backup for which you want to retrieve
    details.

Version: v1

import "magalu.cloud/lib/products/block_storage/backups"
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
	CreatedAt    string                 `json:"created_at"`
	Description  *string                `json:"description,omitempty"`
	Error        *GetResultError        `json:"error,omitempty"`
	Id           string                 `json:"id"`
	Name         string                 `json:"name"`
	Size         int                    `json:"size"`
	SourceBackup *GetResultSourceBackup `json:"source_backup,omitempty"`
	State        string                 `json:"state"`
	Status       string                 `json:"status"`
	Type         string                 `json:"type"`
	UpdatedAt    string                 `json:"updated_at"`
	Volume       *GetResultVolume       `json:"volume"`
}

type GetResultError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type GetResultSourceBackup struct {
	Id string `json:"id"`
}

// any of: GetResultVolume
type GetResultVolume struct {
	CreatedAt  string                    `json:"created_at"`
	Id         string                    `json:"id"`
	Name       string                    `json:"name"`
	Size       int                       `json:"size"`
	State      string                    `json:"state"`
	Status     string                    `json:"status"`
	UpdatedAt  string                    `json:"updated_at"`
	VolumeType GetResultVolumeVolumeType `json:"volume_type"`
}

type GetResultVolumeVolumeType struct {
	Id string `json:"id"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/block-storage/backups/get"), s.client, s.ctx)
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/block-storage/backups/get"), s.client, ctx)
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
