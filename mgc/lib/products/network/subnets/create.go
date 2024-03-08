/*
Executor: create

# Summary

# Create Subnet

# Description

# Create a Subnet

Version: 1.109.0

import "magalu.cloud/lib/products/network/subnets"
*/
package subnets

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	CidrBlock      *string                        `json:"cidr_block,omitempty"`
	Description    *string                        `json:"description,omitempty"`
	DnsNameservers CreateParametersDnsNameservers `json:"dns_nameservers,omitempty"`
	IpVersion      int                            `json:"ip_version"`
	Name           string                         `json:"name"`
	VpcId          string                         `json:"vpc_id"`
}

type CreateParametersDnsNameservers []string

type CreateConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
	Id string `json:"id"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/network/subnets/create"), client, ctx)
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
