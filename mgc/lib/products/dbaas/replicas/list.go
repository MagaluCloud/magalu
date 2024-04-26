/*
Executor: list

# Summary

Replicas List.

# Description

List all replicas for a given instance.

Version: 1.17.2

import "magalu.cloud/lib/products/dbaas/replicas"
*/
package replicas

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit    int    `json:"_limit,omitempty"`
	Offset   int    `json:"_offset,omitempty"`
	Exchange string `json:"exchange,omitempty"`
	SourceId string `json:"source_id,omitempty"`
}

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
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
	Addresses   ListResultResultsItemAddresses `json:"addresses"`
	CreatedAt   string                         `json:"created_at"`
	DatastoreId string                         `json:"datastore_id"`
	EngineId    string                         `json:"engine_id"`
	FinishedAt  string                         `json:"finished_at,omitempty"`
	FlavorId    string                         `json:"flavor_id"`
	Generation  string                         `json:"generation"`
	Id          string                         `json:"id"`
	Name        string                         `json:"name"`
	SourceId    string                         `json:"source_id"`
	StartedAt   string                         `json:"started_at,omitempty"`
	Status      string                         `json:"status"`
	UpdatedAt   string                         `json:"updated_at,omitempty"`
	Volume      ListResultResultsItemVolume    `json:"volume"`
}

type ListResultResultsItemAddressesItem struct {
	Access  string `json:"access"`
	Address string `json:"address,omitempty"`
	Type    string `json:"type,omitempty"`
}

type ListResultResultsItemAddresses []ListResultResultsItemAddressesItem

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
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/dbaas/replicas/list"), s.client, s.ctx)
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
