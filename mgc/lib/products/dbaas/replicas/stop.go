/*
Executor: stop

# Summary

Replica Stop.

# Description

Stop an instance replica.

Version: 1.17.2

import "magalu.cloud/lib/products/dbaas/replicas"
*/
package replicas

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type StopParameters struct {
	Exchange  *string `json:"exchange,omitempty"`
	ReplicaId string  `json:"replica_id"`
}

type StopConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type StopResult struct {
	Addresses   StopResultAddresses `json:"addresses"`
	CreatedAt   string              `json:"created_at"`
	DatastoreId string              `json:"datastore_id"`
	EngineId    string              `json:"engine_id"`
	FinishedAt  *string             `json:"finished_at,omitempty"`
	FlavorId    string              `json:"flavor_id"`
	Generation  string              `json:"generation"`
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	SourceId    string              `json:"source_id"`
	StartedAt   *string             `json:"started_at,omitempty"`
	Status      string              `json:"status"`
	UpdatedAt   *string             `json:"updated_at,omitempty"`
	Volume      StopResultVolume    `json:"volume"`
}

type StopResultAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type StopResultAddresses []StopResultAddressesItem

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
	exec, ctx, err := mgcHelpers.PrepareExecutor("Stop", mgcCore.RefPath("/dbaas/replicas/stop"), s.client, s.ctx)
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
