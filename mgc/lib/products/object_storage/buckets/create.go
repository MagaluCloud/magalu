/*
Executor: create

# Summary

# Create a new Bucket

# Description

Buckets are "containers" that are able to store various Objects inside

import "magalu.cloud/lib/products/object_storage/buckets"
*/
package buckets

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	AuthenticatedRead bool                             `json:"authenticated_read,omitempty"`
	EnableVersioning  bool                             `json:"enable_versioning,omitempty"`
	GrantFullControl  CreateParametersGrantFullControl `json:"grant_full_control,omitempty"`
	GrantRead         CreateParametersGrantRead        `json:"grant_read,omitempty"`
	GrantReadAcp      CreateParametersGrantReadAcp     `json:"grant_read_acp,omitempty"`
	GrantWrite        CreateParametersGrantWrite       `json:"grant_write,omitempty"`
	GrantWriteAcp     CreateParametersGrantWriteAcp    `json:"grant_write_acp,omitempty"`
	Location          string                           `json:"location,omitempty"`
	Name              string                           `json:"name"`
	Private           bool                             `json:"private,omitempty"`
	PublicRead        bool                             `json:"public_read,omitempty"`
	PublicReadWrite   bool                             `json:"public_read_write,omitempty"`
}

type CreateParametersGrantFullControlItem struct {
	Id string `json:"id"`
}

type CreateParametersGrantFullControl []CreateParametersGrantFullControlItem

type CreateParametersGrantRead []CreateParametersGrantFullControlItem

type CreateParametersGrantReadAcp []CreateParametersGrantFullControlItem

type CreateParametersGrantWrite []CreateParametersGrantFullControlItem

type CreateParametersGrantWriteAcp []CreateParametersGrantFullControlItem

type CreateConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type CreateResult any

func Create(
	client *mgcClient.Client,
	ctx context.Context,
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/object-storage/buckets/create"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
