/*
Executor: start

# Summary

Replica Start.

# Description

Start an instance replica.

Version: 1.34.1

import "magalu.cloud/lib/products/dbaas/replicas"
*/
package replicas

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type StartParameters struct {
	ReplicaId string `json:"replica_id"`
}

type StartConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type StartResult struct {
	Addresses              StartResultAddresses  `json:"addresses"`
	CreatedAt              string                `json:"created_at"`
	DatastoreId            string                `json:"datastore_id"`
	EngineId               string                `json:"engine_id"`
	FinishedAt             *string               `json:"finished_at,omitempty"`
	FlavorId               string                `json:"flavor_id"`
	Generation             string                `json:"generation"`
	Id                     string                `json:"id"`
	InstanceTypeId         string                `json:"instance_type_id"`
	MaintenanceScheduledAt *string               `json:"maintenance_scheduled_at,omitempty"`
	Name                   string                `json:"name"`
	Parameters             StartResultParameters `json:"parameters"`
	SourceId               string                `json:"source_id"`
	StartedAt              *string               `json:"started_at,omitempty"`
	Status                 string                `json:"status"`
	UpdatedAt              *string               `json:"updated_at,omitempty"`
	Volume                 StartResultVolume     `json:"volume"`
}

type StartResultAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type StartResultAddresses []StartResultAddressesItem

type StartResultParametersItem struct {
	Name  string                         `json:"name"`
	Value StartResultParametersItemValue `json:"value"`
}

// any of: *float64, *int, *bool, *string
type StartResultParametersItemValue any

type StartResultParameters []StartResultParametersItem

type StartResultVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

func (s *service) Start(
	parameters StartParameters,
	configs StartConfigs,
) (
	result StartResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Start", mgcCore.RefPath("/dbaas/replicas/start"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[StartParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[StartConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[StartResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) StartContext(
	ctx context.Context,
	parameters StartParameters,
	configs StartConfigs,
) (
	result StartResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Start", mgcCore.RefPath("/dbaas/replicas/start"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[StartParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[StartConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[StartResult](r)
}

// TODO: links
// TODO: related
