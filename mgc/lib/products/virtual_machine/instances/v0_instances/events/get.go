/*
Executor: get

# Summary

# Details Instance Event

# Description

# Returns the event details for the instance

Version: 1.199.0

import "magalu.cloud/lib/products/virtual_machine/instances/v0_instances/events"
*/
package events

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	EventId string `json:"event_id"`
	Id      string `json:"id"`
}

type GetConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

type GetResult struct {
	Event     string  `json:"event"`
	EventId   string  `json:"event_id"`
	Message   *string `json:"message,omitempty"`
	StartTime string  `json:"start_time"`
}

func Get(
	client *mgcClient.Client,
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/virtual-machine/instances/v0-instances/events/get"), client, ctx)
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

// TODO: links
// TODO: related
