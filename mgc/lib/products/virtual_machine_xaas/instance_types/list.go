/*
Executor: list

# Summary

# List Instance Types Internal V1

# Description

Internal list all instance types this route will return even the internal instance types.

### Note
This route is used only for internal proposes.

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/instance_types"
*/
package instanceTypes

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	InstanceTypes ListResultInstanceTypes `json:"instance_types"`
}

type ListResultInstanceTypesItem struct {
	Disk     int    `json:"disk"`
	Gpu      *int   `json:"gpu,omitempty"`
	Id       string `json:"id"`
	Internal bool   `json:"internal"`
	Name     string `json:"name"`
	Ram      int    `json:"ram"`
	Sku      string `json:"sku"`
	Status   string `json:"status"`
	Vcpus    int    `json:"vcpus"`
}

type ListResultInstanceTypes []ListResultInstanceTypesItem

func (s *service) List(
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine-xaas/instance types/list"), s.client, s.ctx)
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
