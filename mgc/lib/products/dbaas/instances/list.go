/*
Executor: list

# Summary

List all database instances.

# Description

Returns a list of database instances for a x-tenant-id.

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

type ListParameters struct {
	Expand   string `json:"_expand,omitempty"`
	Limit    int    `json:"_limit,omitempty"`
	Offset   int    `json:"_offset,omitempty"`
	Exchange string `json:"exchange,omitempty"`
	Status   string `json:"status,omitempty"`
}

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Results ListResultResults `json:"results"`
}

type ListResultResultsItem struct {
	Addresses           ListResultResultsItemAddresses `json:"addresses"`
	BackupRetentionDays int                            `json:"backup_retention_days"`
	BackupStartAt       string                         `json:"backup_start_at"`
	CreatedAt           string                         `json:"created_at"`
	DatastoreId         string                         `json:"datastore_id"`
	FinishedAt          string                         `json:"finished_at,omitempty"`
	FlavorId            string                         `json:"flavor_id"`
	Generation          string                         `json:"generation"`
	Id                  string                         `json:"id"`
	Name                string                         `json:"name"`
	Replicas            ListResultResultsItemReplicas  `json:"replicas,omitempty"`
	StartedAt           string                         `json:"started_at,omitempty"`
	Status              string                         `json:"status"`
	UpdatedAt           string                         `json:"updated_at,omitempty"`
	Volume              ListResultResultsItemVolume    `json:"volume"`
}

type ListResultResultsItemAddressesItem struct {
	Access  string `json:"access"`
	Address string `json:"address,omitempty"`
	Type    string `json:"type,omitempty"`
}

type ListResultResultsItemAddresses []ListResultResultsItemAddressesItem

type ListResultResultsItemReplicasItem struct {
	Addresses   ListResultResultsItemReplicasItemAddresses `json:"addresses"`
	CreatedAt   string                                     `json:"created_at"`
	DatastoreId string                                     `json:"datastore_id"`
	FinishedAt  string                                     `json:"finished_at,omitempty"`
	FlavorId    string                                     `json:"flavor_id"`
	Generation  string                                     `json:"generation"`
	Id          string                                     `json:"id"`
	Name        string                                     `json:"name"`
	SourceId    string                                     `json:"source_id"`
	StartedAt   string                                     `json:"started_at,omitempty"`
	Status      string                                     `json:"status"`
	UpdatedAt   string                                     `json:"updated_at,omitempty"`
	Volume      ListResultResultsItemReplicasItemVolume    `json:"volume"`
}

type ListResultResultsItemReplicasItemAddressesItem struct {
	Access  string `json:"access"`
	Address string `json:"address,omitempty"`
	Type    string `json:"type,omitempty"`
}

type ListResultResultsItemReplicasItemAddresses []ListResultResultsItemReplicasItemAddressesItem

type ListResultResultsItemReplicasItemVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

type ListResultResultsItemReplicas []ListResultResultsItemReplicasItem

type ListResultResultsItemVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

type ListResultResults []ListResultResultsItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/dbaas/instances/list"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ListConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
