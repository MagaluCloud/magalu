/*
Executor: enable

# Description

# Enable versioning for a Bucket

import "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/buckets/versioning"
*/
package versioning

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type EnableParameters struct {
	Bucket string `json:"bucket"`
}

type EnableConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type EnableResult any

func (s *service) Enable(
	parameters EnableParameters,
	configs EnableConfigs,
) (
	result EnableResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Enable", mgcCore.RefPath("/object-storage/buckets/versioning/enable"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[EnableParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[EnableConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[EnableResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) EnableContext(
	ctx context.Context,
	parameters EnableParameters,
	configs EnableConfigs,
) (
	result EnableResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Enable", mgcCore.RefPath("/object-storage/buckets/versioning/enable"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[EnableParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[EnableConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[EnableResult](r)
}

// TODO: links
// TODO: related
