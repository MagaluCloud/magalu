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

import "github.com/MagaluCloud/magalu/mgc/lib/products/block_storage/snapshots"
*/
package snapshots

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type ListParameters struct {
	Limit  *int                  `json:"_limit,omitempty"`
	Offset *int                  `json:"_offset,omitempty"`
	Sort   *string               `json:"_sort,omitempty"`
	Expand *ListParametersExpand `json:"expand,omitempty"`
	Name   *string               `json:"name,omitempty"`
	Type   *string               `json:"type,omitempty"`
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
	AvailabilityZones ListResultSnapshotsItemAvailabilityZones `json:"availability_zones"`
	CreatedAt         string                                   `json:"created_at"`
	Description       *string                                  `json:"description"`
	Error             *ListResultSnapshotsItemError            `json:"error,omitempty"`
	Id                string                                   `json:"id"`
	Name              string                                   `json:"name"`
	Size              int                                      `json:"size"`
	State             string                                   `json:"state"`
	Status            string                                   `json:"status"`
	Type              string                                   `json:"type"`
	UpdatedAt         string                                   `json:"updated_at"`
	Volume            *ListResultSnapshotsItemVolume           `json:"volume"`
}

type ListResultSnapshotsItemAvailabilityZones []string

type ListResultSnapshotsItemError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// any of: *ListResultSnapshotsItemVolume
type ListResultSnapshotsItemVolume struct {
	Id   *string                            `json:"id,omitempty"`
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ListContext(
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/block-storage/snapshots/list"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
