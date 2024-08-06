/*
Executor: get

# Summary

Get a container registry repository by repository_name

# Description

Return detailed repository's information filtered by name.

Version: 0.1.0

import "magalu.cloud/lib/products/container_registry/repositories"
*/
package repositories

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	RegistryId     string `json:"registry_id"`
	RepositoryName string `json:"repository_name"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Information about the repository.
type GetResult struct {
	CreatedAt    string `json:"created_at"`
	ImageCount   int    `json:"image_count"`
	Name         string `json:"name"`
	RegistryName string `json:"registry_name"`
	UpdatedAt    string `json:"updated_at"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/container-registry/repositories/get"), s.client, s.ctx)
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
