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
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UploadParameters struct {
	Dst          string  `json:"dst"`
	Src          string  `json:"src"`
	StorageClass *string `json:"storage_class,omitempty"`
}

type UploadConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type UploadResult struct {
	File string `json:"file"`
	Uri  string `json:"uri"`
}

func (s *service) Upload(
	parameters UploadParameters,
	configs UploadConfigs,
) (
	result UploadResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Upload", mgcCore.RefPath("/object-storage/objects/upload"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) UploadContext(
	ctx context.Context,
	parameters UploadParameters,
	configs UploadConfigs,
) (
	result UploadResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Upload", mgcCore.RefPath("/object-storage/objects/upload"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[UploadResult](r)
}

// TODO: links
// TODO: related
