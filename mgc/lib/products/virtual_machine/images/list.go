/*
Executor: list

# Summary

Retrieves all images.

# Description

Retrieve a list of images allowed for the current region.

Version: v1

import "github.com/MagaluCloud/magalu/mgc/lib/products/virtual_machine/images"
*/
package images

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type ListParameters struct {
	Labels           *string `json:"_labels,omitempty"`
	Limit            *int    `json:"_limit,omitempty"`
	Offset           *int    `json:"_offset,omitempty"`
	Sort             *string `json:"_sort,omitempty"`
	AvailabilityZone *string `json:"availability-zone,omitempty"`
}

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Images ListResultImages `json:"images"`
}

type ListResultImagesItem struct {
	AvailabilityZones    *ListResultImagesItemAvailabilityZones  `json:"availability_zones,omitempty"`
	EndLifeAt            *string                                 `json:"end_life_at,omitempty"`
	EndStandardSupportAt *string                                 `json:"end_standard_support_at,omitempty"`
	Id                   string                                  `json:"id"`
	Labels               *ListResultImagesItemLabels             `json:"labels,omitempty"`
	MinimumRequirements  ListResultImagesItemMinimumRequirements `json:"minimum_requirements"`
	Name                 string                                  `json:"name"`
	Platform             *string                                 `json:"platform,omitempty"`
	ReleaseAt            *string                                 `json:"release_at,omitempty"`
	Status               string                                  `json:"status"`
	Version              *string                                 `json:"version,omitempty"`
}

type ListResultImagesItemAvailabilityZones []string

type ListResultImagesItemLabels []string

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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ListContext(
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/images/list"), s.client, ctx)
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
