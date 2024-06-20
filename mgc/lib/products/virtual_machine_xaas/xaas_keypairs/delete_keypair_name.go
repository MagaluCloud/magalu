/*
Executor: delete-keypair-name

# Summary

# Delete Keypair Xaas

# Description

# Delete a keypair

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/xaas_keypairs"
*/
package xaasKeypairs

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteKeypairNameParameters struct {
	KeypairName string `json:"keypair_name"`
	ProjectType string `json:"project_type"`
}

type DeleteKeypairNameConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) DeleteKeypairName(
	parameters DeleteKeypairNameParameters,
	configs DeleteKeypairNameConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("DeleteKeypairName", mgcCore.RefPath("/virtual-machine-xaas/xaas keypairs/delete-keypair-name"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteKeypairNameParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeleteKeypairNameConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

func (s *service) DeleteKeypairNameConfirmPrompt(
	parameters DeleteKeypairNameParameters,
	configs DeleteKeypairNameConfigs,
) (message string) {
	e, err := mgcHelpers.ResolveExecutor("DeleteKeypairName", mgcCore.RefPath("/virtual-machine-xaas/xaas keypairs/delete-keypair-name"), s.client)
	if err != nil {
		return
	}

	exec, ok := e.(mgcCore.ConfirmableExecutor)
	if !ok {
		// Not expected, but let's return an empty message
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteKeypairNameParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[DeleteKeypairNameConfigs](configs); err != nil {
		return
	}

	return exec.ConfirmPrompt(p, c)
}

// TODO: links
// TODO: related
