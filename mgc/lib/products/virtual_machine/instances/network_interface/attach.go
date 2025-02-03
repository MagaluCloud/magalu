/*
Executor: attach

# Summary

Attach network interface to an instance.

# Description

Attach network interface to an instance for a default project.

Version: v1

import "github.com/MagaluCloud/magalu/mgc/lib/products/virtual_machine/instances/network_interface"
*/
package networkInterface

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type AttachParameters struct {
	Instance AttachParametersInstance `json:"instance"`
	Network  AttachParametersNetwork  `json:"network"`
}

// any of: AttachParametersInstance
type AttachParametersInstance struct {
	Id             string                                  `json:"id"`
	Name           *string                                 `json:"name,omitempty"`
	SecurityGroups *AttachParametersInstanceSecurityGroups `json:"security_groups,omitempty"`
}

type AttachParametersInstanceSecurityGroupsItem struct {
	Id string `json:"id"`
}

type AttachParametersInstanceSecurityGroups []AttachParametersInstanceSecurityGroupsItem

type AttachParametersNetwork struct {
	Interface AttachParametersNetworkInterface `json:"interface"`
}

type AttachParametersNetworkInterface struct {
	Id string `json:"id"`
}

type AttachConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Attach(
	parameters AttachParameters,
	configs AttachConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Attach", mgcCore.RefPath("/virtual-machine/instances/network-interface/attach"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[AttachParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[AttachConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) AttachContext(
	ctx context.Context,
	parameters AttachParameters,
	configs AttachConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Attach", mgcCore.RefPath("/virtual-machine/instances/network-interface/attach"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[AttachParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[AttachConfigs](configs); err != nil {
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
