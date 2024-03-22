/*
Executor: download-all

# Description

# Download all objects from a bucket

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DownloadAllParameters struct {
	Dst    string                      `json:"dst,omitempty"`
	Filter DownloadAllParametersFilter `json:"filter,omitempty"`
	Src    string                      `json:"src"`
}

type DownloadAllParametersFilterItem struct {
	Exclude string `json:"exclude,omitempty"`
	Include string `json:"include,omitempty"`
}

type DownloadAllParametersFilter []DownloadAllParametersFilterItem

type DownloadAllConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type DownloadAllResult struct {
	Dst        string `json:"dst,omitempty"`
	ObjVersion string `json:"obj_version,omitempty"`
	Src        string `json:"src"`
}

func DownloadAll(
	client *mgcClient.Client,
	ctx context.Context,
	parameters DownloadAllParameters,
	configs DownloadAllConfigs,
) (
	result DownloadAllResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("DownloadAll", mgcCore.RefPath("/object-storage/objects/download-all"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DownloadAllParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DownloadAllConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DownloadAllResult](r)
}

// TODO: links
// TODO: related