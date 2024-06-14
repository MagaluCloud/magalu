/*
Executor: update

# Summary

Database instance update.

# Description

Updates a database instance.

Version: 1.21.1

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UpdateParameters struct {
	BackupRetentionDays *int    `json:"backup_retention_days,omitempty"`
	BackupStartAt       *string `json:"backup_start_at,omitempty"`
	Exchange            *string `json:"exchange,omitempty"`
	InstanceId          string  `json:"instance_id"`
	Status              *string `json:"status,omitempty"`
}

type UpdateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type UpdateResult struct {
	Addresses           UpdateResultAddresses  `json:"addresses"`
	BackupRetentionDays int                    `json:"backup_retention_days"`
	BackupStartAt       string                 `json:"backup_start_at"`
	CreatedAt           string                 `json:"created_at"`
	DatastoreId         string                 `json:"datastore_id"`
	EngineId            string                 `json:"engine_id"`
	FinishedAt          *string                `json:"finished_at,omitempty"`
	FlavorId            string                 `json:"flavor_id"`
	Generation          string                 `json:"generation"`
	Id                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Parameters          UpdateResultParameters `json:"parameters"`
	Replicas            *UpdateResultReplicas  `json:"replicas,omitempty"`
	StartedAt           *string                `json:"started_at,omitempty"`
	Status              string                 `json:"status"`
	UpdatedAt           *string                `json:"updated_at,omitempty"`
	Volume              UpdateResultVolume     `json:"volume"`
}

type UpdateResultAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type UpdateResultAddresses []UpdateResultAddressesItem

type UpdateResultParametersItem struct {
	Name  string                          `json:"name"`
	Value UpdateResultParametersItemValue `json:"value"`
}

// any of: *float64, *int, *bool, *string
type UpdateResultParametersItemValue any

type UpdateResultParameters []UpdateResultParametersItem

type UpdateResultReplicasItem struct {
	Addresses   UpdateResultReplicasItemAddresses  `json:"addresses"`
	CreatedAt   string                             `json:"created_at"`
	DatastoreId string                             `json:"datastore_id"`
	EngineId    string                             `json:"engine_id"`
	FinishedAt  *string                            `json:"finished_at,omitempty"`
	FlavorId    string                             `json:"flavor_id"`
	Generation  string                             `json:"generation"`
	Id          string                             `json:"id"`
	Name        string                             `json:"name"`
	Parameters  UpdateResultReplicasItemParameters `json:"parameters"`
	SourceId    string                             `json:"source_id"`
	StartedAt   *string                            `json:"started_at,omitempty"`
	Status      string                             `json:"status"`
	UpdatedAt   *string                            `json:"updated_at,omitempty"`
	Volume      UpdateResultReplicasItemVolume     `json:"volume"`
}

type UpdateResultReplicasItemAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type UpdateResultReplicasItemAddresses []UpdateResultReplicasItemAddressesItem

type UpdateResultReplicasItemParametersItem struct {
	Name  string                          `json:"name"`
	Value UpdateResultParametersItemValue `json:"value"`
}

type UpdateResultReplicasItemParameters []UpdateResultReplicasItemParametersItem

type UpdateResultReplicasItemVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

type UpdateResultReplicas []UpdateResultReplicasItem

type UpdateResultVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

func (s *service) Update(
	parameters UpdateParameters,
	configs UpdateConfigs,
) (
	result UpdateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/dbaas/instances/update"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UpdateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UpdateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[UpdateResult](r)
}

// TODO: links
// TODO: related
