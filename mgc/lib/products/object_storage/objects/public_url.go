/*
Executor: public-url

# Description

# Get object public url

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type PublicUrlParameters struct {
	Dst string `json:"dst"`
}

type PublicUrlConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type PublicUrlResult struct {
	Url string `json:"url"`
}

func (s *service) PublicUrl(
	parameters PublicUrlParameters,
	configs PublicUrlConfigs,
) (
	result PublicUrlResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("PublicUrl", mgcCore.RefPath("/object-storage/objects/public-url"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[PublicUrlParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[PublicUrlConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[PublicUrlResult](r)
}

// TODO: links
// TODO: related
