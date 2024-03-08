/*
Executor: get

# Summary

# Retrieve the details of a volume

# Description

Retrieve details of a Volume for the currently authenticated tenant.

#### Notes
  - Use the **expand** argument to obtain additional details about the Volume
    Type.
  - Utilize the **block-storage volume list** command to retrieve a list of all
    Volumes and obtain the ID of the Volume for which you want to retrieve
    details.

Version: v1

import "magalu.cloud/lib/products/block_storage/volumes"
*/
package volumes

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	Expand GetParametersExpand `json:"expand,omitempty"`
	Id     string              `json:"id"`
}

type GetParametersExpand []string

type GetConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	Attachment *GetResultAttachment `json:"attachment,omitempty"`
	CreatedAt  string               `json:"created_at"`
	Error      GetResultError       `json:"error,omitempty"`
	Id         string               `json:"id"`
	Name       string               `json:"name"`
	Size       int                  `json:"size"`
	State      string               `json:"state"`
	Status     string               `json:"status"`
	Type       GetResultType        `json:"type"`
	UpdatedAt  string               `json:"updated_at"`
}

// any of: GetResultAttachment0, GetResultAttachment1
type GetResultAttachment struct {
	GetResultAttachment0 `json:",squash"` // nolint
	GetResultAttachment1 `json:",squash"` // nolint
}

type GetResultAttachment0 struct {
	AttachedAt string `json:"attached_at"`
	MachineId  string `json:"machine_id"`
}

type GetResultAttachment1 struct {
	AttachedAt string                      `json:"attached_at"`
	Machine    GetResultAttachment1Machine `json:"machine"`
	MachineId  string                      `json:"machine_id"`
}

type GetResultAttachment1Machine struct {
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	State     string `json:"state"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type GetResultError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// any of: GetResultType0, GetResultType1
type GetResultType struct {
	GetResultType0 `json:",squash"` // nolint
	GetResultType1 `json:",squash"` // nolint
}

type GetResultType0 struct {
	Id string `json:"id"`
}

type GetResultType1 struct {
	DiskType string             `json:"disk_type"`
	Id       string             `json:"id"`
	Iops     GetResultType1Iops `json:"iops"`
	Name     string             `json:"name"`
	Status   string             `json:"status"`
}

type GetResultType1Iops struct {
	Read  int `json:"read"`
	Write int `json:"write"`
}

func Get(
	client *mgcClient.Client,
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/block-storage/volumes/get"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[GetConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

func GetUntilTermination(
	client *mgcClient.Client,
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	e, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/block-storage/volumes/get"), client, ctx)
	if err != nil {
		return
	}

	exec, ok := e.(mgcCore.TerminatorExecutor)
	if !ok {
		// Not expected, but let's fallback
		return Get(
			client,
			ctx,
			parameters,
			configs,
		)
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[GetConfigs](configs); err != nil {
		return
	}

	r, err := exec.ExecuteUntilTermination(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
