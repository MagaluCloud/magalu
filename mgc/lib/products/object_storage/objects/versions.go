/*
Executor: versions

# Description

# Retrieve all versions of an object

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type VersionsParameters struct {
	Dst string `json:"dst"`
}

type VersionsConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type VersionsResultItem struct {
	ETag         string                  `json:"ETag"`
	IsLatest     bool                    `json:"IsLatest"`
	Key          string                  `json:"Key"`
	LastModified string                  `json:"LastModified"`
	Owner        VersionsResultItemOwner `json:"Owner"`
	Size         int                     `json:"Size"`
	StorageClass string                  `json:"StorageClass"`
	Text         string                  `json:"Text"`
	VersionId    string                  `json:"VersionID"`
}

type VersionsResultItemOwner struct {
	DisplayName string `json:"DisplayName"`
	Id          string `json:"ID"`
}

type VersionsResult []VersionsResultItem

func (s *service) Versions(
	parameters VersionsParameters,
	configs VersionsConfigs,
) (
	result VersionsResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Versions", mgcCore.RefPath("/object-storage/objects/versions"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[VersionsParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[VersionsConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[VersionsResult](r)
}

// TODO: links
// TODO: related
