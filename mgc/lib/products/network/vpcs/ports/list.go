/*
Executor: list

# Summary

# List Ports

# Description

# List VPC ports

Version: 1.119.0

import "magalu.cloud/lib/products/network/vpcs/ports"
*/
package ports

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit      int                      `json:"_limit,omitempty"`
	Offset     int                      `json:"_offset,omitempty"`
	Detailed   bool                     `json:"detailed,omitempty"`
	PortIdList ListParametersPortIdList `json:"port_id_list,omitempty"`
	VpcId      string                   `json:"vpc_id"`
}

type ListParametersPortIdList []string

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

// any of: ListResult0, ListResult1
type ListResult struct {
	ListResult0 `json:",squash"` // nolint
	ListResult1 `json:",squash"` // nolint
}

type ListResult0 struct {
	PortsSimplified ListResult0PortsSimplified `json:"ports_simplified"`
}

type ListResult0PortsSimplifiedItem struct {
	Id        *string                                 `json:"id,omitempty"`
	IpAddress ListResult0PortsSimplifiedItemIpAddress `json:"ip_address,omitempty"`
}

type ListResult0PortsSimplifiedItemIpAddressItem struct {
	IpAddress string `json:"ip_address"`
	SubnetId  string `json:"subnet_id"`
}

type ListResult0PortsSimplifiedItemIpAddress []ListResult0PortsSimplifiedItemIpAddressItem

type ListResult0PortsSimplified []ListResult0PortsSimplifiedItem

type ListResult1 struct {
	Ports ListResult1Ports `json:"ports"`
}

type ListResult1PortsItem struct {
	CreatedAt             *string                            `json:"created_at,omitempty"`
	Description           *string                            `json:"description,omitempty"`
	Id                    *string                            `json:"id,omitempty"`
	IpAddress             ListResult1PortsItemIpAddress      `json:"ip_address,omitempty"`
	IsAdminStateUp        *bool                              `json:"is_admin_state_up,omitempty"`
	IsPortSecurityEnabled *bool                              `json:"is_port_security_enabled,omitempty"`
	Name                  *string                            `json:"name,omitempty"`
	PublicIp              *ListResult1PortsItemPublicIp      `json:"public_ip,omitempty"`
	SecurityGroups        ListResult1PortsItemSecurityGroups `json:"security_groups,omitempty"`
	Updated               *string                            `json:"updated,omitempty"`
	VpcId                 *string                            `json:"vpc_id,omitempty"`
}

type ListResult1PortsItemIpAddress []ListResult0PortsSimplifiedItemIpAddressItem

type ListResult1PortsItemPublicIpItem struct {
	PublicIp   *string `json:"public_ip,omitempty"`
	PublicIpId *string `json:"public_ip_id,omitempty"`
}

type ListResult1PortsItemPublicIp []ListResult1PortsItemPublicIpItem

type ListResult1PortsItemSecurityGroups []string

type ListResult1Ports []ListResult1PortsItem

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
