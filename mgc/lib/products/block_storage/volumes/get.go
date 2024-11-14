/*
Executor: get

# Summary

Retrieve the details of a specific volume.

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
	Attachment        *GetResultAttachment       `json:"attachment,omitempty"`
	AvailabilityZones GetResultAvailabilityZones `json:"availability_zones"`
	CreatedAt         string                     `json:"created_at"`
	Error             *GetResultError            `json:"error,omitempty"`
	Id                string                     `json:"id"`
	Name              string                     `json:"name"`
	Size              int                        `json:"size"`
	State             string                     `json:"state"`
	Status            string                     `json:"status"`
	Type              GetResultType              `json:"type"`
	UpdatedAt         string                     `json:"updated_at"`
}

type GetResultAttachment struct {
	AttachedAt string                      `json:"attached_at"`
	Device     *string                     `json:"device,omitempty"`
	Instance   GetResultAttachmentInstance `json:"instance"`
}

// any of: GetResultAttachmentInstance
type GetResultAttachmentInstance struct {
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type GetResultAvailabilityZones []string

type GetResultError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// any of: GetResultType
type GetResultType struct {
	DiskType *string            `json:"disk_type,omitempty"`
	Id       string             `json:"id"`
	Iops     *GetResultTypeIops `json:"iops,omitempty"`
	Name     *string            `json:"name,omitempty"`
	Status   *string            `json:"status,omitempty"`
}

type GetResultTypeIops struct {
	Read  int `json:"read"`
	Total int `json:"total"`
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) GetContext(
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/block-storage/volumes/get"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
