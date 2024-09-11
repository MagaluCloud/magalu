/*
Executor: password

# Summary

# Retrieve the first windows admin password

# Description

Retrieves the Windows Administrator password for the informed instance.

	The password is accessible only once and has a built-in
	expiration date to enhance security.

Version: v1

import "magalu.cloud/lib/products/virtual_machine/instances"
*/
package instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type PasswordParameters struct {
	Id string `json:"id"`
}

type PasswordConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type PasswordResult struct {
	Instance PasswordResultInstance `json:"instance"`
}

type PasswordResultInstance struct {
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
	Password  string `json:"password"`
}

func (s *service) Password(
	parameters PasswordParameters,
	configs PasswordConfigs,
) (
	result PasswordResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Password", mgcCore.RefPath("/virtual-machine/instances/password"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[PasswordParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[PasswordConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[PasswordResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) PasswordContext(
	ctx context.Context,
	parameters PasswordParameters,
	configs PasswordConfigs,
) (
	result PasswordResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Password", mgcCore.RefPath("/virtual-machine/instances/password"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[PasswordParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[PasswordConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[PasswordResult](r)
}

// TODO: links
// TODO: related
