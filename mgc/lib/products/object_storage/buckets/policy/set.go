/*
Executor: set

# Description

# Set policy document for the specified bucket

import "magalu.cloud/lib/products/object_storage/buckets/policy"
*/
package policy

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type SetParameters struct {
	Dst    string              `json:"dst"`
	Policy SetParametersPolicy `json:"policy"`
}

// Policy file path to be uploaded
type SetParametersPolicy struct {
}

type SetConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type SetResult any

func (s *service) Set(
	parameters SetParameters,
	configs SetConfigs,
) (
	result SetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Set", mgcCore.RefPath("/object-storage/buckets/policy/set"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[SetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[SetConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[SetResult](r)
}

// TODO: links
// TODO: related
