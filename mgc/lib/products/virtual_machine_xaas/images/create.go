/*
Executor: create

# Summary

# Create Image V1

# Description

Register a URP Image on Virtual Machine DB.

### Note
The Image on URP need to be public and protected.

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/images"
*/
package images

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	EndLifeAt            *string `json:"end_life_at,omitempty"`
	EndStandardSupportAt *string `json:"end_standard_support_at,omitempty"`
	ImageId              string  `json:"image_id"`
	ImageUrl             *string `json:"image_url,omitempty"`
	Internal             bool    `json:"internal"`
	MinDisk              int     `json:"min_disk"`
	MinRam               int     `json:"min_ram"`
	MinVcpu              int     `json:"min_vcpu"`
	Name                 string  `json:"name"`
	Platform             *string `json:"platform,omitempty"`
	ReleaseAt            *string `json:"release_at,omitempty"`
	Sku                  string  `json:"sku"`
	Status               *string `json:"status,omitempty"`
	Version              *string `json:"version,omitempty"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine-xaas/images/create"), s.client, s.ctx)
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
