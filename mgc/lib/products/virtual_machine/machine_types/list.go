/*
Executor: list

# Summary

Retrieves all machine-types.

# Description

Retrieves a list of machine types allowed for the current tenant which is logged in.

Version: v1

import "magalu.cloud/lib/products/virtual_machine/machine_types"
*/
package machineTypes

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit            *int    `json:"_limit,omitempty"`
	Offset           *int    `json:"_offset,omitempty"`
	Sort             *string `json:"_sort,omitempty"`
	AvailabilityZone *string `json:"availability-zone,omitempty"`
}

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	MachineTypes ListResultMachineTypes `json:"machine_types"`
}

// any of: ListResultMachineTypesItem
type ListResultMachineTypesItem struct {
	AvailabilityZones *ListResultMachineTypesItemAvailabilityZones `json:"availability_zones,omitempty"`
	Disk              int                                          `json:"disk"`
	Gpu               *int                                         `json:"gpu,omitempty"`
	Id                string                                       `json:"id"`
	Name              string                                       `json:"name"`
	Ram               int                                          `json:"ram"`
	Sku               *string                                      `json:"sku"`
	Status            string                                       `json:"status"`
	Vcpus             int                                          `json:"vcpus"`
}

type ListResultMachineTypesItemAvailabilityZones []string

type ListResultMachineTypes []ListResultMachineTypesItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/machine-types/list"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ListContext(
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/machine-types/list"), s.client, ctx)
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

	sdkConfig := s.client.Sdk().Config().TempConfig()
	if c["serverUrl"] == nil && sdkConfig["serverUrl"] != nil {
		c["serverUrl"] = sdkConfig["serverUrl"]
	}

	if c["env"] == nil && sdkConfig["env"] != nil {
		c["env"] = sdkConfig["env"]
	}

	if c["region"] == nil && sdkConfig["region"] != nil {
		c["region"] = sdkConfig["region"]
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
