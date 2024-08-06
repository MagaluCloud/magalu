/*
Executor: list

# Summary

# List images in container registry repository

# Description

# List all images in container registry repository

Version: 0.1.0

import "magalu.cloud/lib/products/container_registry/images"
*/
package images

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit          *int                  `json:"_limit,omitempty"`
	Offset         *int                  `json:"_offset,omitempty"`
	Sort           *string               `json:"_sort,omitempty"`
	Expand         *ListParametersExpand `json:"expand,omitempty"`
	RegistryId     string                `json:"registry_id"`
	RepositoryName string                `json:"repository_name"`
}

type ListParametersExpand []string

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Repository images response.
type ListResult struct {
	Results ListResultResults `json:"results"`
}

// Repository image response data.
type ListResultResultsItem struct {
	Digest            string                            `json:"digest"`
	ExtraAttr         *ListResultResultsItemExtraAttr   `json:"extra_attr,omitempty"`
	ManifestMediaType *string                           `json:"manifest_media_type,omitempty"`
	MediaType         *string                           `json:"media_type,omitempty"`
	PulledAt          string                            `json:"pulled_at"`
	PushedAt          string                            `json:"pushed_at"`
	SizeBytes         int                               `json:"size_bytes"`
	Tags              ListResultResultsItemTags         `json:"tags"`
	TagsDetails       *ListResultResultsItemTagsDetails `json:"tags_details,omitempty"`
}

// Extra attributes about the image.
type ListResultResultsItemExtraAttr struct {
}

type ListResultResultsItemTags []string

// Tag of an image response.
type ListResultResultsItemTagsDetailsItem struct {
	Name     *string `json:"name,omitempty"`
	PulledAt *string `json:"pulled_at,omitempty"`
	PushedAt *string `json:"pushed_at,omitempty"`
	Signed   *bool   `json:"signed,omitempty"`
}

type ListResultResultsItemTagsDetails []ListResultResultsItemTagsDetailsItem

type ListResultResults []ListResultResultsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/container-registry/images/list"), s.client, s.ctx)
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
