/*
Executor: delete-all

# Description

# Delete all objects from a bucket

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteAllParameters struct {
	BatchSize *int                       `json:"batch_size,omitempty"`
	Bucket    string                     `json:"bucket"`
	Filter    *DeleteAllParametersFilter `json:"filter,omitempty"`
}

type DeleteAllParametersFilterItem struct {
	Exclude *string `json:"exclude,omitempty"`
	Include *string `json:"include,omitempty"`
}

type DeleteAllParametersFilter []DeleteAllParametersFilterItem

type DeleteAllConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type DeleteAllResult any

func (s *service) DeleteAll(
	parameters DeleteAllParameters,
	configs DeleteAllConfigs,
) (
	result DeleteAllResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("DeleteAll", mgcCore.RefPath("/object-storage/objects/delete-all"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteAllParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeleteAllConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DeleteAllResult](r)
}

// TODO: links
// TODO: related
