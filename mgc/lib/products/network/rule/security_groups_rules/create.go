/*
Executor: create

# Summary

# Create Rule

# Description

Create a Rule async, returning its ID. To monitor the creation progress, please check the status in the service message or implement polling.Either a remote_ip_prefix or a remote_group_id can be specified.With remote_ip_prefix, all IPs that match the criteria will be allowed.With remote_group_id, only the specified security group is allowed to communicatefollowing the specified protocol, direction and port_range_min/max

Version: 1.133.0

import "magalu.cloud/lib/products/network/rule/security_groups_rules"
*/
package securityGroupsRules

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Direction       *string `json:"direction,omitempty"`
	Ethertype       string  `json:"ethertype"`
	PortRangeMax    *int    `json:"port_range_max,omitempty"`
	PortRangeMin    *int    `json:"port_range_min,omitempty"`
	Protocol        *string `json:"protocol,omitempty"`
	RemoteGroupId   *string `json:"remote_group_id,omitempty"`
	RemoteIpPrefix  *string `json:"remote_ip_prefix,omitempty"`
	SecurityGroupId string  `json:"security_group_id"`
	ValidateQuota   *bool   `json:"validate_quota,omitempty"`
	Wait            *bool   `json:"wait,omitempty"`
	WaitTimeout     *int    `json:"wait_timeout,omitempty"`
}

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
	Id string `json:"id"`
}

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/network/rule/security-groups-rules/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) CreateContext(
	ctx context.Context,
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/network/rule/security-groups-rules/create"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
