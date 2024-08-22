/*
Executor: list

# Description

# List all available workspaces

import "magalu.cloud/lib/products/profile/workspaces"
*/
package workspaces

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListResultItem struct {
	Name string `json:"name"`
}

type ListResult []ListResultItem

func (s *service) List() (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/profile/workspaces/list"), s.client, s.ctx)
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
