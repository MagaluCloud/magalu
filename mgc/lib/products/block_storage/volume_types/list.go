/*
Executor: list

# Summary

# List volume types

# Description

List Volume Types allowed in the current region.

#### Notes

  - Volume types are managed internally. If you wish to use a Volume Type that
    is not yet available, please contact our support team for assistance.

Version: v1

import "magalu.cloud/lib/products/block_storage/volume_types"
*/
package volumeTypes

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
	Types ListResultTypes `json:"types"`
}

type ListResultTypesItem struct {
	DiskType string                  `json:"disk_type"`
	Id       string                  `json:"id"`
	Iops     ListResultTypesItemIops `json:"iops"`
	Name     string                  `json:"name"`
	Status   string                  `json:"status"`
}

type ListResultTypesItemIops struct {
	Read  int `json:"read"`
	Write int `json:"write"`
}

type ListResultTypes []ListResultTypesItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/block-storage/volume-types/list"), client, ctx)
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
