/*
Executor: get

# Summary

# Rule Details

# Description

# Return a rule details

Version: 1.131.1

import "magalu.cloud/lib/products/network/rule/rules"
*/
package rules

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	RuleId string `json:"rule_id"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	CreatedAt       *string `json:"created_at,omitempty"`
	Direction       *string `json:"direction,omitempty"`
	Error           *string `json:"error,omitempty"`
	Ethertype       *string `json:"ethertype,omitempty"`
	ExternalId      *string `json:"external_id,omitempty"`
	Id              *string `json:"id,omitempty"`
	PortRangeMax    *int    `json:"port_range_max,omitempty"`
	PortRangeMin    *int    `json:"port_range_min,omitempty"`
	Protocol        *string `json:"protocol,omitempty"`
	RemoteGroupId   *string `json:"remote_group_id,omitempty"`
	RemoteIpPrefix  *string `json:"remote_ip_prefix,omitempty"`
	SecurityGroupId *string `json:"security_group_id,omitempty"`
	Status          string  `json:"status"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/network/rule/rules/get"), s.client, s.ctx)
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
