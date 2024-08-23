/*
Executor: list

# Summary

# List Ssh Keys

# Description

Lists all SSH public keys.

Version: 0.1.0

import "magalu.cloud/lib/products/ssh/public_keys"
*/
package publicKeys

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit  *int    `json:"_limit,omitempty"`
	Offset *int    `json:"_offset,omitempty"`
	Sort   *string `json:"_sort,omitempty"`
}

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Results ListResultResults `json:"results"`
}

type ListResultResultsItem struct {
	Id      *string `json:"id,omitempty"`
	KeyType *string `json:"key_type,omitempty"`
	Name    *string `json:"name,omitempty"`
}

type ListResultResults []ListResultResultsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/ssh/public_keys/list"), s.client, s.ctx)
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
