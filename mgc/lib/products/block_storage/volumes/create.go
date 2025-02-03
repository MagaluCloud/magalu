/*
Executor: create

# Summary

Create a new volume.

# Description

Create a Volume for the currently authenticated tenant.

The Volume can be used when it reaches the "available" state and "completed"

	status.

#### Rules
- The Volume name must be unique; otherwise, the creation will be disallowed.
- The Volume type must be available to use.

#### Notes
  - Utilize the **block-storage volume-types list** command to retrieve a list
    of all available Volume Types.
  - Verify the state and status of your Volume using the

**block-storage volume get --id [uuid]** command".

Version: v1

import "github.com/MagaluCloud/magalu/mgc/lib/products/block_storage/volumes"
*/
package volumes

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type CreateParameters struct {
	AvailabilityZone *string                   `json:"availability_zone,omitempty"`
	Encrypted        *bool                     `json:"encrypted,omitempty"`
	Name             string                    `json:"name"`
	Size             int                       `json:"size"`
	Snapshot         *CreateParametersSnapshot `json:"snapshot,omitempty"`
	Type             CreateParametersType      `json:"type"`
}

// any of: *CreateParametersSnapshot
type CreateParametersSnapshot struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// any of: *CreateParametersType
type CreateParametersType struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
	Id   *string           `json:"id,omitempty"`
	Name *string           `json:"name,omitempty"`
	Size *int              `json:"size,omitempty"`
	Type *CreateResultType `json:"type,omitempty"`
}

type CreateResultType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/block-storage/volumes/create"), s.client, s.ctx)
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/block-storage/volumes/create"), s.client, ctx)
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
