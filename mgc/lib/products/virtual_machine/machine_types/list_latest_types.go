/*
Executor: list-latest-types

# Summary

Retrieves all machine-types available in the region.

# Description

Retrieves a list of machine types allowed for the current tenant which is logged in.

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/machine_types"
*/
package machineTypes

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListLatestTypesParameters struct {
	Limit  int    `json:"_limit,omitempty"`
	Offset int    `json:"_offset,omitempty"`
	Sort   string `json:"_sort,omitempty"`
}

type ListLatestTypesConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListLatestTypesResult struct {
	InstanceTypes ListLatestTypesResultInstanceTypes `json:"instance_types"`
}

type ListLatestTypesResultInstanceTypesItem struct {
	Disk   int    `json:"disk"`
	Gpu    int    `json:"gpu,omitempty"`
	Id     string `json:"id"`
	Name   string `json:"name"`
	Ram    int    `json:"ram"`
	Status string `json:"status"`
	Vcpus  int    `json:"vcpus"`
}

type ListLatestTypesResultInstanceTypes []ListLatestTypesResultInstanceTypesItem

func ListLatestTypes(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListLatestTypesParameters,
	configs ListLatestTypesConfigs,
) (
	result ListLatestTypesResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("ListLatestTypes", mgcCore.RefPath("/virtual-machine/machine-types/list-latest-types"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListLatestTypesParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ListLatestTypesConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListLatestTypesResult](r)
}

// TODO: links
// TODO: related
