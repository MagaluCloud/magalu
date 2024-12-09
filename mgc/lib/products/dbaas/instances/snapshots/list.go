/*
Executor: list

# Summary

Snapshots List.

# Description

List all snapshots.

Version: 1.34.1

import "magalu.cloud/lib/products/dbaas/instances/snapshots"
*/
package snapshots

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit      *int    `json:"_limit,omitempty"`
	Offset     *int    `json:"_offset,omitempty"`
	InstanceId string  `json:"instance_id"`
	Status     *string `json:"status,omitempty"`
	Type       *string `json:"type,omitempty"`
}

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Meta    ListResultMeta    `json:"meta"`
	Results ListResultResults `json:"results"`
}

// Page details about the current request pagination.
type ListResultMeta struct {
	Filters ListResultMetaFilters `json:"filters"`
	Page    ListResultMetaPage    `json:"page"`
}

type ListResultMetaFiltersItem struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type ListResultMetaFilters []ListResultMetaFiltersItem

type ListResultMetaPage struct {
	Count    int `json:"count"`
	Limit    int `json:"limit"`
	MaxLimit int `json:"max_limit"`
	Offset   int `json:"offset"`
	Total    int `json:"total"`
}

type ListResultResultsItem struct {
	AllocatedSize int                           `json:"allocated_size"`
	CreatedAt     string                        `json:"created_at"`
	Description   string                        `json:"description"`
	FinishedAt    *string                       `json:"finished_at,omitempty"`
	Id            string                        `json:"id"`
	Instance      ListResultResultsItemInstance `json:"instance"`
	Name          string                        `json:"name"`
	StartedAt     *string                       `json:"started_at,omitempty"`
	Status        string                        `json:"status"`
	Type          string                        `json:"type"`
	UpdatedAt     *string                       `json:"updated_at,omitempty"`
}

// This response object provides details about a database instance associated with a snapshot.

type ListResultResultsItemInstance struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ListResultResults []ListResultResultsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/dbaas/instances/snapshots/list"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ListContext(
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/dbaas/instances/snapshots/list"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
