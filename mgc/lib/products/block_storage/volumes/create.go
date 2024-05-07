/*
Executor: create

# Summary

# Create a new Volume

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

import "magalu.cloud/lib/products/block_storage/volumes"
*/
package volumes

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Name string               `json:"name"`
	Size int                  `json:"size"`
	Type CreateParametersType `json:"type"`
}

// any of: CreateParametersType0, CreateParametersType1
type CreateParametersType struct {
	CreateParametersType0 `json:",squash"` // nolint
	CreateParametersType1 `json:",squash"` // nolint
}

type CreateParametersType0 struct {
	Id string `json:"id"`
}

type CreateParametersType1 struct {
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

// TODO: links
// TODO: related
