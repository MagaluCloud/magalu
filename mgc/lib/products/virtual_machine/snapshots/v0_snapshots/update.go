/*
Executor: update

# Summary

# Update Snapshots With Urp Snapshot Id

# Description

update snapshot with urp snapshot id

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

type UpdateParameters struct {
	Error         *string `json:"error,omitempty"`
	Size          *int    `json:"size,omitempty"`
	Status        string  `json:"status"`
	UrpSnapshotId string  `json:"urp_snapshot_id"`
}

type UpdateConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

func Update(
	client *mgcClient.Client,
	ctx context.Context,
	parameters UpdateParameters,
	configs UpdateConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/virtual-machine/snapshots/v0-snapshots/update"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UpdateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UpdateConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
