/*
Executor: list

# Summary

# Lists all available versions

# Description

Lists all available Kubernetes versions.

Version: 0.1.0

import "magalu.cloud/lib/products/kubernetes/version"
*/
package version

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

// Response object for the Version request.
type ListResult struct {
	Results ListResultResults `json:"results"`
}

// Object representing a Kubernetes version.
type ListResultResultsItem struct {
	Deprecated bool   `json:"deprecated"`
	Version    string `json:"version"`
}

type ListResultResults []ListResultResultsItem

func (s *service) List(
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/kubernetes/version/list"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

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
