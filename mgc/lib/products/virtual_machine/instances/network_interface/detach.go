/*
Executor: detach

# Summary

Detach a non primary network interface from an instance.

# Description

# Detach a non primary network interface from an instance for a default project

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances/network_interface"
*/
package networkInterface

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DetachParameters struct {
	Instance DetachParametersInstance `json:"instance"`
	Network  DetachParametersNetwork  `json:"network"`
}

// any of: DetachParametersInstance
type DetachParametersInstance struct {
	Id             string                                  `json:"id"`
	Name           *string                                 `json:"name,omitempty"`
	SecurityGroups *DetachParametersInstanceSecurityGroups `json:"security_groups,omitempty"`
}

type DetachParametersInstanceSecurityGroupsItem struct {
	Id string `json:"id"`
}

type DetachParametersInstanceSecurityGroups []DetachParametersInstanceSecurityGroupsItem

type DetachParametersNetwork struct {
	Interface DetachParametersNetworkInterface `json:"interface"`
}

type DetachParametersNetworkInterface struct {
	Id string `json:"id"`
}

type DetachConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Detach(
	parameters DetachParameters,
	configs DetachConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Detach", mgcCore.RefPath("/virtual-machine/instances/network-interface/detach"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DetachParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DetachConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) DetachContext(
	ctx context.Context,
	parameters DetachParameters,
	configs DetachConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Detach", mgcCore.RefPath("/virtual-machine/instances/network-interface/detach"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DetachParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DetachConfigs](configs); err != nil {
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

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
