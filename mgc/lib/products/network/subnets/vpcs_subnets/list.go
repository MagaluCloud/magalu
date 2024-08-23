/*
Executor: list

# Summary

# List Subnets

# Description

Returns a list of subnets for a provided vpc_id

Version: 1.131.1

import "magalu.cloud/lib/products/network/subnets/vpcs_subnets"
*/
package vpcsSubnets

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	VpcId string `json:"vpc_id"`
}

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Subnets *ListResultSubnets `json:"subnets,omitempty"`
}

type ListResultSubnetsItem struct {
	CidrBlock   string  `json:"cidr_block"`
	CreatedAt   *string `json:"created_at,omitempty"`
	Description *string `json:"description,omitempty"`
	Id          string  `json:"id"`
	IpVersion   string  `json:"ip_version"`
	Name        *string `json:"name,omitempty"`
	Updated     *string `json:"updated,omitempty"`
	VpcId       string  `json:"vpc_id"`
	Zone        string  `json:"zone"`
}

type ListResultSubnets []ListResultSubnetsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/subnets/vpcs-subnets/list"), s.client, s.ctx)
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
