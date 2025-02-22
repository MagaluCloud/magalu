/*
Executor: list

# Summary

Lists all snapshots.

# Description

List all snapshots in the current tenant which is logged in.

#### Notes
- You can use the **expand** argument to get more details from the inner objects
like image and machine types.

Version: v1

import "github.com/MagaluCloud/magalu/mgc/lib/products/virtual_machine/snapshots"
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
	CreatedAt string                           `json:"created_at"`
	Id        string                           `json:"id"`
	Instance  *ListResultSnapshotsItemInstance `json:"instance"`
	Name      *string                          `json:"name,omitempty"`
	Size      int                              `json:"size"`
	State     string                           `json:"state"`
	Status    string                           `json:"status"`
	UpdatedAt *string                          `json:"updated_at,omitempty"`
}

type ListResultSnapshotsItemInstance struct {
	Id          string                                      `json:"id"`
	Image       *ListResultSnapshotsItemInstanceImage       `json:"image"`
	MachineType *ListResultSnapshotsItemInstanceMachineType `json:"machine_type"`
}

// any of: ListResultSnapshotsItemInstanceImage
type ListResultSnapshotsItemInstanceImage struct {
	Id       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	Platform *string `json:"platform,omitempty"`
}

// any of: ListResultSnapshotsItemInstanceMachineType
type ListResultSnapshotsItemInstanceMachineType struct {
	Disk  *int    `json:"disk,omitempty"`
	Id    string  `json:"id"`
	Name  *string `json:"name,omitempty"`
	Ram   *int    `json:"ram,omitempty"`
	Vcpus *int    `json:"vcpus,omitempty"`
}

type ListResultSnapshots []ListResultSnapshotsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/snapshots/list"), s.client, s.ctx)
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/snapshots/list"), s.client, ctx)
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
