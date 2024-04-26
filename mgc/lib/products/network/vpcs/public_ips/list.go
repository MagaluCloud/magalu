/*
Executor: list

# Summary

# List Public IPs

# Description

Returns a list of Public IPs for a provided vpc_id

Version: 1.114.3

import "magalu.cloud/lib/products/network/vpcs/public_ips"
*/
package publicIps

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	VpcId string `json:"vpc_id"`
}

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

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/vpcs/public-ips/list"), s.client, s.ctx)
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
