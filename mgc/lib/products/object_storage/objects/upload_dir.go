/*
Executor: upload-dir

# Description

# Upload a directory to a bucket

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UploadDirParameters struct {
	Dst          string                     `json:"dst"`
	Filter       *UploadDirParametersFilter `json:"filter,omitempty"`
	Shallow      *bool                      `json:"shallow,omitempty"`
	Src          string                     `json:"src"`
	StorageClass *string                    `json:"storage_class,omitempty"`
}

type UploadDirParametersFilterItem struct {
	Exclude *string `json:"exclude,omitempty"`
	Include *string `json:"include,omitempty"`
}

type UploadDirParametersFilter []UploadDirParametersFilterItem

type UploadDirConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type UploadDirResult struct {
	Dir string `json:"dir"`
	Uri string `json:"uri"`
}

func (s *service) UploadDir(
	parameters UploadDirParameters,
	configs UploadDirConfigs,
) (
	result UploadDirResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("UploadDir", mgcCore.RefPath("/object-storage/objects/upload-dir"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UploadDirParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UploadDirConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[UploadDirResult](r)
}

// TODO: links
// TODO: related
