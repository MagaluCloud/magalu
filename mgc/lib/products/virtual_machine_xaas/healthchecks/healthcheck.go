/*
Executor: healthcheck

# Summary

# Healthcheck

# Description

# Check api status

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/healthchecks"
*/
package healthchecks

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type HealthcheckConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type HealthcheckResult struct {
	Status string `json:"status"`
}

func (s *service) Healthcheck(
	configs HealthcheckConfigs,
) (
	result HealthcheckResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Healthcheck", mgcCore.RefPath("/virtual-machine-xaas/healthchecks/healthcheck"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[HealthcheckConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[HealthcheckResult](r)
}

// TODO: links
// TODO: related
