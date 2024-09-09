/*
Executor: delete

# Description

# Delete policy document for the specified bucket

import "magalu.cloud/lib/products/object_storage/buckets/policy"
*/
package policy

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteParameters struct {
	Dst string `json:"dst"`
}

type DeleteConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type DeleteResult any

/*func (s *service) Delete(
	parameters DeleteParameters,
	configs DeleteConfigs,
) (
	result DeleteResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Delete", mgcCore.RefPath("/object-storage/buckets/policy/delete"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeleteConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DeleteResult](r)
}*/

// Context from caller is used to allow cancellation of long-running requests
func (s *service) DeleteContext(
	ctx context.Context,
	parameters DeleteParameters,
	configs DeleteConfigs,
) (
	result DeleteResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Delete", mgcCore.RefPath("/object-storage/buckets/policy/delete"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeleteConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DeleteResult](r)
}

// TODO: links
// TODO: related
