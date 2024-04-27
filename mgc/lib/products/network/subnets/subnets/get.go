/*
Executor: get

# Summary

# Subnet Details

# Description

# Returns a subnet details

Version: 1.119.0

import "magalu.cloud/lib/products/network/subnets/subnets"
*/
package subnets

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	SubnetId string `json:"subnet_id"`
}

type GetConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	CidrBlock      string                  `json:"cidr_block"`
	CreatedAt      *string                 `json:"created_at,omitempty"`
	Description    *string                 `json:"description,omitempty"`
	DhcpPools      GetResultDhcpPools      `json:"dhcp_pools"`
	DnsNameservers GetResultDnsNameservers `json:"dns_nameservers"`
	GatewayIp      string                  `json:"gateway_ip"`
	Id             string                  `json:"id"`
	IpVersion      string                  `json:"ip_version"`
	Name           *string                 `json:"name,omitempty"`
	Updated        *string                 `json:"updated,omitempty"`
	VpcId          string                  `json:"vpc_id"`
	Zone           string                  `json:"zone"`
}

type GetResultDhcpPoolsItem struct {
	End   string `json:"end"`
	Start string `json:"start"`
}

type GetResultDhcpPools []GetResultDhcpPoolsItem

type GetResultDnsNameservers []string

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/network/subnets/subnets/get"), s.client, s.ctx)
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
