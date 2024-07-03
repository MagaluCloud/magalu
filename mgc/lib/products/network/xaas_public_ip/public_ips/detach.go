/*
Executor: detach

# Summary

# Detach Public IP

# Description

# Detach a Public IP to a Port

Version: 1.126.1

import "magalu.cloud/lib/products/network/xaas_public_ip/public_ips"
*/
package publicIps

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DetachParameters struct {
	PortId      string `json:"port_id"`
	ProjectType string `json:"project_type"`
	PublicIpId  string `json:"public_ip_id"`
}

type DetachConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type DetachResult any

func (s *service) Detach(
	parameters DetachParameters,
	configs DetachConfigs,
) (
	result DetachResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Detach", mgcCore.RefPath("/network/xaas_public_ip/public-ips/detach"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DetachParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DetachConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DetachResult](r)
}

// TODO: links
// TODO: related
