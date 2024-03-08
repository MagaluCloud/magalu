/*
Executor: list

# Summary

Retrieves all machine-types available in the region.

# Description

Retrieves a list of machine types allowed for the current tenant which is logged in.

Version: 0.1.0

import "magalu.cloud/lib/products/virtual_machine/machine_types"
*/
package machineTypes

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
	InstanceTypes ListResultInstanceTypes `json:"instance_types"`
}

type ListResultInstanceTypesItem struct {
	Disk   int    `json:"disk"`
	Gpu    int    `json:"gpu,omitempty"`
	Id     string `json:"id"`
	Name   string `json:"name"`
	Ram    int    `json:"ram"`
	Status string `json:"status"`
	Vcpus  int    `json:"vcpus"`
}

type ListResultInstanceTypes []ListResultInstanceTypesItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/machine-types/list"), client, ctx)
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
