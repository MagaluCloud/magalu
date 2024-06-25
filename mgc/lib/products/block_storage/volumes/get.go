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
	Attachment *GetResultAttachment `json:"attachment,omitempty"`
	CreatedAt  string               `json:"created_at"`
	Error      *GetResultError      `json:"error,omitempty"`
	Id         string               `json:"id"`
	Name       string               `json:"name"`
	Size       int                  `json:"size"`
	State      string               `json:"state"`
	Status     string               `json:"status"`
	Type       GetResultType        `json:"type"`
	UpdatedAt  string               `json:"updated_at"`
}

type GetResultAttachment struct {
	AttachedAt string                      `json:"attached_at"`
	Device     *string                     `json:"device,omitempty"`
	Instance   GetResultAttachmentInstance `json:"instance"`
}

// any of: , GetResultAttachmentInstance1
type GetResultAttachmentInstance struct {
	GetResultAttachmentInstance1 `json:",squash"` // nolint
}

type GetResultAttachmentInstance1 struct {
	Id string `json:"id"`
}

type GetResultError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// any of: , GetResultType1
type GetResultType struct {
	GetResultType1 `json:",squash"` // nolint
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

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/block-storage/volumes/get"), s.client, s.ctx)
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

func (s *service) GetUntilTermination(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	e, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/block-storage/volumes/get"), s.client, s.ctx)
	if err != nil {
		return
	}

	exec, ok := e.(mgcCore.TerminatorExecutor)
	if !ok {
		// Not expected, but let's fallback
		return s.Get(
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
