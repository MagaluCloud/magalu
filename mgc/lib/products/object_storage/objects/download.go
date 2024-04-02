/*
Executor: download

# Summary

# Download an object from a bucket

# Description

Download an object from a bucket. If no destination is specified, the default is the current working directory

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DownloadParameters struct {
	Dst        string `json:"dst,omitempty"`
	ObjVersion string `json:"obj_version,omitempty"`
	Src        string `json:"src"`
}

type DownloadConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type DownloadResult any

func Download(
	client *mgcClient.Client,
	ctx context.Context,
	parameters DownloadParameters,
	configs DownloadConfigs,
) (
	result DownloadResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Download", mgcCore.RefPath("/object-storage/objects/download"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DownloadParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DownloadConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DownloadResult](r)
}

// TODO: links
// TODO: related
