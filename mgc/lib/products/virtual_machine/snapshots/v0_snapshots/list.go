/*
Executor: list

# Summary

# List Snapshots

# Description

# List snapshots by tenant

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/snapshots/v0_snapshots"
*/
package v0Snapshots

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
	Results ListResultResults `json:"results"`
}

type ListResultResultsItem struct {
	CreatedAt        string                     `json:"created_at"`
	Error            *string                    `json:"error,omitempty"`
	Id               string                     `json:"id"`
	Image            ListResultResultsItemImage `json:"image"`
	InstanceTypeId   string                     `json:"instance_type_id"`
	Name             string                     `json:"name"`
	Size             *int                       `json:"size,omitempty"`
	Status           string                     `json:"status"`
	TenantId         string                     `json:"tenant_id"`
	UpdatedAt        *string                    `json:"updated_at,omitempty"`
	VirtualMachineId string                     `json:"virtual_machine_id"`
}

type ListResultResultsItemImage struct {
	Name string `json:"name"`
}

type ListResultResults []ListResultResultsItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine/snapshots/v0-snapshots/list"), client, ctx)
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
