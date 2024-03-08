/*
Executor: list

# Summary

Retrieves all images available in the region.

# Description

Retrieve a list of images allowed for the current tenant which is logged in.

Version: 0.1.0

import "magalu.cloud/lib/products/virtual_machine/images"
*/
package images

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
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
	EndLifeAt            *string `json:"end_life_at,omitempty"`
	EndStandardSupportAt *string `json:"end_standard_support_at,omitempty"`
	Id                   string  `json:"id"`
	MinDisk              int     `json:"min_disk"`
	MinRam               int     `json:"min_ram"`
	MinVcpu              int     `json:"min_vcpu"`
	Name                 string  `json:"name"`
	Platform             *string `json:"platform,omitempty"`
	ReleaseAt            *string `json:"release_at,omitempty"`
	Status               string  `json:"status"`
	Version              *string `json:"version,omitempty"`
}

type ListResultImages []ListResultImagesItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/images/list"), client, ctx)
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
