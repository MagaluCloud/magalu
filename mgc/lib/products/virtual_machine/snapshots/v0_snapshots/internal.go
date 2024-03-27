/*
Executor: internal

# Summary

# Update Internal Snapshot

# Description

update internal snapshot

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

type InternalParameters struct {
	Error         *string `json:"error,omitempty"`
	Id            string  `json:"id"`
	Size          *int    `json:"size,omitempty"`
	Status        string  `json:"status"`
	UrpSnapshotId *string `json:"urp_snapshot_id"`
}

type InternalConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

func Internal(
	client *mgcClient.Client,
	ctx context.Context,
	parameters InternalParameters,
	configs InternalConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Internal", mgcCore.RefPath("/virtual-machine/snapshots/v0-snapshots/internal"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[InternalParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[InternalConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
