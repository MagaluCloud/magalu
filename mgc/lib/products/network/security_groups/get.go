/*
Executor: get

# Summary

# Security Group Details

# Description

# Return a security group details

Version: 1.125.3

import "magalu.cloud/lib/products/network/security_groups"
*/
package securityGroups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	SecurityGroupId string `json:"security_group_id"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	CreatedAt   *string         `json:"created_at,omitempty"`
	Description *string         `json:"description,omitempty"`
	Error       *string         `json:"error,omitempty"`
	ExternalId  *string         `json:"external_id,omitempty"`
	Id          *string         `json:"id,omitempty"`
	IsDefault   *bool           `json:"is_default,omitempty"`
	Name        *string         `json:"name,omitempty"`
	ProjectType *string         `json:"project_type,omitempty"`
	Rules       *GetResultRules `json:"rules,omitempty"`
	Status      string          `json:"status"`
	TenantId    *string         `json:"tenant_id,omitempty"`
	Updated     *string         `json:"updated,omitempty"`
	VpcId       *string         `json:"vpc_id,omitempty"`
}

type GetResultRulesItem struct {
	CreatedAt       *string `json:"created_at,omitempty"`
	Direction       *string `json:"direction,omitempty"`
	Error           *string `json:"error,omitempty"`
	Ethertype       *string `json:"ethertype,omitempty"`
	Id              *string `json:"id,omitempty"`
	PortRangeMax    *int    `json:"port_range_max,omitempty"`
	PortRangeMin    *int    `json:"port_range_min,omitempty"`
	Protocol        *string `json:"protocol,omitempty"`
	RemoteGroupId   *string `json:"remote_group_id,omitempty"`
	RemoteIpPrefix  *string `json:"remote_ip_prefix,omitempty"`
	SecurityGroupId *string `json:"security_group_id,omitempty"`
	Status          *string `json:"status,omitempty"`
}

type GetResultRules []GetResultRulesItem

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/network/security_groups/get"), s.client, s.ctx)
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
