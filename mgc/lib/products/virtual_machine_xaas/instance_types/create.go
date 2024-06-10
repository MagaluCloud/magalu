/*
Executor: create

# Summary

# Create Instance Type V1

# Description

This internal method is for create a new flavor (instance type)

	some information will pass to URP, but some information

	like ratio or generation for now will be only in the database.

	For gpu will have 2 infos:

	1- the gpu model name (today will olny have a100 as model)

	2 - number of gpu cores this flavor will use.

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/instance_types"
*/
package instanceTypes

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Disk     int     `json:"disk"`
	Gpu      *int    `json:"gpu,omitempty"`
	GpuModel *string `json:"gpu_model,omitempty"`
	Id       string  `json:"id"`
	Internal bool    `json:"internal"`
	Name     string  `json:"name"`
	Ram      int     `json:"ram"`
	Sku      string  `json:"sku"`
	Vcpus    int     `json:"vcpus"`
}

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
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

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine-xaas/instance types/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
