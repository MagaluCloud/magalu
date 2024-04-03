/*
Executor: list

# Summary

# Tenant's public IP list

# Description

# Return a tenant's public ip list

Version: 1.114.3

import "magalu.cloud/lib/products/network/public_ips"
*/
package publicIps

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	PublicIps ListResultPublicIps `json:"public_ips"`
}

type ListResultPublicIpsItem struct {
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

type ListResultPublicIps []ListResultPublicIpsItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/public_ips/list"), client, ctx)
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
