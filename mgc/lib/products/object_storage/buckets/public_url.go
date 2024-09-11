/*
Executor: public-url

# Description

# Get bucket public url

import "magalu.cloud/lib/products/object_storage/buckets"
*/
package buckets

import (
	"context"

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
	exec, ctx, err := mgcHelpers.PrepareExecutor("PublicUrl", mgcCore.RefPath("/object-storage/buckets/public-url"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) PublicUrlContext(
	ctx context.Context,
	parameters PublicUrlParameters,
	configs PublicUrlConfigs,
) (
	result PublicUrlResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("PublicUrl", mgcCore.RefPath("/object-storage/buckets/public-url"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[PublicUrlResult](r)
}

// TODO: links
// TODO: related
