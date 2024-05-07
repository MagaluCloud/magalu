/*
Executor: get

# Summary

Engine detail.

# Description

Returns a engine detail.

Version: 1.17.2

import "magalu.cloud/lib/products/dbaas/engines"
*/
package engines

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	EngineId string `json:"engine_id"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	Engine  string `json:"engine"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/dbaas/engines/get"), s.client, s.ctx)
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
