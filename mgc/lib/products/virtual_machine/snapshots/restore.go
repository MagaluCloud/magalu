/*
Executor: restore

# Summary

Restore a snapshot to an instance.

# Description

Restore a snapshot of an instance with the current tenant which is logged in. </br>

#### Notes
- You can check the snapshot state using snapshot list command.
- Use "machine-types list" to see all machine types available.

#### Rules
- A Snapshot is ready to restore when it's in available state.
- To restore a snapshot you have to inform the new instance settings.
- You must choose a machine-type that has a disk equal or larger
than the original instance.

Version: v1

import "magalu.cloud/lib/products/virtual_machine/snapshots"
*/
package snapshots

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RestoreParameters struct {
	AvailabilityZone *string                      `json:"availability_zone,omitempty"`
	Id               string                       `json:"id"`
	MachineType      RestoreParametersMachineType `json:"machine_type"`
	Name             string                       `json:"name"`
	Network          *RestoreParametersNetwork    `json:"network,omitempty"`
	SshKeyName       *string                      `json:"ssh_key_name,omitempty"`
	UserData         *string                      `json:"user_data,omitempty"`
}

// any of: RestoreParametersMachineType
type RestoreParametersMachineType struct {
	Id             string                                      `json:"id"`
	Name           *string                                     `json:"name,omitempty"`
	SecurityGroups *RestoreParametersMachineTypeSecurityGroups `json:"security_groups,omitempty"`
}

type RestoreParametersMachineTypeSecurityGroupsItem struct {
	Id string `json:"id"`
}

type RestoreParametersMachineTypeSecurityGroups []RestoreParametersMachineTypeSecurityGroupsItem

type RestoreParametersNetwork struct {
	AssociatePublicIp *bool                              `json:"associate_public_ip,omitempty"`
	Interface         *RestoreParametersNetworkInterface `json:"interface,omitempty"`
	Vpc               *RestoreParametersNetworkVpc       `json:"vpc,omitempty"`
}

// any of: RestoreParametersNetworkInterface
type RestoreParametersNetworkInterface struct {
	Id             string                                      `json:"id"`
	Name           *string                                     `json:"name,omitempty"`
	SecurityGroups *RestoreParametersMachineTypeSecurityGroups `json:"security_groups,omitempty"`
}

// any of: RestoreParametersNetworkVpc
type RestoreParametersNetworkVpc struct {
	Id             string                                      `json:"id"`
	Name           *string                                     `json:"name,omitempty"`
	SecurityGroups *RestoreParametersMachineTypeSecurityGroups `json:"security_groups,omitempty"`
}

type RestoreConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type RestoreResult struct {
	Id string `json:"id"`
}

func (s *service) Restore(
	parameters RestoreParameters,
	configs RestoreConfigs,
) (
	result RestoreResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Restore", mgcCore.RefPath("/virtual-machine/snapshots/restore"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RestoreParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[RestoreConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[RestoreResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) RestoreContext(
	ctx context.Context,
	parameters RestoreParameters,
	configs RestoreConfigs,
) (
	result RestoreResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Restore", mgcCore.RefPath("/virtual-machine/snapshots/restore"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[RestoreParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[RestoreConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[RestoreResult](r)
}

// TODO: links
// TODO: related
