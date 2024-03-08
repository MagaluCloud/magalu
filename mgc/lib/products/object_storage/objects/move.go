/*
Executor: move

# Summary

# Moves one object from source to destination

# Description

Moves one object from a source to a destination.
It can be either local or remote but not both local (Local -> Remote, Remote -> Local, Remote -> Remote)

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type MoveParameters struct {
	Dst string `json:"dst"`
	Src string `json:"src"`
}

type MoveConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type MoveResult struct {
	Dst string `json:"dst"`
	Src string `json:"src"`
}

func Move(
	client *mgcClient.Client,
	ctx context.Context,
	parameters MoveParameters,
	configs MoveConfigs,
) (
	result MoveResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Move", mgcCore.RefPath("/object-storage/objects/move"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[MoveParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[MoveConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[MoveResult](r)
}

// TODO: links
// TODO: related
