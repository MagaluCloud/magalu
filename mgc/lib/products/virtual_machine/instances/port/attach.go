/*
Executor: attach

# Summary

# Attach port to instance

# Description

# Attach port to instance for a default project

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances/port"
*/
package port

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type AttachParameters struct {
	Id     string `json:"id"`
	PortId string `json:"port_id"`
}

type AttachConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Attach(
	parameters AttachParameters,
	configs AttachConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Attach", mgcCore.RefPath("/virtual-machine/instances/port/attach"), s.client, s.ctx)
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

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
