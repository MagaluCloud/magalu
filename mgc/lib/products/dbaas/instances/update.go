/*
Executor: update

# Summary

Database instance update.

# Description

Updates a database instance.

Version: 1.17.2

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UpdateParameters struct {
	BackupRetentionDays int    `json:"backup_retention_days,omitempty"`
	BackupStartAt       string `json:"backup_start_at,omitempty"`
	Exchange            string `json:"exchange,omitempty"`
	InstanceId          string `json:"instance_id"`
	Status              string `json:"status,omitempty"`
}

type UpdateConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type UpdateResult struct {
	Addresses           UpdateResultAddresses `json:"addresses"`
	BackupRetentionDays int                   `json:"backup_retention_days"`
	BackupStartAt       string                `json:"backup_start_at"`
	CreatedAt           string                `json:"created_at"`
	DatastoreId         string                `json:"datastore_id"`
	EngineId            string                `json:"engine_id"`
	FinishedAt          string                `json:"finished_at,omitempty"`
	FlavorId            string                `json:"flavor_id"`
	Generation          string                `json:"generation"`
	Id                  string                `json:"id"`
	Name                string                `json:"name"`
	Replicas            UpdateResultReplicas  `json:"replicas,omitempty"`
	StartedAt           string                `json:"started_at,omitempty"`
	Status              string                `json:"status"`
	UpdatedAt           string                `json:"updated_at,omitempty"`
	Volume              UpdateResultVolume    `json:"volume"`
}

type UpdateResultAddressesItem struct {
	Access  string `json:"access"`
	Address string `json:"address,omitempty"`
	Type    string `json:"type,omitempty"`
}

type UpdateResultAddresses []UpdateResultAddressesItem

type UpdateResultReplicasItem struct {
	Addresses   UpdateResultReplicasItemAddresses `json:"addresses"`
	CreatedAt   string                            `json:"created_at"`
	DatastoreId string                            `json:"datastore_id"`
	EngineId    string                            `json:"engine_id"`
	FinishedAt  string                            `json:"finished_at,omitempty"`
	FlavorId    string                            `json:"flavor_id"`
	Generation  string                            `json:"generation"`
	Id          string                            `json:"id"`
	Name        string                            `json:"name"`
	SourceId    string                            `json:"source_id"`
	StartedAt   string                            `json:"started_at,omitempty"`
	Status      string                            `json:"status"`
	UpdatedAt   string                            `json:"updated_at,omitempty"`
	Volume      UpdateResultReplicasItemVolume    `json:"volume"`
}

type UpdateResultReplicasItemAddressesItem struct {
	Access  string `json:"access"`
	Address string `json:"address,omitempty"`
	Type    string `json:"type,omitempty"`
}

type UpdateResultReplicasItemAddresses []UpdateResultReplicasItemAddressesItem

type UpdateResultReplicasItemVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

type UpdateResultReplicas []UpdateResultReplicasItem

type UpdateResultVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

func Update(
	client *mgcClient.Client,
	ctx context.Context,
	parameters UpdateParameters,
	configs UpdateConfigs,
) (
	result UpdateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/dbaas/instances/update"), client, ctx)
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
