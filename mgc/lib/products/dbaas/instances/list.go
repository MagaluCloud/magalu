/*
Executor: list

# Summary

List all database instances.

# Description

Returns a list of database instances for a x-tenant-id.

Version: 1.27.1

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Expand        *string `json:"_expand,omitempty"`
	Limit         *int    `json:"_limit,omitempty"`
	Offset        *int    `json:"_offset,omitempty"`
	EngineId      *string `json:"engine_id,omitempty"`
	Status        *string `json:"status,omitempty"`
	Volumesize    *int    `json:"volume.size,omitempty"`
	VolumesizeGt  *int    `json:"volume.size__gt,omitempty"`
	VolumesizeGte *int    `json:"volume.size__gte,omitempty"`
	VolumesizeLt  *int    `json:"volume.size__lt,omitempty"`
	VolumesizeLte *int    `json:"volume.size__lte,omitempty"`
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

// Page details about the current request pagination.
type ListResultMeta struct {
	Filters ListResultMetaFilters `json:"filters"`
	Page    ListResultMetaPage    `json:"page"`
}

type ListResultMetaFiltersItem struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type ListResultMetaFilters []ListResultMetaFiltersItem

type ListResultMetaPage struct {
	Count    int `json:"count"`
	Limit    int `json:"limit"`
	MaxLimit int `json:"max_limit"`
	Offset   int `json:"offset"`
	Total    int `json:"total"`
}

type ListResultResultsItem struct {
	Addresses           ListResultResultsItemAddresses          `json:"addresses"`
	BackupRetentionDays int                                     `json:"backup_retention_days"`
	BackupStartAt       string                                  `json:"backup_start_at"`
	CreatedAt           string                                  `json:"created_at"`
	DatastoreId         string                                  `json:"datastore_id"`
	EngineId            string                                  `json:"engine_id"`
	FinishedAt          *string                                 `json:"finished_at,omitempty"`
	FlavorId            string                                  `json:"flavor_id"`
	Generation          string                                  `json:"generation"`
	Id                  string                                  `json:"id"`
	InstanceTypeId      string                                  `json:"instance_type_id"`
	Name                string                                  `json:"name"`
	Parameters          ListResultResultsItemParameters         `json:"parameters"`
	Replicas            *ListResultResultsItemReplicas          `json:"replicas,omitempty"`
	StartedAt           *string                                 `json:"started_at,omitempty"`
	Status              string                                  `json:"status"`
	UpdatedAt           *string                                 `json:"updated_at,omitempty"`
	Volume              ListResultResultsItemReplicasItemVolume `json:"volume"`
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
	Addresses      ListResultResultsItemReplicasItemAddresses  `json:"addresses"`
	CreatedAt      string                                      `json:"created_at"`
	DatastoreId    string                                      `json:"datastore_id"`
	EngineId       string                                      `json:"engine_id"`
	FinishedAt     *string                                     `json:"finished_at,omitempty"`
	FlavorId       string                                      `json:"flavor_id"`
	Generation     string                                      `json:"generation"`
	Id             string                                      `json:"id"`
	InstanceTypeId string                                      `json:"instance_type_id"`
	Name           string                                      `json:"name"`
	Parameters     ListResultResultsItemReplicasItemParameters `json:"parameters"`
	SourceId       string                                      `json:"source_id"`
	StartedAt      *string                                     `json:"started_at,omitempty"`
	Status         string                                      `json:"status"`
	UpdatedAt      *string                                     `json:"updated_at,omitempty"`
	Volume         ListResultResultsItemReplicasItemVolume     `json:"volume"`
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) ListContext(
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/dbaas/instances/list"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
