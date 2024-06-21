/*
Executor: create

# Summary

# Create Security Group

# Description

# Create a Security Group

Version: 1.125.3

import "magalu.cloud/lib/products/network/security_groups"
*/
package securityGroups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Description   *string `json:"description,omitempty"`
	Name          *string `json:"name"`
	ValidateQuota *bool   `json:"validate_quota,omitempty"`
	Wait          *bool   `json:"wait,omitempty"`
	WaitTimeout   *int    `json:"wait_timeout,omitempty"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/network/security_groups/create"), s.client, s.ctx)
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

// TODO: links
// TODO: related
