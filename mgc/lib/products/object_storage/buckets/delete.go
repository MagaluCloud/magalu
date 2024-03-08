/*
Executor: delete

# Description

# Delete an existing Bucket

import "magalu.cloud/lib/products/object_storage/buckets"
*/
package buckets

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteParameters struct {
	Bucket string `json:"bucket"`
	Force  bool   `json:"force"`
}

type DeleteConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type DeleteResult any

func Delete(
	client *mgcClient.Client,
	ctx context.Context,
	parameters DeleteParameters,
	configs DeleteConfigs,
) (
	result DeleteResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Delete", mgcCore.RefPath("/object-storage/buckets/delete"), client, ctx)
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

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DeleteResult](r)
}

// TODO: links
// TODO: related
