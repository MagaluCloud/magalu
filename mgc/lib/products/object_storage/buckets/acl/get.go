/*
Executor: get

# Description

# Get the ACL for the specified bucket

import "magalu.cloud/lib/products/object_storage/buckets/acl"
*/
package acl

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	Bucket string `json:"bucket"`
}

type GetConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type GetResult struct {
	AccessControlList GetResultAccessControlList `json:"AccessControlList"`
	Owner             GetResultOwner             `json:"Owner"`
}

type GetResultAccessControlList struct {
	Grant GetResultAccessControlListGrant `json:"Grant"`
}

type GetResultAccessControlListGrant struct {
	Grantee    GetResultAccessControlListGrantGrantee `json:"Grantee"`
	Permission string                                 `json:"Permission"`
}

type GetResultAccessControlListGrantGrantee struct {
	DisplayName  string `json:"DisplayName"`
	EmailAddress string `json:"EmailAddress"`
	Id           string `json:"ID"`
	Uri          string `json:"URI"`
}

type GetResultOwner struct {
	DisplayName string `json:"DisplayName"`
	Id          string `json:"ID"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/object-storage/buckets/acl/get"), client, ctx)
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

// TODO: links
// TODO: related
