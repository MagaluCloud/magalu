/*
Executor: list

# Description

# List all available Configs

import "magalu.cloud/lib/products/config"
*/
package config

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListResult map[string]ListResultProperty

type ListResultProperty struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

func (s *service) List() (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/config/list"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
