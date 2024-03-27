/*
Executor: list-v1-types

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

type ListV1TypesParameters struct {
	Limit  int    `json:"_limit,omitempty"`
	Offset int    `json:"_offset,omitempty"`
	Sort   string `json:"_sort,omitempty"`
}

type ListV1TypesConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListV1TypesResult struct {
	InstanceTypes ListV1TypesResultInstanceTypes `json:"instance_types"`
}

type ListV1TypesResultInstanceTypesItem struct {
	Disk   int    `json:"disk"`
	Gpu    int    `json:"gpu,omitempty"`
	Id     string `json:"id"`
	Name   string `json:"name"`
	Ram    int    `json:"ram"`
	Status string `json:"status"`
	Vcpus  int    `json:"vcpus"`
}

type ListV1TypesResultInstanceTypes []ListV1TypesResultInstanceTypesItem

func ListV1Types(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListV1TypesParameters,
	configs ListV1TypesConfigs,
) (
	result ListV1TypesResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("ListV1Types", mgcCore.RefPath("/virtual-machine/machine-types/list-v1-types"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListV1TypesParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ListV1TypesConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListV1TypesResult](r)
}

// TODO: links
// TODO: related
