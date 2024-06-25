/*
Executor: create

# Summary

# Create a snapshot

# Description

Create a Snapshot for the currently authenticated tenant.

The Snapshot can be used when it reaches the "available" state and the

	"completed" status.

#### Rules
  - The Snapshot name must be unique; otherwise, the creation will be disallowed.
  - Creating Snapshots from restored Volumes may lead to future conflicts as
    you can't delete a Volume with an Snapshot and can't delete a Snapshot with a
    restored Volume, so we recommend avoiding it.

#### Notes
  - Use the **block-storage volume list** command to retrieve a list of all
    Volumes and obtain the ID of the Volume that will be used to create the
    Snapshot.

Version: v1

import "magalu.cloud/lib/products/block_storage/snapshots"
*/
package snapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Description *string                `json:"description"`
	Name        string                 `json:"name"`
	Volume      CreateParametersVolume `json:"volume"`
}

// any of: , CreateParametersVolume1
type CreateParametersVolume struct {
	CreateParametersVolume1 `json:",squash"` // nolint
}

type CreateParametersVolume1 struct {
	Name string `json:"name"`
}

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
	Id string `json:"id"`
}

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/block-storage/snapshots/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
