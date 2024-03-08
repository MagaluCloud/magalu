/*
Executor: detach

# Summary

# Detach Security Group

# Description

# Detach a Security Group to a Port

Version: 1.109.0

import "magalu.cloud/lib/products/network/ports"
*/
package ports

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DetachParameters struct {
	PortId          string `json:"port_id"`
	SecurityGroupId string `json:"security_group_id"`
}

type DetachConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type DetachResult any

func Detach(
	client *mgcClient.Client,
	ctx context.Context,
	parameters DetachParameters,
	configs DetachConfigs,
) (
	result DetachResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Detach", mgcCore.RefPath("/network/ports/detach"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DetachParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DetachConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DetachResult](r)
}

// TODO: links
// TODO: related
