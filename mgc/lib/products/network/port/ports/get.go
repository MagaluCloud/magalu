/*
Executor: get

# Summary

# Port Details

# Description

Return a port details from the provided tenant_id

Version: 1.130.0

import "magalu.cloud/lib/products/network/port/ports"
*/
package ports

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	PortId string `json:"port_id"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	CreatedAt             *string                  `json:"created_at,omitempty"`
	Description           *string                  `json:"description,omitempty"`
	Id                    *string                  `json:"id,omitempty"`
	IpAddress             *GetResultIpAddress      `json:"ip_address,omitempty"`
	IsAdminStateUp        *bool                    `json:"is_admin_state_up,omitempty"`
	IsPortSecurityEnabled *bool                    `json:"is_port_security_enabled,omitempty"`
	Name                  *string                  `json:"name,omitempty"`
	PublicIp              *GetResultPublicIp       `json:"public_ip,omitempty"`
	SecurityGroups        *GetResultSecurityGroups `json:"security_groups,omitempty"`
	Updated               *string                  `json:"updated,omitempty"`
	VpcId                 *string                  `json:"vpc_id,omitempty"`
}

type GetResultIpAddressItem struct {
	IpAddress string `json:"ip_address"`
	SubnetId  string `json:"subnet_id"`
}

type GetResultIpAddress []GetResultIpAddressItem

type GetResultPublicIpItem struct {
	PublicIp   *string `json:"public_ip,omitempty"`
	PublicIpId *string `json:"public_ip_id,omitempty"`
}

type GetResultPublicIp []GetResultPublicIpItem

type GetResultSecurityGroups []string

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/network/port/ports/get"), s.client, s.ctx)
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
