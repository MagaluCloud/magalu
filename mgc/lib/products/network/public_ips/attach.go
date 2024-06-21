/*
Executor: attach

# Summary

# Attach Public IP

# Description

# Attach a Public IP to a Port

Version: 1.125.3

import "magalu.cloud/lib/products/network/public_ips"
*/
package publicIps

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type AttachParameters struct {
	PortId     string `json:"port_id"`
	PublicIpId string `json:"public_ip_id"`
}

type AttachConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type AttachResult any

func (s *service) Attach(
	parameters AttachParameters,
	configs AttachConfigs,
) (
	result AttachResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Attach", mgcCore.RefPath("/network/public_ips/attach"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[AttachParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[AttachConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[AttachResult](r)
}

// TODO: links
// TODO: related
