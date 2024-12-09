/*
Executor: get

# Summary

Backup Detail.

# Description

Get a backup detail.

Version: 1.34.1

import "magalu.cloud/lib/products/dbaas/instances/backups"
*/
package backups

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	BackupId   string `json:"backup_id"`
	InstanceId string `json:"instance_id"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	CreatedAt  string             `json:"created_at"`
	DbSize     *int               `json:"db_size,omitempty"`
	EngineId   string             `json:"engine_id"`
	FinishedAt *string            `json:"finished_at,omitempty"`
	Id         string             `json:"id"`
	Instance   *GetResultInstance `json:"instance,omitempty"`
	InstanceId string             `json:"instance_id"`
	Location   *string            `json:"location,omitempty"`
	Mode       string             `json:"mode"`
	Name       *string            `json:"name,omitempty"`
	Size       *int               `json:"size,omitempty"`
	StartedAt  *string            `json:"started_at,omitempty"`
	Status     string             `json:"status"`
	Type       string             `json:"type"`
	UpdatedAt  *string            `json:"updated_at,omitempty"`
}

// This response object provides details about a database instance associated with a backup.  It is provided only if the originating database instance of the backup is not deleted.  If the originating instance is deleted, no instance details will be provided.

type GetResultInstance struct {
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/dbaas/instances/backups/get"), s.client, s.ctx)
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/dbaas/instances/backups/get"), s.client, ctx)
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
