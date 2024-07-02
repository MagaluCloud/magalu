/*
Executor: list

# Summary

# List Ports

# Description

# List VPC ports

Version: 1.126.1

import "magalu.cloud/lib/products/network/vpcs/ports"
*/
package ports

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit      *int                      `json:"_limit,omitempty"`
	Offset     *int                      `json:"_offset,omitempty"`
	Detailed   *bool                     `json:"detailed,omitempty"`
	PortIdList *ListParametersPortIdList `json:"port_id_list,omitempty"`
	VpcId      string                    `json:"vpc_id"`
}

type ListParametersPortIdList []string

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	PortsSimplified ListResultPortsSimplified `json:"ports_simplified"`
	Ports           ListResultPorts           `json:"ports"`
}

type ListResultPortsSimplifiedItem struct {
	Id        *string                                 `json:"id,omitempty"`
	IpAddress *ListResultPortsSimplifiedItemIpAddress `json:"ip_address,omitempty"`
}

type ListResultPortsSimplifiedItemIpAddressItem struct {
	IpAddress string `json:"ip_address"`
	SubnetId  string `json:"subnet_id"`
}

type ListResultPortsSimplifiedItemIpAddress []ListResultPortsSimplifiedItemIpAddressItem

type ListResultPortsSimplified []ListResultPortsSimplifiedItem

type ListResultPortsItem struct {
	CreatedAt             *string                            `json:"created_at,omitempty"`
	Description           *string                            `json:"description,omitempty"`
	Id                    *string                            `json:"id,omitempty"`
	IpAddress             *ListResultPortsItemIpAddress      `json:"ip_address,omitempty"`
	IsAdminStateUp        *bool                              `json:"is_admin_state_up,omitempty"`
	IsPortSecurityEnabled *bool                              `json:"is_port_security_enabled,omitempty"`
	Name                  *string                            `json:"name,omitempty"`
	PublicIp              *ListResultPortsItemPublicIp       `json:"public_ip,omitempty"`
	SecurityGroups        *ListResultPortsItemSecurityGroups `json:"security_groups,omitempty"`
	Updated               *string                            `json:"updated,omitempty"`
	VpcId                 *string                            `json:"vpc_id,omitempty"`
}

type ListResultPortsItemIpAddressItem struct {
	IpAddress string `json:"ip_address"`
	SubnetId  string `json:"subnet_id"`
}

type ListResultPortsItemIpAddress []ListResultPortsItemIpAddressItem

type ListResultPortsItemPublicIpItem struct {
	PublicIp   *string `json:"public_ip,omitempty"`
	PublicIpId *string `json:"public_ip_id,omitempty"`
}

type ListResultPortsItemPublicIp []ListResultPortsItemPublicIpItem

type ListResultPortsItemSecurityGroups []string

type ListResultPorts []ListResultPortsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/vpcs/ports/list"), s.client, s.ctx)
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
