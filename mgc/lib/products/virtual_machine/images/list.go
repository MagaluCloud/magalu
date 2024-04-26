/*
Executor: list

# Summary

Retrieves all images available in the region.

# Description

Retrieve a list of images allowed for the current tenant which is logged in.

Version: v1

import "magalu.cloud/lib/products/virtual_machine/images"
*/
package images

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit  int    `json:"_limit,omitempty"`
	Offset int    `json:"_offset,omitempty"`
	Sort   string `json:"_sort,omitempty"`
}

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Images ListResultImages `json:"images"`
}

type ListResultImagesItem struct {
	EndLifeAt            *string                                 `json:"end_life_at,omitempty"`
	EndStandardSupportAt *string                                 `json:"end_standard_support_at,omitempty"`
	Id                   string                                  `json:"id"`
	MinimumRequirements  ListResultImagesItemMinimumRequirements `json:"minimum_requirements"`
	Name                 string                                  `json:"name"`
	Platform             *string                                 `json:"platform,omitempty"`
	ReleaseAt            *string                                 `json:"release_at,omitempty"`
	Status               string                                  `json:"status"`
	Version              *string                                 `json:"version,omitempty"`
}

type ListResultImagesItemMinimumRequirements struct {
	Disk int `json:"disk"`
	Ram  int `json:"ram"`
	Vcpu int `json:"vcpu"`
}

type ListResultImages []ListResultImagesItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/images/list"), s.client, s.ctx)
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
