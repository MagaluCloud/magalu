/*
Executor: update

# Summary

# Patches a cluster

# Description

# Patches the mutable fields of a cluster

Version: 0.1.0

import "magalu.cloud/lib/products/kubernetes/cluster"
*/
package cluster

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UpdateParameters struct {
	AllowedCidrs *UpdateParametersAllowedCidrs `json:"allowed_cidrs,omitempty"`
	ClusterId    string                        `json:"cluster_id"`
}

type UpdateParametersAllowedCidrs []string

type UpdateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Object for cluster patch response
type UpdateResult struct {
	AllowedCidrs *UpdateResultAllowedCidrs `json:"allowed_cidrs,omitempty"`
}

type UpdateResultAllowedCidrs []string

func (s *service) Update(
	parameters UpdateParameters,
	configs UpdateConfigs,
) (
	result UpdateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/kubernetes/cluster/update"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UpdateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UpdateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[UpdateResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) UpdateContext(
	ctx context.Context,
	parameters UpdateParameters,
	configs UpdateConfigs,
) (
	result UpdateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/kubernetes/cluster/update"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UpdateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UpdateConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[UpdateResult](r)
}

// TODO: links
// TODO: related
