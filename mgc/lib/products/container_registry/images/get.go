/*
Executor: get

# Summary

# Get image details

# Description

Show detailed information about the image.

Version: 0.1.0

import "github.com/MagaluCloud/magalu/mgc/lib/products/container_registry/images"
*/
package images

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type GetParameters struct {
	DigestOrTag    string `json:"digest_or_tag"`
	RegistryId     string `json:"registry_id"`
	RepositoryName string `json:"repository_name"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Repository image response data.
type GetResult struct {
	Digest            string                `json:"digest"`
	ExtraAttr         *GetResultExtraAttr   `json:"extra_attr,omitempty"`
	ManifestMediaType *string               `json:"manifest_media_type,omitempty"`
	MediaType         *string               `json:"media_type,omitempty"`
	PulledAt          string                `json:"pulled_at"`
	PushedAt          string                `json:"pushed_at"`
	SizeBytes         int                   `json:"size_bytes"`
	Tags              GetResultTags         `json:"tags"`
	TagsDetails       *GetResultTagsDetails `json:"tags_details,omitempty"`
}

// Extra attributes about the image.
type GetResultExtraAttr struct {
}

type GetResultTags []string

// Tag of an image response.
type GetResultTagsDetailsItem struct {
	Name     *string `json:"name,omitempty"`
	PulledAt *string `json:"pulled_at,omitempty"`
	PushedAt *string `json:"pushed_at,omitempty"`
	Signed   *bool   `json:"signed,omitempty"`
}

type GetResultTagsDetails []GetResultTagsDetailsItem

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/container-registry/images/get"), s.client, s.ctx)
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/container-registry/images/get"), s.client, ctx)
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
