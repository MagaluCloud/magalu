/*
Executor: get

# Summary

# Public IP Details

# Description

# Return a Public IP details

Version: 1.109.0

import "magalu.cloud/lib/products/network/public_ips"
*/
package publicIps

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	PublicIpId string `json:"public_ip_id"`
}

type GetConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	CreatedAt   *string `json:"created_at,omitempty"`
	Description *string `json:"description,omitempty"`
	Error       *string `json:"error,omitempty"`
	ExternalId  *string `json:"external_id,omitempty"`
	Id          *string `json:"id,omitempty"`
	PortId      *string `json:"port_id,omitempty"`
	PublicIp    *string `json:"public_ip,omitempty"`
	Status      *string `json:"status,omitempty"`
	Updated     *string `json:"updated,omitempty"`
	VpcId       *string `json:"vpc_id,omitempty"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/network/public_ips/get"), client, ctx)
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
