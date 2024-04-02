/*
Executor: delete

# Summary

# Delete Snapshot

# Description

# Delete a snapshot

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/snapshots/v0_snapshots"
*/
package v0Snapshots

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteParameters struct {
	Id string `json:"id"`
}

type DeleteConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

func Delete(
	client *mgcClient.Client,
	ctx context.Context,
	parameters DeleteParameters,
	configs DeleteConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Delete", mgcCore.RefPath("/virtual-machine/snapshots/v0-snapshots/delete"), client, ctx)
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

func DeleteConfirmPrompt(
	client *mgcClient.Client,
	parameters DeleteParameters,
	configs DeleteConfigs,
) (message string) {
	e, err := mgcHelpers.ResolveExecutor("Delete", mgcCore.RefPath("/virtual-machine/snapshots/v0-snapshots/delete"), client)
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
