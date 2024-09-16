/*
Executor: conciliator

# Summary

# Start conciliator

# Description

# Start conciliator validation

Version: 1.133.0

import "magalu.cloud/lib/products/network/internal"
*/
package internal

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ConciliatorConfigs struct {
	Env          *string `json:"env,omitempty"`
	Limit        *int    `json:"limit,omitempty"`
	Region       *string `json:"region,omitempty"`
	ResourceType *string `json:"resource-type,omitempty"`
	ServerUrl    *string `json:"serverUrl,omitempty"`
	Skip         *int    `json:"skip,omitempty"`
}

type ConciliatorResult any

func (s *service) Conciliator(
	configs ConciliatorConfigs,
) (
	result ConciliatorResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Conciliator", mgcCore.RefPath("/network/internal/conciliator"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ConciliatorConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ConciliatorResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ConciliatorContext(
	ctx context.Context,
	configs ConciliatorConfigs,
) (
	result ConciliatorResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Conciliator", mgcCore.RefPath("/network/internal/conciliator"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ConciliatorConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[ConciliatorResult](r)
}

// TODO: links
// TODO: related
