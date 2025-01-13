/*
Executor: move

# Summary

# Moves one object from source to destination

# Description

Moves one object from a source to a destination.
It can be either local or remote but not both local (Local -> Remote, Remote -> Local, Remote -> Remote)

import "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type MoveParameters struct {
	Dst string `json:"dst"`
	Src string `json:"src"`
}

type MoveConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type MoveResult struct {
	Dst string `json:"dst"`
	Src string `json:"src"`
}

func (s *service) Move(
	parameters MoveParameters,
	configs MoveConfigs,
) (
	result MoveResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Move", mgcCore.RefPath("/object-storage/objects/move"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) MoveContext(
	ctx context.Context,
	parameters MoveParameters,
	configs MoveConfigs,
) (
	result MoveResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Move", mgcCore.RefPath("/object-storage/objects/move"), s.client, ctx)
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

	sdkConfig := s.client.Sdk().Config().TempConfig()
	if c["serverUrl"] == nil && sdkConfig["serverUrl"] != nil {
		c["serverUrl"] = sdkConfig["serverUrl"]
	}

	if c["env"] == nil && sdkConfig["env"] != nil {
		c["env"] = sdkConfig["env"]
	}

	if c["region"] == nil && sdkConfig["region"] != nil {
		c["region"] = sdkConfig["region"]
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[MoveResult](r)
}

// TODO: links
// TODO: related
