/*
Executor: list

# Description

# List all existing Buckets

import "magalu.cloud/lib/products/object_storage/buckets"
*/
package buckets

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type ListResult struct {
	Buckets ListResultBuckets `json:"Buckets"`
	Owner   ListResultOwner   `json:"Owner"`
}

type ListResultBucketsItem struct {
	CreationDate string `json:"CreationDate"`
	Name         string `json:"Name"`
}

type ListResultBuckets []ListResultBucketsItem

type ListResultOwner struct {
	DisplayName string `json:"DisplayName"`
	Id          string `json:"ID"`
}

func List(
	client *mgcClient.Client,
	ctx context.Context,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/object-storage/buckets/list"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

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
