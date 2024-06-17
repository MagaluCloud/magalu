/*
Executor: stop

# Summary

Stops a database instance.

# Description

Stops a database instance.

Version: 1.21.2

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type StopParameters struct {
	Exchange   *string `json:"exchange,omitempty"`
	InstanceId string  `json:"instance_id"`
}

type StopConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type StopResult struct {
	Addresses           StopResultAddresses  `json:"addresses"`
	BackupRetentionDays int                  `json:"backup_retention_days"`
	BackupStartAt       string               `json:"backup_start_at"`
	CreatedAt           string               `json:"created_at"`
	DatastoreId         string               `json:"datastore_id"`
	EngineId            string               `json:"engine_id"`
	FinishedAt          *string              `json:"finished_at,omitempty"`
	FlavorId            string               `json:"flavor_id"`
	Generation          string               `json:"generation"`
	Id                  string               `json:"id"`
	Name                string               `json:"name"`
	Parameters          StopResultParameters `json:"parameters"`
	Replicas            *StopResultReplicas  `json:"replicas,omitempty"`
	StartedAt           *string              `json:"started_at,omitempty"`
	Status              string               `json:"status"`
	UpdatedAt           *string              `json:"updated_at,omitempty"`
	Volume              StopResultVolume     `json:"volume"`
}

type StopResultAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type StopResultAddresses []StopResultAddressesItem

type StopResultParametersItem struct {
	Name  string                        `json:"name"`
	Value StopResultParametersItemValue `json:"value"`
}

// any of: *float64, *int, *bool, *string
type StopResultParametersItemValue any

type StopResultParameters []StopResultParametersItem

type StopResultReplicasItem struct {
	Addresses   StopResultReplicasItemAddresses  `json:"addresses"`
	CreatedAt   string                           `json:"created_at"`
	DatastoreId string                           `json:"datastore_id"`
	EngineId    string                           `json:"engine_id"`
	FinishedAt  *string                          `json:"finished_at,omitempty"`
	FlavorId    string                           `json:"flavor_id"`
	Generation  string                           `json:"generation"`
	Id          string                           `json:"id"`
	Name        string                           `json:"name"`
	Parameters  StopResultReplicasItemParameters `json:"parameters"`
	SourceId    string                           `json:"source_id"`
	StartedAt   *string                          `json:"started_at,omitempty"`
	Status      string                           `json:"status"`
	UpdatedAt   *string                          `json:"updated_at,omitempty"`
	Volume      StopResultReplicasItemVolume     `json:"volume"`
}

type StopResultReplicasItemAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type StopResultReplicasItemAddresses []StopResultReplicasItemAddressesItem

type StopResultReplicasItemParametersItem struct {
	Name  string                        `json:"name"`
	Value StopResultParametersItemValue `json:"value"`
}

type StopResultReplicasItemParameters []StopResultReplicasItemParametersItem

type StopResultReplicasItemVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

type StopResultReplicas []StopResultReplicasItem

type StopResultVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

func (s *service) Stop(
	parameters StopParameters,
	configs StopConfigs,
) (
	result StopResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Stop", mgcCore.RefPath("/dbaas/instances/stop"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[StopParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[StopConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[StopResult](r)
}

// TODO: links
// TODO: related
