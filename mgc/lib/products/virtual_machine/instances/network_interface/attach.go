/*
Executor: attach

# Summary

Attach network interface to an instance.

# Description

Attach network interface to an instance for a default project.

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances/network_interface"
*/
package networkInterface

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
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

// TODO: links
// TODO: related
