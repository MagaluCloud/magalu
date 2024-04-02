/*
Executor: list-v1-images

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

type ListV1ImagesParameters struct {
	Limit  int    `json:"_limit,omitempty"`
	Offset int    `json:"_offset,omitempty"`
	Sort   string `json:"_sort,omitempty"`
}

type ListV1ImagesConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListV1ImagesResult struct {
	Images ListV1ImagesResultImages `json:"images"`
}

type ListV1ImagesResultImagesItem struct {
	EndLifeAt            *string                                         `json:"end_life_at,omitempty"`
	EndStandardSupportAt *string                                         `json:"end_standard_support_at,omitempty"`
	Id                   string                                          `json:"id"`
	MinimumRequirements  ListV1ImagesResultImagesItemMinimumRequirements `json:"minimum_requirements"`
	Name                 string                                          `json:"name"`
	Platform             *string                                         `json:"platform,omitempty"`
	ReleaseAt            *string                                         `json:"release_at,omitempty"`
	Status               string                                          `json:"status"`
	Version              *string                                         `json:"version,omitempty"`
}

type ListV1ImagesResultImagesItemMinimumRequirements struct {
	Disk int `json:"disk"`
	Ram  int `json:"ram"`
	Vcpu int `json:"vcpu"`
}

type ListV1ImagesResultImages []ListV1ImagesResultImagesItem

func ListV1Images(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListV1ImagesParameters,
	configs ListV1ImagesConfigs,
) (
	result ListV1ImagesResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("ListV1Images", mgcCore.RefPath("/virtual-machine/images/list-v1-images"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListV1ImagesParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ListV1ImagesConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListV1ImagesResult](r)
}

// TODO: links
// TODO: related
