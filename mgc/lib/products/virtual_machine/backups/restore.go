/*
Executor: restore

# Summary

# Restore a backup to a virtual machine

# Description

Restore a backup of a Virtual Machine with the current tenant which is logged in. </br>
A Backup is ready for restore when it's in completed status.

#### Notes
- You can verify the status of backup using the backup list command.
- Use machine-types list to see all machine types available.

#### Rules
- To restore a backup you have to inform the new virtual machine information.
- You can choose a machine-type that has a disk equal or larger
than the minimum disk of the backup.

Version: v1

import "magalu.cloud/lib/products/virtual_machine/backups"
*/
package backups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type RestoreParameters struct {
	AvailabilityZone *string                      `json:"availability_zone,omitempty"`
	Id               string                       `json:"id"`
	MachineType      RestoreParametersMachineType `json:"machine_type"`
	Name             string                       `json:"name"`
	Network          RestoreParametersNetwork     `json:"network"`
	SshKeyName       string                       `json:"ssh_key_name"`
	UserData         *string                      `json:"user_data,omitempty"`
}

// any of: RestoreParametersMachineType
type RestoreParametersMachineType struct {
	Id   string  `json:"id"`
	Name *string `json:"name,omitempty"`
}

type RestoreParametersNetwork struct {
	AssociatePublicIp *bool                              `json:"associate_public_ip,omitempty"`
	Interface         *RestoreParametersNetworkInterface `json:"interface,omitempty"`
	Vpc               *RestoreParametersNetworkVpc       `json:"vpc,omitempty"`
}

// any of: RestoreParametersNetworkInterface
type RestoreParametersNetworkInterface struct {
	Id             string                                           `json:"id"`
	Name           *string                                          `json:"name,omitempty"`
	SecurityGroups *RestoreParametersNetworkInterfaceSecurityGroups `json:"security_groups,omitempty"`
}

type RestoreParametersNetworkInterfaceSecurityGroupsItem struct {
	Id string `json:"id"`
}

type RestoreParametersNetworkInterfaceSecurityGroups []RestoreParametersNetworkInterfaceSecurityGroupsItem

// any of: RestoreParametersNetworkVpc
type RestoreParametersNetworkVpc struct {
	Id             string                                           `json:"id"`
	Name           *string                                          `json:"name,omitempty"`
	SecurityGroups *RestoreParametersNetworkInterfaceSecurityGroups `json:"security_groups,omitempty"`
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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Restore", mgcCore.RefPath("/virtual-machine/backups/restore"), s.client, s.ctx)
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
