/*
Executor: detach

# Summary

# Detach a non primary port from instance

# Description

# Detach a non primary port from instance for a default project

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances/port"
*/
package port

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DetachParameters struct {
	Id     string `json:"id"`
	PortId string `json:"port_id"`
}

type DetachConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Detach(
	parameters DetachParameters,
	configs DetachConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Detach", mgcCore.RefPath("/virtual-machine/instances/port/detach"), s.client, s.ctx)
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

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
