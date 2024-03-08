/*
Executor: copy

# Description

# Copy an object from a bucket to another bucket

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CopyParameters struct {
	Dst        string `json:"dst"`
	ObjVersion string `json:"obj_version,omitempty"`
	Src        string `json:"src"`
}

type CopyConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type CopyResult any

func Copy(
	client *mgcClient.Client,
	ctx context.Context,
	parameters CopyParameters,
	configs CopyConfigs,
) (
	result CopyResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Copy", mgcCore.RefPath("/object-storage/objects/copy"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CopyParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CopyConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CopyResult](r)
}

// TODO: links
// TODO: related
