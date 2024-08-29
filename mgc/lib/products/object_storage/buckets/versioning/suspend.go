/*
Executor: suspend

# Description

# Suspend versioning for a Bucket

import "magalu.cloud/lib/products/object_storage/buckets/versioning"
*/
package versioning

import (
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
	if c, err = mgcHelpers.ConvertConfigs[map[string]interface{}](s.client.Sdk().Config().TempConfig()); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[SuspendResult](r)
}

// TODO: links
// TODO: related
