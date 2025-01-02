/*
Executor: get

# Summary

# Port Details

# Description

Return a port details from the provided tenant_id

Version: 1.141.3

import "github.com/MagaluCloud/magalu/mgc/lib/products/network/ports"
*/
package ports

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
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
	Ethertype *string `json:"ethertype,omitempty"`
	IpAddress string  `json:"ip_address"`
	SubnetId  string  `json:"subnet_id"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/network/ports/get"), s.client, s.ctx)
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) GetContext(
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/network/ports/get"), s.client, ctx)
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

	sdkConfig := s.client.Sdk().Config().TempConfig()
	if c["serverUrl"] == nil && sdkConfig["serverUrl"] != nil {
		c["serverUrl"] = sdkConfig["serverUrl"]
	}

	if c["env"] == nil && sdkConfig["env"] != nil {
		c["env"] = sdkConfig["env"]
	}

	if c["region"] == nil && sdkConfig["region"] != nil {
		c["region"] = sdkConfig["region"]
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
