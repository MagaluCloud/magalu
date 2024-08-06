/*
Executor: list

# Description

# List all objects from a bucket

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	ContinuationToken *string               `json:"continuation-token,omitempty"`
	Dst               string                `json:"dst"`
	Filter            *ListParametersFilter `json:"filter,omitempty"`
	MaxItems          *int                  `json:"max-items,omitempty"`
	Recursive         *bool                 `json:"recursive,omitempty"`
}

type ListParametersFilterItem struct {
	Exclude *string `json:"exclude,omitempty"`
	Include *string `json:"include,omitempty"`
}

type ListParametersFilter []ListParametersFilterItem

type ListConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type ListResult struct {
	CommonPrefixes ListResultCommonPrefixes `json:"CommonPrefixes"`
	Contents       ListResultContents       `json:"Contents"`
}

type ListResultCommonPrefixesItem struct {
	Path string `json:"Path"`
}

type ListResultCommonPrefixes []ListResultCommonPrefixesItem

type ListResultContentsItem struct {
	ContentSize  int    `json:"ContentSize"`
	Key          string `json:"Key"`
	LastModified string `json:"LastModified"`
	StorageClass string `json:"StorageClass"`
}

type ListResultContents []ListResultContentsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/object-storage/objects/list"), s.client, s.ctx)
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
