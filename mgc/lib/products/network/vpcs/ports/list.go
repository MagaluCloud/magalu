/*
Executor: list

# Summary

# List Ports

# Description

Returns a list of ports for a provided vpc_id and x-tenant-id. The list will be paginated, it means you can easily find what you need just setting the page number(_offset) and the quantity of items per page(_limit). The level of detail can also be set

Version: 1.141.3

import "magalu.cloud/lib/products/network/vpcs/ports"
*/
package ports

import (
	"context"

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

// any of: ListResult
type ListResult struct {
	Ports           *ListResultPorts          `json:"ports,omitempty"`
	PortsSimplified ListResultPortsSimplified `json:"ports_simplified"`
}

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
	Ethertype *string `json:"ethertype,omitempty"`
	IpAddress string  `json:"ip_address"`
	SubnetId  string  `json:"subnet_id"`
}

type ListResultPortsItemIpAddress []ListResultPortsItemIpAddressItem

type ListResultPortsItemPublicIpItem struct {
	PublicIp   *string `json:"public_ip,omitempty"`
	PublicIpId *string `json:"public_ip_id,omitempty"`
}

type ListResultPortsItemPublicIp []ListResultPortsItemPublicIpItem

type ListResultPortsItemSecurityGroups []string

type ListResultPorts []ListResultPortsItem

type ListResultPortsSimplifiedItem struct {
	Id        *string                                 `json:"id,omitempty"`
	IpAddress *ListResultPortsSimplifiedItemIpAddress `json:"ip_address,omitempty"`
}

type ListResultPortsSimplifiedItemIpAddressItem struct {
	Ethertype *string `json:"ethertype,omitempty"`
	IpAddress string  `json:"ip_address"`
	SubnetId  string  `json:"subnet_id"`
}

type ListResultPortsSimplifiedItemIpAddress []ListResultPortsSimplifiedItemIpAddressItem

type ListResultPortsSimplified []ListResultPortsSimplifiedItem

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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ListContext(
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/vpcs/ports/list"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
