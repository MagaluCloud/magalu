/*
Executor: create

# Summary

# Start conciliator

# Description

# Start conciliator validation

Version: 1.131.1

import "magalu.cloud/lib/products/network/backoffice_conciliator"
*/
package backofficeConciliator

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateConfigs struct {
	Env          *string `json:"env,omitempty"`
	Limit        *int    `json:"limit,omitempty"`
	Region       *string `json:"region,omitempty"`
	ResourceType *string `json:"resource-type,omitempty"`
	ServerUrl    *string `json:"serverUrl,omitempty"`
	Skip         *int    `json:"skip,omitempty"`
}

type CreateResult any

func (s *service) Create(
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/network/backoffice_conciliator/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

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
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/network/backoffice_conciliator/create"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

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
