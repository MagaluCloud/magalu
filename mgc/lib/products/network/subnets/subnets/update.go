/*
Executor: update

# Summary

# Update Subnet

# Description

Update a subnet from the provided tenant_id

Version: 1.131.0

import "magalu.cloud/lib/products/network/subnets/subnets"
*/
package subnets

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UpdateParameters struct {
	DnsNameservers *UpdateParametersDnsNameservers `json:"dns_nameservers,omitempty"`
	SubnetId       string                          `json:"subnet_id"`
}

type UpdateParametersDnsNameservers []string

type UpdateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type UpdateResult struct {
	Id string `json:"id"`
}

func (s *service) Update(
	parameters UpdateParameters,
	configs UpdateConfigs,
) (
	result UpdateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/network/subnets/subnets/update"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UpdateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UpdateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[UpdateResult](r)
}

// TODO: links
// TODO: related
