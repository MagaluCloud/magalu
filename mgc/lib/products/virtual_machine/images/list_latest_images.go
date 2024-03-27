/*
Executor: list-latest-images

# Summary

Retrieves all images available in the region.

# Description

Retrieve a list of images allowed for the current tenant which is logged in.

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/images"
*/
package images

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListLatestImagesParameters struct {
	Limit  int    `json:"_limit,omitempty"`
	Offset int    `json:"_offset,omitempty"`
	Sort   string `json:"_sort,omitempty"`
}

type ListLatestImagesConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListLatestImagesResult struct {
	Images ListLatestImagesResultImages `json:"images"`
}

type ListLatestImagesResultImagesItem struct {
	EndLifeAt            *string                                             `json:"end_life_at,omitempty"`
	EndStandardSupportAt *string                                             `json:"end_standard_support_at,omitempty"`
	Id                   string                                              `json:"id"`
	MinimumRequirements  ListLatestImagesResultImagesItemMinimumRequirements `json:"minimum_requirements"`
	Name                 string                                              `json:"name"`
	Platform             *string                                             `json:"platform,omitempty"`
	ReleaseAt            *string                                             `json:"release_at,omitempty"`
	Status               string                                              `json:"status"`
	Version              *string                                             `json:"version,omitempty"`
}

type ListLatestImagesResultImagesItemMinimumRequirements struct {
	Disk int `json:"disk"`
	Ram  int `json:"ram"`
	Vcpu int `json:"vcpu"`
}

type ListLatestImagesResultImages []ListLatestImagesResultImagesItem

func ListLatestImages(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListLatestImagesParameters,
	configs ListLatestImagesConfigs,
) (
	result ListLatestImagesResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("ListLatestImages", mgcCore.RefPath("/virtual-machine/images/list-latest-images"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListLatestImagesParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ListLatestImagesConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListLatestImagesResult](r)
}

// TODO: links
// TODO: related
