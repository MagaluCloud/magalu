/*
Executor: list

# Summary

# List Images

# Description

# Retrive a list of images

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/images/v0_images"
*/
package v0Images

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Images ListResultImages `json:"images,omitempty"`
}

type ListResultImagesItem struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	OsDistribution *string `json:"os_distribution,omitempty"`
	Size           int     `json:"size"`
	Sku            *string `json:"sku,omitempty"`
	Version        *string `json:"version,omitempty"`
}

type ListResultImages []ListResultImagesItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/images/v0-images/list"), client, ctx)
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
