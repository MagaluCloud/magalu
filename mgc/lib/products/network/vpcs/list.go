/*
Executor: list

# Summary

# List VPC

# Description

Returns a list of VPCs for a provided tenant_id

Version: 1.126.1

import "magalu.cloud/lib/products/network/vpcs"
*/
package vpcs

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Vpcs *ListResultVpcs `json:"vpcs,omitempty"`
}

type ListResultVpcsItem struct {
	CreatedAt       *string                           `json:"created_at,omitempty"`
	Description     *string                           `json:"description,omitempty"`
	ExternalNetwork *string                           `json:"external_network,omitempty"`
	Id              *string                           `json:"id,omitempty"`
	IsDefault       *bool                             `json:"is_default,omitempty"`
	Name            *string                           `json:"name,omitempty"`
	NetworkId       *string                           `json:"network_id,omitempty"`
	RouterId        *string                           `json:"router_id,omitempty"`
	SecurityGroups  *ListResultVpcsItemSecurityGroups `json:"security_groups,omitempty"`
	Subnets         *ListResultVpcsItemSubnets        `json:"subnets,omitempty"`
	TenantId        *string                           `json:"tenant_id,omitempty"`
	Updated         *string                           `json:"updated,omitempty"`
}

type ListResultVpcsItemSecurityGroups []string

type ListResultVpcsItemSubnets []string

type ListResultVpcs []ListResultVpcsItem

func (s *service) List(
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/vpcs/list"), s.client, s.ctx)
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

// TODO: links
// TODO: related
