/*
Executor: restores

# Summary

Backup restore.

# Description

Restores a backup for an instance asynchronously.

Version: 1.17.2

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RestoresParameters struct {
	BackupId   string `json:"backup_id"`
	Exchange   string `json:"exchange,omitempty"`
	InstanceId string `json:"instance_id"`
}

type RestoresConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type RestoresResult struct {
	Id string `json:"id"`
}

func Restores(
	client *mgcClient.Client,
	ctx context.Context,
	parameters RestoresParameters,
	configs RestoresConfigs,
) (
	result RestoresResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Restores", mgcCore.RefPath("/dbaas/instances/restores"), client, ctx)
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

// TODO: links
// TODO: related
