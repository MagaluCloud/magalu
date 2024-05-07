/*
Executor: create

# Summary

Backup Create.

# Description

Creates a new backup asynchronously.

Version: 1.17.2

import "magalu.cloud/lib/products/dbaas/instances/backups"
*/
package backups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Exchange   *string `json:"exchange,omitempty"`
	InstanceId string  `json:"instance_id"`
	Mode       string  `json:"mode"`
}

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
	Id string `json:"id"`
}

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/dbaas/instances/backups/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
