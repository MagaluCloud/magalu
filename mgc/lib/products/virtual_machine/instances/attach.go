/*
Executor: attach

# Summary

# Instance Attach Port

# Description

# Attach new port at Instance

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances"
*/
package instances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type AttachParameters struct {
	ForceAuthentication *bool   `json:"force-authentication,omitempty"`
	Id                  string  `json:"id"`
	ProjectType         *string `json:"project_type,omitempty"`
}

type AttachConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type AttachResult struct {
	Id string `json:"id"`
}

func (s *service) Attach(
	parameters AttachParameters,
	configs AttachConfigs,
) (
	result AttachResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Attach", mgcCore.RefPath("/virtual-machine/instances/attach"), s.client, s.ctx)
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

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[AttachResult](r)
}

// TODO: links
// TODO: related
