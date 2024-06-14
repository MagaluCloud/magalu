/*
Executor: create

# Summary

# Create Subnet

# Description

# Create a Subnet

Version: 1.124.1

import "magalu.cloud/lib/products/network/subnets/vpcs_subnets"
*/
package vpcsSubnets

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	CidrBlock      *string                         `json:"cidr_block,omitempty"`
	Description    *string                         `json:"description,omitempty"`
	DnsNameservers *CreateParametersDnsNameservers `json:"dns_nameservers,omitempty"`
	IpVersion      int                             `json:"ip_version"`
	Name           string                          `json:"name"`
	VpcId          any                             `json:"vpc_id"`
}

type CreateParametersDnsNameservers []string

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	XZone     *string `json:"x-zone,omitempty"`
}

type CreateResult struct {
	Id string `json:"id"`
}

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/network/subnets/vpcs-subnets/create"), s.client, s.ctx)
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
