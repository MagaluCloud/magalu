/*
Executor: list

# Summary

# Details of a Port list

# Description

Return a detailed list of ports from the provided tenant_id

Version: 1.124.1

import "magalu.cloud/lib/products/network/ports"
*/
package ports

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	PortIdList *ListParametersPortIdList `json:"port_id_list,omitempty"`
}

type ListParametersPortIdList []string

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResultItem struct {
	CreatedAt             *string                       `json:"created_at,omitempty"`
	Description           *string                       `json:"description,omitempty"`
	Id                    *string                       `json:"id,omitempty"`
	IpAddress             *ListResultItemIpAddress      `json:"ip_address,omitempty"`
	IsAdminStateUp        *bool                         `json:"is_admin_state_up,omitempty"`
	IsPortSecurityEnabled *bool                         `json:"is_port_security_enabled,omitempty"`
	Name                  *string                       `json:"name,omitempty"`
	PublicIp              *ListResultItemPublicIp       `json:"public_ip,omitempty"`
	SecurityGroups        *ListResultItemSecurityGroups `json:"security_groups,omitempty"`
	Updated               *string                       `json:"updated,omitempty"`
	VpcId                 *string                       `json:"vpc_id,omitempty"`
}

type ListResultItemIpAddressItem struct {
	IpAddress string `json:"ip_address"`
	SubnetId  string `json:"subnet_id"`
}

type ListResultItemIpAddress []ListResultItemIpAddressItem

type ListResultItemPublicIpItem struct {
	PublicIp   *string `json:"public_ip,omitempty"`
	PublicIpId *string `json:"public_ip_id,omitempty"`
}

type ListResultItemPublicIp []ListResultItemPublicIpItem

type ListResultItemSecurityGroups []string

type ListResult []ListResultItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/ports/list"), s.client, s.ctx)
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
