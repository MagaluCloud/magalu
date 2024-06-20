/*
Executor: list

# Summary

# List all container registries

# Description

List user's container registries.

Version: 0.1.0

import "magalu.cloud/lib/products/container_registry/registries/registries"
*/
package registries

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit  *int    `json:"_limit,omitempty"`
	Offset *int    `json:"_offset,omitempty"`
	Sort   *string `json:"_sort,omitempty"`
}

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Container registry information response object.
type ListResult struct {
	Results ListResultResults `json:"results"`
}

// Container Registry's response data.
type ListResultResultsItem struct {
	CreatedAt         string `json:"created_at"`
	Id                string `json:"id"`
	Name              string `json:"name"`
	StorageUsageBytes int    `json:"storage_usage_bytes"`
	UpdatedAt         string `json:"updated_at"`
}

type ListResultResults []ListResultResultsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/container-registry/registries/registries/list"), s.client, s.ctx)
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

// TODO: links
// TODO: related
