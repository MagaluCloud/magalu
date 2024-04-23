/*
Executor: list

# Summary

# List all Volumes

# Description

Retrieve a list of Volumes for the currently authenticated tenant.

#### Notes
- Use the expand argument to obtain additional details about the Volume Type.

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

type ListParameters struct {
	Limit  int                  `json:"_limit,omitempty"`
	Offset int                  `json:"_offset,omitempty"`
	Sort   string               `json:"_sort,omitempty"`
	Expand ListParametersExpand `json:"expand,omitempty"`
}

type ListParametersExpand []string

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Volumes ListResultVolumes `json:"volumes"`
}

type ListResultVolumesItem struct {
	Attachment ListResultVolumesItemAttachment `json:"attachment,omitempty"`
	CreatedAt  string                          `json:"created_at"`
	Error      ListResultVolumesItemError      `json:"error,omitempty"`
	Id         string                          `json:"id"`
	Name       string                          `json:"name"`
	Size       int                             `json:"size"`
	State      string                          `json:"state"`
	Status     string                          `json:"status"`
	Type       ListResultVolumesItemType       `json:"type"`
	UpdatedAt  string                          `json:"updated_at"`
}

type ListResultVolumesItemAttachment struct {
	AttachedAt string                                  `json:"attached_at"`
	Device     string                                  `json:"device,omitempty"`
	Instance   ListResultVolumesItemAttachmentInstance `json:"instance"`
}

// any of: ListResultVolumesItemAttachmentInstance0, ListResultVolumesItemAttachmentInstance1
type ListResultVolumesItemAttachmentInstance struct {
	ListResultVolumesItemAttachmentInstance0 `json:",squash"` // nolint
	ListResultVolumesItemAttachmentInstance1 `json:",squash"` // nolint
}

type ListResultVolumesItemAttachmentInstance0 struct {
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type ListResultVolumesItemAttachmentInstance1 struct {
	Id string `json:"id"`
}

type ListResultVolumesItemError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// any of: ListResultVolumesItemAttachmentInstance1, ListResultVolumesItemType1
type ListResultVolumesItemType struct {
	ListResultVolumesItemAttachmentInstance1 `json:",squash"` // nolint
	ListResultVolumesItemType1               `json:",squash"` // nolint
}

type ListResultVolumesItemType1 struct {
	DiskType string                         `json:"disk_type"`
	Id       string                         `json:"id"`
	Iops     ListResultVolumesItemType1Iops `json:"iops"`
	Name     string                         `json:"name"`
	Status   string                         `json:"status"`
}

type ListResultVolumesItemType1Iops struct {
	Read  int `json:"read"`
	Write int `json:"write"`
}

type ListResultVolumes []ListResultVolumesItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/block-storage/volumes/list"), client, ctx)
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
