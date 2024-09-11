/*
Executor: delete

# Description

# Deletes the workspace with the specified name

import "magalu.cloud/lib/products/workspace"
*/
package workspace

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type DeleteParameters struct {
	Name string `json:"name"`
}

type DeleteResult struct {
	Name string `json:"name"`
}

func (s *service) Delete(
	parameters DeleteParameters,
) (
	result DeleteResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Delete", mgcCore.RefPath("/workspace/delete"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DeleteResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) DeleteContext(
	ctx context.Context,
	parameters DeleteParameters,
) (
	result DeleteResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Delete", mgcCore.RefPath("/workspace/delete"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[DeleteParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[DeleteResult](r)
}

// TODO: links
// TODO: related
