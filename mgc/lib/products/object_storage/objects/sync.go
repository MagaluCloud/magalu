/*
Executor: sync

# Summary

# Synchronizes a local path to a bucket

# Description

This command uploads any file from the source to the destination if it's not present or has a different size. Additionally any file in the destination not present on the source is deleted.

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SyncParameters struct {
	BatchSize int    `json:"batch_size,omitempty"`
	Delete    bool   `json:"delete,omitempty"`
	Dst       string `json:"dst"`
	Src       string `json:"src"`
}

type SyncConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type SyncResult any

func Sync(
	client *mgcClient.Client,
	ctx context.Context,
	parameters SyncParameters,
	configs SyncConfigs,
) (
	result SyncResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Sync", mgcCore.RefPath("/object-storage/objects/sync"), client, ctx)
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

// TODO: links
// TODO: related
