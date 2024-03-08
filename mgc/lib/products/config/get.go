/*
Executor: get

# Summary

# Get a specific Config value that has been previously set

# Description

Get a specific Config value that has been previously set. If there's an env variable
matching the key (in uppercase and with the 'MGC_' prefix), it'll be retreived.
Otherwise, the value will be searched for in the YAML file

import "magalu.cloud/lib/products/config"
*/
package config

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	Key string `json:"key"`
}

type GetResult any

func Get(
	client *mgcClient.Client,
	ctx context.Context,
	parameters GetParameters,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/config/get"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
