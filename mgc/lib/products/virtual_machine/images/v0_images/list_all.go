/*
Executor: list-all

# Summary

# List Images All

# Description

# Retrive a list of all images

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

type ListAllConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListAllResult struct {
	Images ListAllResultImages `json:"images"`
}

type ListAllResultImagesItem struct {
	Id       string  `json:"id"`
	Internal bool    `json:"internal"`
	Name     string  `json:"name"`
	Sku      *string `json:"sku,omitempty"`
}

type ListAllResultImages []ListAllResultImagesItem

func ListAll(
	client *mgcClient.Client,
	ctx context.Context,
	configs ListAllConfigs,
) (
	result ListAllResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("ListAll", mgcCore.RefPath("/virtual-machine/images/v0-images/list-all"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ListAllConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListAllResult](r)
}

// TODO: links
// TODO: related
