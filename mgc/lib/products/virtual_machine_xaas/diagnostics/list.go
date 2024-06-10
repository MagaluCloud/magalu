/*
Executor: list

# Summary

# List Instances Diagnostic

# Description

Get a list of instances informing a tenant ID.

This route only list instances inserted on virtual machine database.

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/diagnostics"
*/
package diagnostics

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit       *int    `json:"_limit,omitempty"`
	Offset      *int    `json:"_offset,omitempty"`
	Sort        *string `json:"_sort,omitempty"`
	ProjectType *string `json:"project_type,omitempty"`
}

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Instances ListResultInstances `json:"instances"`
}

type ListResultInstancesItem struct {
	CreatedAt   string                       `json:"created_at"`
	Error       *string                      `json:"error,omitempty"`
	Id          string                       `json:"id"`
	ImageName   string                       `json:"image_name"`
	KeyPairName string                       `json:"key_pair_name"`
	Name        string                       `json:"name"`
	Ports       ListResultInstancesItemPorts `json:"ports"`
	ProjectType string                       `json:"project_type"`
	Retries     int                          `json:"retries"`
	SnapshotId  *string                      `json:"snapshot_id,omitempty"`
	Status      string                       `json:"status"`
	Step        int                          `json:"step"`
	UpdatedAt   *string                      `json:"updated_at,omitempty"`
	UserData    *string                      `json:"user_data,omitempty"`
}

type ListResultInstancesItemPortsItem struct {
	Port string `json:"port"`
}

type ListResultInstancesItemPorts []ListResultInstancesItemPortsItem

type ListResultInstances []ListResultInstancesItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/virtual-machine-xaas/diagnostics/list"), s.client, s.ctx)
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
