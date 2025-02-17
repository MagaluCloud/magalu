/*
Executor: create

# Summary

Creates a new database instance.

# Description

Creates a new database instance asynchronously for a tenant.

Version: 1.34.1

import "github.com/MagaluCloud/magalu/mgc/lib/products/dbaas/instances"
*/
package instances

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type CreateParameters struct {
	BackupRetentionDays *int                        `json:"backup_retention_days,omitempty"`
	BackupStartAt       *string                     `json:"backup_start_at,omitempty"`
	DatastoreId         *string                     `json:"datastore_id,omitempty"`
	EngineId            *string                     `json:"engine_id,omitempty"`
	FlavorId            *string                     `json:"flavor_id,omitempty"`
	InstanceTypeId      *string                     `json:"instance_type_id,omitempty"`
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) CreateContext(
	ctx context.Context,
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/dbaas/instances/create"), s.client, ctx)
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
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
