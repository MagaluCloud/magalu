/*
Executor: delete-id

# Summary

# Instance Delete Worker

# Description

This route is when the vm-instace worker already
deleted the instance from: urp and vpc api to mark the instance

to 'deleted' on virtual machine DB.

### Note
This route is used only for internal proposes.

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/instances/internal_instances"
*/
package internalInstances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteIdParameters struct {
	Id          string  `json:"id"`
	ProjectType *string `json:"project_type,omitempty"`
}

type DeleteIdConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) DeleteId(
	parameters DeleteIdParameters,
	configs DeleteIdConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("DeleteId", mgcCore.RefPath("/virtual-machine-xaas/instances/internal-instances/delete-id"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteIdParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeleteIdConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

func (s *service) DeleteIdConfirmPrompt(
	parameters DeleteIdParameters,
	configs DeleteIdConfigs,
) (message string) {
	e, err := mgcHelpers.ResolveExecutor("DeleteId", mgcCore.RefPath("/virtual-machine-xaas/instances/internal-instances/delete-id"), s.client)
	if err != nil {
		return
	}

	exec, ok := e.(mgcCore.ConfirmableExecutor)
	if !ok {
		// Not expected, but let's return an empty message
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteIdParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeleteIdConfigs](configs); err != nil {
		return
	}

	return exec.ConfirmPrompt(p, c)
}

// TODO: links
// TODO: related
