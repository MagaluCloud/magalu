/*
Executor: resize

# Summary

Resizes a database instance.

# Description

Resizes a database instance.

Version: 1.34.1

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ResizeParameters struct {
	FlavorId       *string                 `json:"flavor_id,omitempty"`
	InstanceId     string                  `json:"instance_id"`
	InstanceTypeId *string                 `json:"instance_type_id,omitempty"`
	Volume         *ResizeParametersVolume `json:"volume,omitempty"`
}

type ResizeParametersVolume struct {
	Size int     `json:"size"`
	Type *string `json:"type,omitempty"`
}

type ResizeConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ResizeResult struct {
	Addresses              ResizeResultAddresses          `json:"addresses"`
	BackupRetentionDays    int                            `json:"backup_retention_days"`
	BackupStartAt          string                         `json:"backup_start_at"`
	CreatedAt              string                         `json:"created_at"`
	DatastoreId            string                         `json:"datastore_id"`
	EngineId               string                         `json:"engine_id"`
	FinishedAt             *string                        `json:"finished_at,omitempty"`
	FlavorId               string                         `json:"flavor_id"`
	Generation             string                         `json:"generation"`
	Id                     string                         `json:"id"`
	InstanceTypeId         string                         `json:"instance_type_id"`
	MaintenanceScheduledAt *string                        `json:"maintenance_scheduled_at,omitempty"`
	Name                   string                         `json:"name"`
	Parameters             ResizeResultParameters         `json:"parameters"`
	Replicas               *ResizeResultReplicas          `json:"replicas,omitempty"`
	StartedAt              *string                        `json:"started_at,omitempty"`
	Status                 string                         `json:"status"`
	UpdatedAt              *string                        `json:"updated_at,omitempty"`
	Volume                 ResizeResultReplicasItemVolume `json:"volume"`
}

type ResizeResultAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type ResizeResultAddresses []ResizeResultAddressesItem

type ResizeResultParametersItem struct {
	Name  string                          `json:"name"`
	Value ResizeResultParametersItemValue `json:"value"`
}

// any of: *float64, *int, *bool, *string
type ResizeResultParametersItemValue any

type ResizeResultParameters []ResizeResultParametersItem

type ResizeResultReplicasItem struct {
	Addresses              ResizeResultReplicasItemAddresses  `json:"addresses"`
	CreatedAt              string                             `json:"created_at"`
	DatastoreId            string                             `json:"datastore_id"`
	EngineId               string                             `json:"engine_id"`
	FinishedAt             *string                            `json:"finished_at,omitempty"`
	FlavorId               string                             `json:"flavor_id"`
	Generation             string                             `json:"generation"`
	Id                     string                             `json:"id"`
	InstanceTypeId         string                             `json:"instance_type_id"`
	MaintenanceScheduledAt *string                            `json:"maintenance_scheduled_at,omitempty"`
	Name                   string                             `json:"name"`
	Parameters             ResizeResultReplicasItemParameters `json:"parameters"`
	SourceId               string                             `json:"source_id"`
	StartedAt              *string                            `json:"started_at,omitempty"`
	Status                 string                             `json:"status"`
	UpdatedAt              *string                            `json:"updated_at,omitempty"`
	Volume                 ResizeResultReplicasItemVolume     `json:"volume"`
}

type ResizeResultReplicasItemAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type ResizeResultReplicasItemAddresses []ResizeResultReplicasItemAddressesItem

type ResizeResultReplicasItemParametersItem struct {
	Name  string                          `json:"name"`
	Value ResizeResultParametersItemValue `json:"value"`
}

type ResizeResultReplicasItemParameters []ResizeResultReplicasItemParametersItem

type ResizeResultReplicasItemVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

type ResizeResultReplicas []ResizeResultReplicasItem

func (s *service) Resize(
	parameters ResizeParameters,
	configs ResizeConfigs,
) (
	result ResizeResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Resize", mgcCore.RefPath("/dbaas/instances/resize"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ResizeParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ResizeConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ResizeResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ResizeContext(
	ctx context.Context,
	parameters ResizeParameters,
	configs ResizeConfigs,
) (
	result ResizeResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Resize", mgcCore.RefPath("/dbaas/instances/resize"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ResizeParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ResizeConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[ResizeResult](r)
}

// TODO: links
// TODO: related
