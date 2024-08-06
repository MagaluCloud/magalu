/*
Executor: list

# Summary

# List all container registry repositories

# Description

List all user's repositories in the container registry.

Version: 0.1.0

import "magalu.cloud/lib/products/container_registry/repositories"
*/
package repositories

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit      *int    `json:"_limit,omitempty"`
	Offset     *int    `json:"_offset,omitempty"`
	Sort       *string `json:"_sort,omitempty"`
	RegistryId string  `json:"registry_id"`
}

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Information returned about the container registry repository.
type ListResult struct {
	Goal    *ListResultGoal   `json:"goal,omitempty"`
	Results ListResultResults `json:"results"`
}

// User's repositories quantity.
type ListResultGoal struct {
	Total *int `json:"total,omitempty"`
}

// Information about the repository.
type ListResultResultsItem struct {
	CreatedAt    string `json:"created_at"`
	ImageCount   int    `json:"image_count"`
	Name         string `json:"name"`
	RegistryName string `json:"registry_name"`
	UpdatedAt    string `json:"updated_at"`
}

type ListResultResults []ListResultResultsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/container-registry/repositories/list"), s.client, s.ctx)
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
