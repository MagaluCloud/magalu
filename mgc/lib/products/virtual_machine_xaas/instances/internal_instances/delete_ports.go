/*
Executor: delete-ports

# Summary

# Delete Port

# Description

# Delete a not primary port

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/instances/internal_instances"
*/
package internalInstances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeletePortsParameters struct {
	Id          string  `json:"id"`
	PortId      string  `json:"port_id"`
	ProjectType *string `json:"project_type,omitempty"`
}

type DeletePortsConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) DeletePorts(
	parameters DeletePortsParameters,
	configs DeletePortsConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("DeletePorts", mgcCore.RefPath("/virtual-machine-xaas/instances/internal-instances/delete-ports"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeletePortsParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeletePortsConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

func (s *service) DeletePortsConfirmPrompt(
	parameters DeletePortsParameters,
	configs DeletePortsConfigs,
) (message string) {
	e, err := mgcHelpers.ResolveExecutor("DeletePorts", mgcCore.RefPath("/virtual-machine-xaas/instances/internal-instances/delete-ports"), s.client)
	if err != nil {
		return
	}

	exec, ok := e.(mgcCore.ConfirmableExecutor)
	if !ok {
		// Not expected, but let's return an empty message
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeletePortsParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeletePortsConfigs](configs); err != nil {
		return
	}

	return exec.ConfirmPrompt(p, c)
}

// TODO: links
// TODO: related
