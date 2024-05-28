/*
Executor: list

# Summary

# List Keypairs Xaas

# Description

Returns a list of keypairs from a provided tenant_id

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/xaas_keypairs"
*/
package xaasKeypairs

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	ProjectType string `json:"project_type"`
}

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Keypairs ListResultKeypairs `json:"keypairs"`
}

type ListResultKeypairsItem struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

type ListResultKeypairs []ListResultKeypairsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine-xaas/xaas keypairs/list"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListParameters](parameters); err != nil {
		return
	}

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
