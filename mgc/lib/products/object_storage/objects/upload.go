/*
Executor: upload

# Description

# Upload a file to a bucket

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UploadParameters struct {
	Dst string `json:"dst"`
	Src string `json:"src"`
}

type UploadConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type UploadResult struct {
	File string `json:"file"`
	Uri  string `json:"uri"`
}

func Upload(
	client *mgcClient.Client,
	ctx context.Context,
	parameters UploadParameters,
	configs UploadConfigs,
) (
	result UploadResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Upload", mgcCore.RefPath("/object-storage/objects/upload"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UploadParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UploadConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[UploadResult](r)
}

// TODO: links
// TODO: related
