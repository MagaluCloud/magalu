/*
Executor: get

# Summary

Backup Detail.

# Description

Get a backup detail.

Version: 1.20.0

import "magalu.cloud/lib/products/dbaas/instances/backups"
*/
package backups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	BackupId   string  `json:"backup_id"`
	Exchange   *string `json:"exchange,omitempty"`
	InstanceId string  `json:"instance_id"`
}

type GetConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	CreatedAt  string  `json:"created_at"`
	DbSize     *int    `json:"db_size,omitempty"`
	FinishedAt *string `json:"finished_at,omitempty"`
	Id         string  `json:"id"`
	InstanceId string  `json:"instance_id"`
	Location   *string `json:"location,omitempty"`
	Mode       string  `json:"mode"`
	Name       *string `json:"name,omitempty"`
	Size       *int    `json:"size,omitempty"`
	StartedAt  *string `json:"started_at,omitempty"`
	Status     string  `json:"status"`
	Type       string  `json:"type"`
	UpdatedAt  *string `json:"updated_at,omitempty"`
}

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/dbaas/instances/backups/get"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[GetConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

func (s *service) GetUntilTermination(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	e, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/dbaas/instances/backups/get"), s.client, s.ctx)
	if err != nil {
		return
	}

	exec, ok := e.(mgcCore.TerminatorExecutor)
	if !ok {
		// Not expected, but let's fallback
		return s.Get(
			parameters,
			configs,
		)
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[GetConfigs](configs); err != nil {
		return
	}

	r, err := exec.ExecuteUntilTermination(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
