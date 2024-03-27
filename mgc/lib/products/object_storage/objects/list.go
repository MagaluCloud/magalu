/*
Executor: list

# Description

# List all objects from a bucket

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	ContinuationToken string               `json:"continuation-token,omitempty"`
	Dst               string               `json:"dst"`
	Filter            ListParametersFilter `json:"filter,omitempty"`
	MaxItems          int                  `json:"max-items,omitempty"`
	Recursive         bool                 `json:"recursive,omitempty"`
}

type ListParametersFilterItem struct {
	Exclude string `json:"exclude,omitempty"`
	Include string `json:"include,omitempty"`
}

type ListParametersFilter []ListParametersFilterItem

type ListConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
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
}

type ListResultContents []ListResultContentsItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/object-storage/objects/list"), client, ctx)
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
