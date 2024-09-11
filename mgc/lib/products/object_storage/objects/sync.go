/*
Executor: sync

# Summary

# Synchronizes a local path with a bucket

# Description

This command uploads any file from the local path to the bucket if it is not already present or has modified time changed.

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SyncParameters struct {
	BatchSize *int   `json:"batch_size,omitempty"`
	Bucket    string `json:"bucket"`
	Delete    *bool  `json:"delete,omitempty"`
	Local     string `json:"local"`
}

type SyncConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type SyncResult any

func (s *service) Sync(
	parameters SyncParameters,
	configs SyncConfigs,
) (
	result SyncResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Sync", mgcCore.RefPath("/object-storage/objects/sync"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[SyncParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[SyncConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[SyncResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) SyncContext(
	ctx context.Context,
	parameters SyncParameters,
	configs SyncConfigs,
) (
	result SyncResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Sync", mgcCore.RefPath("/object-storage/objects/sync"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[SyncParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[SyncConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[SyncResult](r)
}

// TODO: links
// TODO: related
