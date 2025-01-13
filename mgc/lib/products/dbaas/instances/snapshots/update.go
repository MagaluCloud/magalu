/*
Executor: update

# Summary

Snapshot Update.

# Description

Updates a snapshot.

Version: 1.34.1

import "github.com/MagaluCloud/magalu/mgc/lib/products/dbaas/instances/snapshots"
*/
package snapshots

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type UpdateParameters struct {
	Description *string `json:"description,omitempty"`
	InstanceId  string  `json:"instance_id"`
	Name        *string `json:"name,omitempty"`
	SnapshotId  string  `json:"snapshot_id"`
}

type UpdateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type UpdateResult struct {
	AllocatedSize int                  `json:"allocated_size"`
	CreatedAt     string               `json:"created_at"`
	Description   string               `json:"description"`
	FinishedAt    *string              `json:"finished_at,omitempty"`
	Id            string               `json:"id"`
	Instance      UpdateResultInstance `json:"instance"`
	Name          string               `json:"name"`
	StartedAt     *string              `json:"started_at,omitempty"`
	Status        string               `json:"status"`
	Type          string               `json:"type"`
	UpdatedAt     *string              `json:"updated_at,omitempty"`
}

// This response object provides details about a database instance associated with a snapshot.

type UpdateResultInstance struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (s *service) Update(
	parameters UpdateParameters,
	configs UpdateConfigs,
) (
	result UpdateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/dbaas/instances/snapshots/update"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UpdateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UpdateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[UpdateResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) UpdateContext(
	ctx context.Context,
	parameters UpdateParameters,
	configs UpdateConfigs,
) (
	result UpdateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/dbaas/instances/snapshots/update"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UpdateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UpdateConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[UpdateResult](r)
}

// TODO: links
// TODO: related
