/*
Executor: list

# Summary

# List Subnet Pools by Tenant

# Description

# Returns a list of Subnet Pools for the current tenant's project

Version: 1.133.0

import "magalu.cloud/lib/products/network/subnetpools"
*/
package subnetpools

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// A Pydantic model representing a list of Subnet Pools.
type ListResult struct {
	Subnetpools *ListResultSubnetpools `json:"subnetpools,omitempty"`
}

// A Pydantic model representing a response for Subnet Pool creation.
type ListResultSubnetpoolsItem struct {
	Cidr        *string `json:"cidr,omitempty"`
	Description *string `json:"description,omitempty"`
	ExternalId  *string `json:"external_id,omitempty"`
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	TenantId    string  `json:"tenant_id"`
}

type ListResultSubnetpools []ListResultSubnetpoolsItem

func (s *service) List(
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
