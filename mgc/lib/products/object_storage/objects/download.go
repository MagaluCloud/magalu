/*
Executor: download

# Summary

# Download an object from a bucket

# Description

Download an object from a bucket. If no destination is specified, the default is the current working directory

import "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type DownloadParameters struct {
	Dst        *string `json:"dst,omitempty"`
	ObjVersion *string `json:"obj_version,omitempty"`
	Src        string  `json:"src"`
}

type DownloadConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type DownloadResult any

func (s *service) Download(
	parameters DownloadParameters,
	configs DownloadConfigs,
) (
	result DownloadResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Download", mgcCore.RefPath("/object-storage/objects/download"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) DownloadContext(
	ctx context.Context,
	parameters DownloadParameters,
	configs DownloadConfigs,
) (
	result DownloadResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Download", mgcCore.RefPath("/object-storage/objects/download"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[DownloadResult](r)
}

// TODO: links
// TODO: related
