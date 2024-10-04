/*
Executor: list

# Summary

# List Subnet Pools by Tenant

# Description

# Returns a list of Subnet Pools for the current tenant's project

Version: 1.138.0

import "magalu.cloud/lib/products/network/subnetpools"
*/
package subnetpools

import (
	"context"

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

// A Schema representing a list of Subnet Pools.
type ListResult struct {
	Meta    ListResultMeta     `json:"meta"`
	Results *ListResultResults `json:"results,omitempty"`
}

type ListResultMeta struct {
	Links ListResultMetaLinks `json:"links"`
	Page  ListResultMetaPage  `json:"page"`
}

type ListResultMetaLinks struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Self     string `json:"self"`
}

type ListResultMetaPage struct {
	Count  int  `json:"count"`
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
	Total  int  `json:"total"`
}

// A Schema representing a response for Subnet Pool creation.
type ListResultResultsItem struct {
	Cidr        *string `json:"cidr,omitempty"`
	Description *string `json:"description,omitempty"`
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	TenantId    string  `json:"tenant_id"`
}

type ListResultResults []ListResultResultsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/subnetpools/list"), s.client, s.ctx)
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/subnetpools/list"), s.client, ctx)
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
