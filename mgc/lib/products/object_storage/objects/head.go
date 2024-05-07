/*
Executor: head

# Description

# Get object metadata

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type HeadParameters struct {
	Dst        string  `json:"dst"`
	ObjVersion *string `json:"objVersion,omitempty"`
}

type HeadConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type HeadResult struct {
	AcceptRanges  string `json:"AcceptRanges"`
	ContentLength int    `json:"ContentLength"`
	ContentType   string `json:"ContentType"`
	ETag          string `json:"ETag"`
	LastModified  string `json:"LastModified"`
}

func (s *service) Head(
	parameters HeadParameters,
	configs HeadConfigs,
) (
	result HeadResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Head", mgcCore.RefPath("/object-storage/objects/head"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[HeadParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[HeadConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[HeadResult](r)
}

// TODO: links
// TODO: related
