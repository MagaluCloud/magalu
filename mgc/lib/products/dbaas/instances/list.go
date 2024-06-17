/*
Executor: list

# Summary

List all database instances.

# Description

Returns a list of database instances for a x-tenant-id.

Version: 1.21.2

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Expand   *string `json:"_expand,omitempty"`
	Limit    *int    `json:"_limit,omitempty"`
	Offset   *int    `json:"_offset,omitempty"`
	Exchange *string `json:"exchange,omitempty"`
	Status   *string `json:"status,omitempty"`
}

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Meta    ListResultMeta    `json:"meta"`
	Results ListResultResults `json:"results"`
}

type ListResultMeta struct {
	Page ListResultMetaPage `json:"page"`
}

type ListResultMetaPage struct {
	Count    int `json:"count"`
	Limit    int `json:"limit"`
	MaxLimit int `json:"max_limit"`
	Offset   int `json:"offset"`
	Total    int `json:"total"`
}

type ListResultResultsItem struct {
	Addresses           ListResultResultsItemAddresses  `json:"addresses"`
	BackupRetentionDays int                             `json:"backup_retention_days"`
	BackupStartAt       string                          `json:"backup_start_at"`
	CreatedAt           string                          `json:"created_at"`
	DatastoreId         string                          `json:"datastore_id"`
	EngineId            string                          `json:"engine_id"`
	FinishedAt          *string                         `json:"finished_at,omitempty"`
	FlavorId            string                          `json:"flavor_id"`
	Generation          string                          `json:"generation"`
	Id                  string                          `json:"id"`
	Name                string                          `json:"name"`
	Parameters          ListResultResultsItemParameters `json:"parameters"`
	Replicas            *ListResultResultsItemReplicas  `json:"replicas,omitempty"`
	StartedAt           *string                         `json:"started_at,omitempty"`
	Status              string                          `json:"status"`
	UpdatedAt           *string                         `json:"updated_at,omitempty"`
	Volume              ListResultResultsItemVolume     `json:"volume"`
}

type ListResultResultsItemAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type ListResultResultsItemAddresses []ListResultResultsItemAddressesItem

type ListResultResultsItemParametersItem struct {
	Name  string                                   `json:"name"`
	Value ListResultResultsItemParametersItemValue `json:"value"`
}

// any of: *float64, *int, *bool, *string
type ListResultResultsItemParametersItemValue any

type ListResultResultsItemParameters []ListResultResultsItemParametersItem

type ListResultResultsItemReplicasItem struct {
	Addresses   ListResultResultsItemReplicasItemAddresses  `json:"addresses"`
	CreatedAt   string                                      `json:"created_at"`
	DatastoreId string                                      `json:"datastore_id"`
	EngineId    string                                      `json:"engine_id"`
	FinishedAt  *string                                     `json:"finished_at,omitempty"`
	FlavorId    string                                      `json:"flavor_id"`
	Generation  string                                      `json:"generation"`
	Id          string                                      `json:"id"`
	Name        string                                      `json:"name"`
	Parameters  ListResultResultsItemReplicasItemParameters `json:"parameters"`
	SourceId    string                                      `json:"source_id"`
	StartedAt   *string                                     `json:"started_at,omitempty"`
	Status      string                                      `json:"status"`
	UpdatedAt   *string                                     `json:"updated_at,omitempty"`
	Volume      ListResultResultsItemReplicasItemVolume     `json:"volume"`
}

type ListResultResultsItemReplicasItemAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type ListResultResultsItemReplicasItemAddresses []ListResultResultsItemReplicasItemAddressesItem

type ListResultResultsItemReplicasItemParametersItem struct {
	Name  string                                   `json:"name"`
	Value ListResultResultsItemParametersItemValue `json:"value"`
}

type ListResultResultsItemReplicasItemParameters []ListResultResultsItemReplicasItemParametersItem

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

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/dbaas/instances/list"), s.client, s.ctx)
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
