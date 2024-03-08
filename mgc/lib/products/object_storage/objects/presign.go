/*
Executor: presign

# Description

# Generate a pre-signed URL for accessing an object

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type PresignParameters struct {
	Dst       string `json:"dst"`
	ExpiresIn string `json:"expires-in"`
	Method    string `json:"method"`
}

type PresignConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type PresignResult struct {
	Url string `json:"url"`
}

func Presign(
	client *mgcClient.Client,
	ctx context.Context,
	parameters PresignParameters,
	configs PresignConfigs,
) (
	result PresignResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Presign", mgcCore.RefPath("/object-storage/objects/presign"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[PresignParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[PresignConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[PresignResult](r)
}

// TODO: links
// TODO: related
