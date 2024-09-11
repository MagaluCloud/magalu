/*
Executor: suspend

# Description

# Suspend versioning for a Bucket

import "magalu.cloud/lib/products/object_storage/buckets/versioning"
*/
package versioning

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SuspendParameters struct {
	Bucket string `json:"bucket"`
}

type SuspendConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type SuspendResult any

func (s *service) Suspend(
	parameters SuspendParameters,
	configs SuspendConfigs,
) (
	result SuspendResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Suspend", mgcCore.RefPath("/object-storage/buckets/versioning/suspend"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[SuspendParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[SuspendConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[SuspendResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) SuspendContext(
	ctx context.Context,
	parameters SuspendParameters,
	configs SuspendConfigs,
) (
	result SuspendResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Suspend", mgcCore.RefPath("/object-storage/buckets/versioning/suspend"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[SuspendParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[SuspendConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[SuspendResult](r)
}

// TODO: links
// TODO: related
