/*
Executor: create

# Summary

Create a backup.

# Description

Create a backup for the currently authenticated tenant.

The Backup can be used when it reaches the "available" state and the

	"completed" status.

#### Rules
- The Backup name must be unique; otherwise, the creation will be disallowed.
- The Volume can be either in in-use or available states.
- The Volume must not have an operation in execution.

#### Notes
  - Use the **block-storage volume list** command to retrieve a list of all
    Volumes and obtain the ID of the Volume that will be used to create the
    Backup.

Version: v1

import "magalu.cloud/lib/products/block_storage/backups"
*/
package backups

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Description *string                `json:"description,omitempty"`
	Name        string                 `json:"name"`
	Type        *string                `json:"type,omitempty"`
	Volume      CreateParametersVolume `json:"volume"`
}

// any of: CreateParametersVolume
type CreateParametersVolume struct {
	Id   string  `json:"id"`
	Name *string `json:"name,omitempty"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/block-storage/backups/create"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) CreateContext(
	ctx context.Context,
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/block-storage/backups/create"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
