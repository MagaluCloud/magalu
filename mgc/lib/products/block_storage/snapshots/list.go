/*
Executor: list

# Summary

List all snapshots.

# Description

Retrieve a list of Snapshots for the currently authenticated tenant.

#### Notes
  - Use the expand argument to obtain additional details about the Volume used to
    create each Snapshot.

Version: v1

import "magalu.cloud/lib/products/block_storage/snapshots"
*/
package snapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit  *int                  `json:"_limit,omitempty"`
	Offset *int                  `json:"_offset,omitempty"`
	Sort   *string               `json:"_sort,omitempty"`
	Expand *ListParametersExpand `json:"expand,omitempty"`
}

type ListParametersExpand []string

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Snapshots ListResultSnapshots `json:"snapshots"`
}

type ListResultSnapshotsItem struct {
	CreatedAt   string                        `json:"created_at"`
	Description *string                       `json:"description"`
	Error       *ListResultSnapshotsItemError `json:"error,omitempty"`
	Id          string                        `json:"id"`
	Name        string                        `json:"name"`
	Size        int                           `json:"size"`
	State       string                        `json:"state"`
	Status      string                        `json:"status"`
	UpdatedAt   string                        `json:"updated_at"`
	Volume      ListResultSnapshotsItemVolume `json:"volume"`
}

type ListResultSnapshotsItemError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// any of: ListResultSnapshotsItemVolume
type ListResultSnapshotsItemVolume struct {
	Id   string                             `json:"id"`
	Name *string                            `json:"name,omitempty"`
	Size *int                               `json:"size,omitempty"`
	Type *ListResultSnapshotsItemVolumeType `json:"type,omitempty"`
}

type ListResultSnapshotsItemVolumeType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ListResultSnapshots []ListResultSnapshotsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/block-storage/snapshots/list"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ListConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
