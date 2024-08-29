/*
Executor: get

# Summary

Retrieve the details of an snapshot.

# Description

Get an snapshot details for the current tenant which is logged in.

#### Notes
- You can use the snapshots list command to retrieve all snapshots,
so you can get the id of the snapshot that you want to get details.

- You can use the **expand** argument to get more details from the inner objects
like image or type.

Version: v1

import "magalu.cloud/lib/products/virtual_machine/snapshots"
*/
package snapshots

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	Expand *GetParametersExpand `json:"expand,omitempty"`
	Id     string               `json:"id"`
}

type GetParametersExpand []string

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	CreatedAt string            `json:"created_at"`
	Id        string            `json:"id"`
	Instance  GetResultInstance `json:"instance"`
	Name      *string           `json:"name,omitempty"`
	Size      int               `json:"size"`
	State     string            `json:"state"`
	Status    string            `json:"status"`
	UpdatedAt *string           `json:"updated_at,omitempty"`
}

type GetResultInstance struct {
	Id          string                       `json:"id"`
	Image       GetResultInstanceImage       `json:"image"`
	MachineType GetResultInstanceMachineType `json:"machine_type"`
}

// any of: GetResultInstanceImage
type GetResultInstanceImage struct {
	Id       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	Platform *string `json:"platform,omitempty"`
}

// any of: GetResultInstanceMachineType
type GetResultInstanceMachineType struct {
	Disk  *int    `json:"disk,omitempty"`
	Id    string  `json:"id"`
	Name  *string `json:"name,omitempty"`
	Ram   *int    `json:"ram,omitempty"`
	Vcpus *int    `json:"vcpus,omitempty"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine/snapshots/get"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[map[string]interface{}](s.client.Sdk().Config().TempConfig()); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
