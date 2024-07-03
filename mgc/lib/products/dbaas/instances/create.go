/*
Executor: create

# Summary

Creates a new database instance.

# Description

Creates a new database instance asynchronously for a tenant.

Version: 1.23.0

import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	BackupRetentionDays *int                        `json:"backup_retention_days,omitempty"`
	BackupStartAt       *string                     `json:"backup_start_at,omitempty"`
	DatastoreId         *string                     `json:"datastore_id,omitempty"`
	EngineId            *string                     `json:"engine_id,omitempty"`
	FlavorId            string                      `json:"flavor_id"`
	Name                string                      `json:"name"`
	Parameters          *CreateParametersParameters `json:"parameters,omitempty"`
	Password            string                      `json:"password"`
	User                string                      `json:"user"`
	Volume              CreateParametersVolume      `json:"volume"`
}

type CreateParametersParametersItem struct {
	Name  string                              `json:"name"`
	Value CreateParametersParametersItemValue `json:"value"`
}

// any of: *float64, *int, *bool, *string
type CreateParametersParametersItemValue any

type CreateParametersParameters []CreateParametersParametersItem

type CreateParametersVolume struct {
	Size int     `json:"size"`
	Type *string `json:"type,omitempty"`
}

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type CreateResult struct {
	Id string `json:"id"`
}

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/dbaas/instances/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
