/*
Executor: password

# Summary

# Reset password

# Description

Reset container registry user's password.

Version: 0.1.0

import "magalu.cloud/lib/products/container_registry/credentials"
*/
package credentials

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type PasswordConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// User's credentials for authentication to the container registry.
type PasswordResult struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (s *service) Password(
	configs PasswordConfigs,
) (
	result PasswordResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Password", mgcCore.RefPath("/container-registry/credentials/password"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

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
