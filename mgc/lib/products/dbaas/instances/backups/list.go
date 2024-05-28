/*
Executor: list

# Summary

Backups List.

# Description

List all backups.

Version: 1.20.0

import "magalu.cloud/lib/products/dbaas/instances/backups"
*/
package backups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit      *int    `json:"_limit,omitempty"`
	Offset     *int    `json:"_offset,omitempty"`
	Exchange   *string `json:"exchange,omitempty"`
	InstanceId string  `json:"instance_id"`
	Mode       *string `json:"mode,omitempty"`
	Status     *string `json:"status,omitempty"`
	Type       *string `json:"type,omitempty"`
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

type ListResultResults []ListResultResultsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/dbaas/instances/backups/list"), s.client, s.ctx)
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
