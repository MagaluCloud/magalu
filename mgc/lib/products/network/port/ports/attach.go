/*
Executor: attach

# Summary

# Attach Security Group

# Description

Attach a Security Group to a Port with provided port_id, security_group_id, x-tenant-id of an specific project type

Version: 1.131.1

import "magalu.cloud/lib/products/network/port/ports"
*/
package ports

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type AttachParameters struct {
	PortId          string `json:"port_id"`
	SecurityGroupId string `json:"security_group_id"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Attach", mgcCore.RefPath("/network/port/ports/attach"), s.client, s.ctx)
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
