/*
Executor: list

# Summary

List available datastores (Deprecated).

# Description

**Deprecated**: This endpoint is being phased out. Please use the `/v1/engines` endpoint to retrieve the list of available engines instead.

Returns a list of available datastores. It is recommended to update your integration to use the newer `/v1/engines` endpoint for improved functionality and future compatibility.

Version: 1.26.1

import "magalu.cloud/lib/products/dbaas/datastores"
*/
package datastores

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit  *int    `json:"_limit,omitempty"`
	Offset *int    `json:"_offset,omitempty"`
	Status *string `json:"status,omitempty"`
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
	Engine  string `json:"engine"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

type ListResultResults []ListResultResultsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/dbaas/datastores/list"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[map[string]interface{}](s.client.Sdk().Config().TempConfig()); err != nil {
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
