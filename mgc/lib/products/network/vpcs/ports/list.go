/*
Executor: list

# Summary

# List Ports

# Description

# List VPC ports

Version: 1.109.0

import "magalu.cloud/lib/products/network/vpcs/ports"
*/
package ports

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
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
	Ports ListResult0Ports `json:"ports"`
}

type ListResult0PortsItem struct {
	CreatedAt             *string                            `json:"created_at,omitempty"`
	Description           *string                            `json:"description,omitempty"`
	Id                    *string                            `json:"id,omitempty"`
	IpAddress             ListResult0PortsItemIpAddress      `json:"ip_address,omitempty"`
	IsAdminStateUp        *bool                              `json:"is_admin_state_up,omitempty"`
	IsPortSecurityEnabled *bool                              `json:"is_port_security_enabled,omitempty"`
	Name                  *string                            `json:"name,omitempty"`
	PublicIp              *ListResult0PortsItemPublicIp      `json:"public_ip,omitempty"`
	SecurityGroups        ListResult0PortsItemSecurityGroups `json:"security_groups,omitempty"`
	Updated               *string                            `json:"updated,omitempty"`
	VpcId                 *string                            `json:"vpc_id,omitempty"`
}

type ListResult0PortsItemIpAddressItem struct {
	IpAddress string `json:"ip_address"`
	SubnetId  string `json:"subnet_id"`
}

type ListResult0PortsItemIpAddress []ListResult0PortsItemIpAddressItem

type ListResult0PortsItemPublicIpItem struct {
	PublicIp   *string `json:"public_ip,omitempty"`
	PublicIpId *string `json:"public_ip_id,omitempty"`
}

type ListResult0PortsItemPublicIp []ListResult0PortsItemPublicIpItem

type ListResult0PortsItemSecurityGroups []string

type ListResult0Ports []ListResult0PortsItem

type ListResult1 struct {
	PortsSimplified ListResult1PortsSimplified `json:"ports_simplified"`
}

type ListResult1PortsSimplifiedItem struct {
	Id        *string                                 `json:"id,omitempty"`
	IpAddress ListResult1PortsSimplifiedItemIpAddress `json:"ip_address,omitempty"`
}

type ListResult1PortsSimplifiedItemIpAddress []ListResult0PortsItemIpAddressItem

type ListResult1PortsSimplified []ListResult1PortsSimplifiedItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/vpcs/ports/list"), client, ctx)
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
