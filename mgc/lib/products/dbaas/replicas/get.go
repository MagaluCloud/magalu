/*
Executor: get

# Summary

Replica Detail.

# Description

Get an instance replica detail.

Version: 1.34.1

import "github.com/MagaluCloud/magalu/mgc/lib/products/dbaas/replicas"
*/
package replicas

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type GetParameters struct {
	ReplicaId string `json:"replica_id"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	Addresses              GetResultAddresses  `json:"addresses"`
	CreatedAt              string              `json:"created_at"`
	DatastoreId            string              `json:"datastore_id"`
	EngineId               string              `json:"engine_id"`
	FinishedAt             *string             `json:"finished_at,omitempty"`
	FlavorId               string              `json:"flavor_id"`
	Generation             string              `json:"generation"`
	Id                     string              `json:"id"`
	InstanceTypeId         string              `json:"instance_type_id"`
	MaintenanceScheduledAt *string             `json:"maintenance_scheduled_at,omitempty"`
	Name                   string              `json:"name"`
	Parameters             GetResultParameters `json:"parameters"`
	SourceId               string              `json:"source_id"`
	StartedAt              *string             `json:"started_at,omitempty"`
	Status                 string              `json:"status"`
	UpdatedAt              *string             `json:"updated_at,omitempty"`
	Volume                 GetResultVolume     `json:"volume"`
}

type GetResultAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type GetResultAddresses []GetResultAddressesItem

type GetResultParametersItem struct {
	Name  string                       `json:"name"`
	Value GetResultParametersItemValue `json:"value"`
}

// any of: *float64, *int, *bool, *string
type GetResultParametersItemValue any

type GetResultParameters []GetResultParametersItem

type GetResultVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/dbaas/replicas/get"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[GetConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) GetContext(
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/dbaas/replicas/get"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[GetConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
