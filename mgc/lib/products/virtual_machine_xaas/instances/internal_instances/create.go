/*
Executor: create

# Summary

# Instance Finish Create V1

# Description

This route is for internaly to update the information
when finished to request a instance creation on URP.

After requested successfully this route is called to save
the network information and urp instance ID.

### Note
This route is used only for internal proposes.

Version: 1.230.0

import "magalu.cloud/lib/products/virtual_machine_xaas/instances/internal_instances"
*/
package internalInstances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	Error         *string                     `json:"error,omitempty"`
	Id            string                      `json:"id"`
	NetworkIds    *CreateParametersNetworkIds `json:"network_ids,omitempty"`
	ProjectType   *string                     `json:"project_type,omitempty"`
	Status        string                      `json:"status"`
	UrpInstanceId *string                     `json:"urp_instance_id,omitempty"`
}

type CreateParametersNetworkIdsItem struct {
	Id string `json:"id"`
}

type CreateParametersNetworkIds []CreateParametersNetworkIdsItem

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/virtual-machine-xaas/instances/internal-instances/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
