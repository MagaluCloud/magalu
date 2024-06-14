/*
Executor: start

# Summary

Starts a database instance.

# Description

Starts a database instance.

Version: 1.21.1

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type StartParameters struct {
	Exchange   *string `json:"exchange,omitempty"`
	InstanceId string  `json:"instance_id"`
}

type StartConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type StartResult struct {
	Addresses           StartResultAddresses  `json:"addresses"`
	BackupRetentionDays int                   `json:"backup_retention_days"`
	BackupStartAt       string                `json:"backup_start_at"`
	CreatedAt           string                `json:"created_at"`
	DatastoreId         string                `json:"datastore_id"`
	EngineId            string                `json:"engine_id"`
	FinishedAt          *string               `json:"finished_at,omitempty"`
	FlavorId            string                `json:"flavor_id"`
	Generation          string                `json:"generation"`
	Id                  string                `json:"id"`
	Name                string                `json:"name"`
	Parameters          StartResultParameters `json:"parameters"`
	Replicas            *StartResultReplicas  `json:"replicas,omitempty"`
	StartedAt           *string               `json:"started_at,omitempty"`
	Status              string                `json:"status"`
	UpdatedAt           *string               `json:"updated_at,omitempty"`
	Volume              StartResultVolume     `json:"volume"`
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

type StartResultReplicasItem struct {
	Addresses   StartResultReplicasItemAddresses  `json:"addresses"`
	CreatedAt   string                            `json:"created_at"`
	DatastoreId string                            `json:"datastore_id"`
	EngineId    string                            `json:"engine_id"`
	FinishedAt  *string                           `json:"finished_at,omitempty"`
	FlavorId    string                            `json:"flavor_id"`
	Generation  string                            `json:"generation"`
	Id          string                            `json:"id"`
	Name        string                            `json:"name"`
	Parameters  StartResultReplicasItemParameters `json:"parameters"`
	SourceId    string                            `json:"source_id"`
	StartedAt   *string                           `json:"started_at,omitempty"`
	Status      string                            `json:"status"`
	UpdatedAt   *string                           `json:"updated_at,omitempty"`
	Volume      StartResultReplicasItemVolume     `json:"volume"`
}

type StartResultReplicasItemAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type StartResultReplicasItemAddresses []StartResultReplicasItemAddressesItem

type StartResultReplicasItemParametersItem struct {
	Name  string                         `json:"name"`
	Value StartResultParametersItemValue `json:"value"`
}

type StartResultReplicasItemParameters []StartResultReplicasItemParametersItem

type StartResultReplicasItemVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

type StartResultReplicas []StartResultReplicasItem

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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Start", mgcCore.RefPath("/dbaas/instances/start"), s.client, s.ctx)
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

// TODO: links
// TODO: related
