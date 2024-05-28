/*
Executor: resize

# Summary

Replica Resize.

# Description

Resize an instance replica.

Version: 1.20.0

import "magalu.cloud/lib/products/dbaas/replicas"
*/
package replicas

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ResizeParameters struct {
	Exchange  *string `json:"exchange,omitempty"`
	FlavorId  string  `json:"flavor_id"`
	ReplicaId string  `json:"replica_id"`
}

type ResizeConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ResizeResult struct {
	Addresses   ResizeResultAddresses  `json:"addresses"`
	CreatedAt   string                 `json:"created_at"`
	DatastoreId string                 `json:"datastore_id"`
	EngineId    string                 `json:"engine_id"`
	FinishedAt  *string                `json:"finished_at,omitempty"`
	FlavorId    string                 `json:"flavor_id"`
	Generation  string                 `json:"generation"`
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Parameters  ResizeResultParameters `json:"parameters"`
	SourceId    string                 `json:"source_id"`
	StartedAt   *string                `json:"started_at,omitempty"`
	Status      string                 `json:"status"`
	UpdatedAt   *string                `json:"updated_at,omitempty"`
	Volume      ResizeResultVolume     `json:"volume"`
}

type ResizeResultAddressesItem struct {
	Access  string  `json:"access"`
	Address *string `json:"address,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type ResizeResultAddresses []ResizeResultAddressesItem

type ResizeResultParametersItem struct {
	Name  string                          `json:"name"`
	Value ResizeResultParametersItemValue `json:"value"`
}

// any of: *float64, *int, *bool, *string
type ResizeResultParametersItemValue any

type ResizeResultParameters []ResizeResultParametersItem

type ResizeResultVolume struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

func (s *service) Resize(
	parameters ResizeParameters,
	configs ResizeConfigs,
) (
	result ResizeResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Resize", mgcCore.RefPath("/dbaas/replicas/resize"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ResizeParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ResizeConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ResizeResult](r)
}

// TODO: links
// TODO: related
