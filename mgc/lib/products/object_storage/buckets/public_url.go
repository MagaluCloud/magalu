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
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type PublicUrlParameters struct {
	Dst string `json:"dst"`
}

type PublicUrlConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type PublicUrlResult struct {
	Url string `json:"url"`
}

func PublicUrl(
	client *mgcClient.Client,
	ctx context.Context,
	parameters PublicUrlParameters,
	configs PublicUrlConfigs,
) (
	result PublicUrlResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("PublicUrl", mgcCore.RefPath("/object-storage/buckets/public-url"), client, ctx)
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
