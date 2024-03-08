/*
Executor: list

# Summary

# List Security Groups by Tenant

# Description

Returns a list of security groups for a provided tenant_id

Version: 1.109.0

import "magalu.cloud/lib/products/network/security_groups"
*/
package securityGroups

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	SecurityGroups ListResultSecurityGroups `json:"security_groups"`
}

type ListResultSecurityGroupsItem struct {
	CreatedAt   *string `json:"created_at,omitempty"`
	Description *string `json:"description,omitempty"`
	Error       *string `json:"error,omitempty"`
	Id          *string `json:"id,omitempty"`
	IsDefault   *bool   `json:"is_default,omitempty"`
	Name        *string `json:"name,omitempty"`
	ProjectType *string `json:"project_type,omitempty"`
	Status      string  `json:"status"`
	TenantId    *string `json:"tenant_id,omitempty"`
	Updated     *string `json:"updated,omitempty"`
	VpcId       *string `json:"vpc_id,omitempty"`
}

type ListResultSecurityGroups []ListResultSecurityGroupsItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/network/security_groups/list"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

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
