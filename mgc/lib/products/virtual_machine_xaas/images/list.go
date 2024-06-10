/*
Executor: list

# Summary

# List Images All V1

# Description

Retrieve a list of images allowed for the current tenant which is logged in.

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/images"
*/
package images

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Images ListResultImages `json:"images"`
}

type ListResultImagesItem struct {
	EndLifeAt            *string `json:"end_life_at,omitempty"`
	EndStandardSupportAt *string `json:"end_standard_support_at,omitempty"`
	Id                   string  `json:"id"`
	Internal             bool    `json:"internal"`
	MinDisk              int     `json:"min_disk"`
	MinRam               int     `json:"min_ram"`
	MinVcpu              int     `json:"min_vcpu"`
	Name                 string  `json:"name"`
	Platform             *string `json:"platform,omitempty"`
	ReleaseAt            *string `json:"release_at,omitempty"`
	Sku                  string  `json:"sku"`
	Status               string  `json:"status"`
	Version              *string `json:"version,omitempty"`
}

type ListResultImages []ListResultImagesItem

func (s *service) List(
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine-xaas/images/list"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

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
