/*
Executor: stop

# Summary

Stops a database instance.

# Description

Stops a database instance.

Version: 1.15.3

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type StopParameters struct {
	Exchange   string `json:"exchange,omitempty"`
	InstanceId string `json:"instance_id"`
}

type StopConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type StopResult struct {
	Addresses           StopResultAddresses `json:"addresses"`
	BackupRetentionDays int                 `json:"backup_retention_days"`
	BackupStartAt       string              `json:"backup_start_at"`
	CreatedAt           string              `json:"created_at"`
	DatastoreId         string              `json:"datastore_id"`
	FinishedAt          string              `json:"finished_at,omitempty"`
	FlavorId            string              `json:"flavor_id"`
	Generation          string              `json:"generation"`
	Id                  string              `json:"id"`
	Name                string              `json:"name"`
	Replicas            StopResultReplicas  `json:"replicas,omitempty"`
	StartedAt           string              `json:"started_at,omitempty"`
	Status              string              `json:"status"`
	UpdatedAt           string              `json:"updated_at,omitempty"`
	Volume              StopResultVolume    `json:"volume"`
}

type StopResultAddressesItem struct {
	Access  string `json:"access"`
	Address string `json:"address,omitempty"`
	Type    string `json:"type,omitempty"`
}

type StopResultAddresses []StopResultAddressesItem

type StopResultReplicasItem struct {
	Addresses   StopResultReplicasItemAddresses `json:"addresses"`
	CreatedAt   string                          `json:"created_at"`
	DatastoreId string                          `json:"datastore_id"`
	FinishedAt  string                          `json:"finished_at,omitempty"`
	FlavorId    string                          `json:"flavor_id"`
	Generation  string                          `json:"generation"`
	Id          string                          `json:"id"`
	Name        string                          `json:"name"`
	SourceId    string                          `json:"source_id"`
	StartedAt   string                          `json:"started_at,omitempty"`
	Status      string                          `json:"status"`
	UpdatedAt   string                          `json:"updated_at,omitempty"`
	Volume      StopResultReplicasItemVolume    `json:"volume"`
}

type StopResultReplicasItemAddressesItem struct {
	Access  string `json:"access"`
	Address string `json:"address,omitempty"`
	Type    string `json:"type,omitempty"`
}

type StopResultReplicasItemAddresses []StopResultReplicasItemAddressesItem

type StopResultReplicasItemVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

type StopResultReplicas []StopResultReplicasItem

type StopResultVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

func Stop(
	client *mgcClient.Client,
	ctx context.Context,
	parameters StopParameters,
	configs StopConfigs,
) (
	result StopResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Stop", mgcCore.RefPath("/dbaas/instances/stop"), client, ctx)
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
