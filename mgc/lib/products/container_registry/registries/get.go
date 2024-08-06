/*
Executor: get

# Summary

# Get registry information

# Description

Show detailed information about the user's container registry.

Version: 0.1.0

import "magalu.cloud/lib/products/container_registry/registries"
*/
package registries

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	RegistryId string `json:"registry_id"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Container Registry's response data.
type GetResult struct {
	CreatedAt         string `json:"created_at"`
	Id                string `json:"id"`
	Name              string `json:"name"`
	StorageUsageBytes int    `json:"storage_usage_bytes"`
	UpdatedAt         string `json:"updated_at"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/container-registry/registries/get"), s.client, s.ctx)
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
