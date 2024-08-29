/*
Executor: delete

# Summary

Delete a backup.

# Description

Delete a Backup for the currently authenticated tenant.

#### Rules
- The Backup's status must be "completed".
- The Backup's state must be "available".

#### Notes
  - Utilize the **block-storage backups** list command to retrieve a list of
    all Backups and obtain the ID of the Backup you wish to delete.

Version: v1

import "magalu.cloud/lib/products/block_storage/backups"
*/
package backups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteParameters struct {
	Id string `json:"id"`
}

type DeleteConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Delete(
	parameters DeleteParameters,
	configs DeleteConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Delete", mgcCore.RefPath("/block-storage/backups/delete"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[map[string]interface{}](s.client.Sdk().Config().TempConfig()); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

func (s *service) DeleteConfirmPrompt(
	parameters DeleteParameters,
	configs DeleteConfigs,
) (message string) {
	e, err := mgcHelpers.ResolveExecutor("Delete", mgcCore.RefPath("/block-storage/backups/delete"), s.client)
	if err != nil {
		return
	}

	exec, ok := e.(mgcCore.ConfirmableExecutor)
	if !ok {
		// Not expected, but let's return an empty message
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[map[string]interface{}](s.client.Sdk().Config().TempConfig()); err != nil {
		return
	}

	return exec.ConfirmPrompt(p, c)
}

// TODO: links
// TODO: related
