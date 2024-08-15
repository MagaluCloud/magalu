/*
Executor: get

# Summary

# Get Ssh Key

Version: 0.1.0

import "magalu.cloud/lib/products/ssh/ssh_keys"
*/
package sshKeys

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	KeyId string `json:"key_id"`
}

type GetConfigs struct {
	XTenantId string  `json:"X-Tenant-ID"`
	Env       *string `json:"env,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	KeyType string `json:"key_type"`
	Name    string `json:"name"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/ssh/ssh_keys/get"), s.client, s.ctx)
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
