/*
Executor: versions

# Summary

# Lists all available versions

# Description

Lists all available Kubernetes versions.

Version: 0.1.0

import "magalu.cloud/lib/products/kubernetes/info"
*/
package info

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type VersionsConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Response object for the Version request.
type VersionsResult struct {
	Results VersionsResultResults `json:"results"`
}

// Object representing a Kubernetes version.
type VersionsResultResultsItem struct {
	Deprecated bool   `json:"deprecated"`
	Version    string `json:"version"`
}

type VersionsResultResults []VersionsResultResultsItem

func (s *service) Versions(
	configs VersionsConfigs,
) (
	result VersionsResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Versions", mgcCore.RefPath("/kubernetes/info/versions"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[VersionsConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[VersionsResult](r)
}

// TODO: links
// TODO: related
