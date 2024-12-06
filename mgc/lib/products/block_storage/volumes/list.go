/*
Executor: list

# Summary

List all volumes.

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
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit  *int                  `json:"_limit,omitempty"`
	Offset *int                  `json:"_offset,omitempty"`
	Sort   *string               `json:"_sort,omitempty"`
	Expand *ListParametersExpand `json:"expand,omitempty"`
	Name   *string               `json:"name,omitempty"`
}

type ListParametersExpand []string

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Volumes ListResultVolumes `json:"volumes"`
}

type ListResultVolumesItem struct {
	Attachment       *ListResultVolumesItemAttachment `json:"attachment,omitempty"`
	AvailabilityZone string                           `json:"availability_zone"`
	CreatedAt        string                           `json:"created_at"`
	Error            *ListResultVolumesItemError      `json:"error,omitempty"`
	Id               string                           `json:"id"`
	Name             string                           `json:"name"`
	Size             int                              `json:"size"`
	State            string                           `json:"state"`
	Status           string                           `json:"status"`
	Type             ListResultVolumesItemType        `json:"type"`
	UpdatedAt        string                           `json:"updated_at"`
}

type ListResultVolumesItemAttachment struct {
	AttachedAt string                                  `json:"attached_at"`
	Device     *string                                 `json:"device,omitempty"`
	Instance   ListResultVolumesItemAttachmentInstance `json:"instance"`
}

// any of: ListResultVolumesItemAttachmentInstance
type ListResultVolumesItemAttachmentInstance struct {
	CreatedAt string                                       `json:"created_at"`
	DiskType  *string                                      `json:"disk_type,omitempty"`
	Id        string                                       `json:"id"`
	Iops      *ListResultVolumesItemAttachmentInstanceIops `json:"iops,omitempty"`
	Name      string                                       `json:"name"`
	State     string                                       `json:"state"`
	Status    string                                       `json:"status"`
	UpdatedAt string                                       `json:"updated_at"`
}

type ListResultVolumesItemAttachmentInstanceIops struct {
	Read  int `json:"read"`
	Total int `json:"total"`
	Write int `json:"write"`
}

type ListResultVolumesItemError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// any of: ListResultVolumesItemType
type ListResultVolumesItemType struct {
	DiskType *string                                      `json:"disk_type,omitempty"`
	Id       string                                       `json:"id"`
	Iops     *ListResultVolumesItemAttachmentInstanceIops `json:"iops,omitempty"`
	Name     *string                                      `json:"name,omitempty"`
	Status   *string                                      `json:"status,omitempty"`
}

type ListResultVolumes []ListResultVolumesItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/block-storage/volumes/list"), s.client, s.ctx)
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/block-storage/volumes/list"), s.client, ctx)
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
