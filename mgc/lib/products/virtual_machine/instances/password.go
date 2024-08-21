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

// TODO: links
// TODO: related
