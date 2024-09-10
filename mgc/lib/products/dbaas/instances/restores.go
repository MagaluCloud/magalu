/*
Executor: restores

# Summary

Backup restore.

# Description

Restores a backup for an instance asynchronously.

Version: 1.26.1

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RestoresParameters struct {
	BackupId   string `json:"backup_id"`
	InstanceId string `json:"instance_id"`
}

type RestoresConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type RestoresResult struct {
	Id string `json:"id"`
}

func (s *service) Restores(
	parameters RestoresParameters,
	configs RestoresConfigs,
) (
	result RestoresResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Restores", mgcCore.RefPath("/dbaas/instances/restores"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RestoresParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[RestoresConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[RestoresResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) RestoresContext(
	ctx context.Context,
	parameters RestoresParameters,
	configs RestoresConfigs,
) (
	result RestoresResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Restores", mgcCore.RefPath("/dbaas/instances/restores"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RestoresParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[RestoresConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[RestoresResult](r)
}

// TODO: links
// TODO: related
