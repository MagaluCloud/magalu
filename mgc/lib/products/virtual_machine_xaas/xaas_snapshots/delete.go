/*
Executor: delete

# Summary

# Delete a Snapshot asynchronously

# Description

Deletes an snapshot with the id provided in the current tenant
which is logged in.

#### Notes
- You can use the Snapshots list command to retrieve all snapshots, so
- you can get the id of the snapshot that you want to delete.

Version: 1.229.0

import "magalu.cloud/lib/products/virtual_machine_xaas/xaas_snapshots"
*/
package xaasSnapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteParameters struct {
	Id          string `json:"id"`
	ProjectType string `json:"project_type"`
}

type DeleteConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Delete(
	parameters DeleteParameters,
	configs DeleteConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Delete", mgcCore.RefPath("/virtual-machine-xaas/xaas snapshots/delete"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeleteConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

func (s *service) DeleteConfirmPrompt(
	parameters DeleteParameters,
	configs DeleteConfigs,
) (message string) {
	e, err := mgcHelpers.ResolveExecutor("Delete", mgcCore.RefPath("/virtual-machine-xaas/xaas snapshots/delete"), s.client)
	if err != nil {
		return
	}

	exec, ok := e.(mgcCore.ConfirmableExecutor)
	if !ok {
		// Not expected, but let's return an empty message
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeleteConfigs](configs); err != nil {
		return
	}

	return exec.ConfirmPrompt(p, c)
}

// TODO: links
// TODO: related
