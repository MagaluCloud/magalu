/*
Executor: create-id

# Summary

# Restore Snapshot

# Description

# Restore a snapshot

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/snapshots/v0_snapshots"
*/
package v0Snapshots

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateIdParameters struct {
	AllocateFip       bool                                 `json:"allocate_fip,omitempty"`
	AvailabilityZone  *string                              `json:"availability_zone,omitempty"`
	Id                string                               `json:"id"`
	KeyName           string                               `json:"key_name"`
	Name              string                               `json:"name"`
	NetworkInterfaces *CreateIdParametersNetworkInterfaces `json:"network_interfaces,omitempty"`
	Type              string                               `json:"type"`
	UserData          *string                              `json:"user_data,omitempty"`
}

type CreateIdParametersNetworkInterfacesItem struct {
	Id string `json:"id"`
}

type CreateIdParametersNetworkInterfaces []CreateIdParametersNetworkInterfacesItem

type CreateIdConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type CreateIdResult struct {
	Id string `json:"id"`
}

func CreateId(
	client *mgcClient.Client,
	ctx context.Context,
	parameters CreateIdParameters,
	configs CreateIdConfigs,
) (
	result CreateIdResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("CreateId", mgcCore.RefPath("/virtual-machine/snapshots/v0-snapshots/create-id"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateIdParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateIdConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateIdResult](r)
}

// TODO: links
// TODO: related
