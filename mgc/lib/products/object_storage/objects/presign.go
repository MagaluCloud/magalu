/*
Executor: presign

# Description

# Generate a pre-signed URL for accessing an object

import "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type PresignParameters struct {
	Dst       string  `json:"dst"`
	ExpiresIn *string `json:"expires-in,omitempty"`
	Method    string  `json:"method"`
}

type PresignConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type PresignResult struct {
	Url string `json:"url"`
}

func (s *service) Presign(
	parameters PresignParameters,
	configs PresignConfigs,
) (
	result PresignResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Presign", mgcCore.RefPath("/object-storage/objects/presign"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) PresignContext(
	ctx context.Context,
	parameters PresignParameters,
	configs PresignConfigs,
) (
	result PresignResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Presign", mgcCore.RefPath("/object-storage/objects/presign"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[PresignResult](r)
}

// TODO: links
// TODO: related
