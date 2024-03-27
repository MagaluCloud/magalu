/*
Executor: create

# Summary

# Create Image

# Description

# Creates images only in DB

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/images/v0_images"
*/
package v0Images

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Id       string `json:"id"`
	Internal bool   `json:"internal"`
	Sku      string `json:"sku"`
}

type CreateConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	OsDistribution *string `json:"os_distribution,omitempty"`
	Size           int     `json:"size"`
	Sku            *string `json:"sku,omitempty"`
	Version        *string `json:"version,omitempty"`
}

func Create(
	client *mgcClient.Client,
	ctx context.Context,
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine/images/v0-images/create"), client, ctx)
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
