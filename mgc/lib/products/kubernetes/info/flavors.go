/*
Executor: flavors

# Summary

# Lists all available flavors

# Description

Lists all available flavors.

Version: 0.1.0

import "magalu.cloud/lib/products/kubernetes/info"
*/
package info

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type FlavorsConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Response object for the Flavor request.
type FlavorsResult struct {
	Results FlavorsResultResults `json:"results"`
}

// Lists of available flavors provided by the application.
type FlavorsResultResultsItem struct {
	Bastion      FlavorsResultResultsItemBastion      `json:"bastion"`
	Controlplane FlavorsResultResultsItemControlplane `json:"controlplane"`
	Nodepool     FlavorsResultResultsItemNodepool     `json:"nodepool"`
}

// Definition of CPU capacity, RAM, and storage for nodes.
type FlavorsResultResultsItemBastionItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Ram  int    `json:"ram"`
	Size int    `json:"size"`
	Sku  string `json:"sku"`
	Vcpu int    `json:"vcpu"`
}

type FlavorsResultResultsItemBastion []FlavorsResultResultsItemBastionItem

type FlavorsResultResultsItemControlplane []FlavorsResultResultsItemBastionItem

type FlavorsResultResultsItemNodepool []FlavorsResultResultsItemBastionItem

type FlavorsResultResults []FlavorsResultResultsItem

func (s *service) Flavors(
	configs FlavorsConfigs,
) (
	result FlavorsResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Flavors", mgcCore.RefPath("/kubernetes/info/flavors"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[FlavorsConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[FlavorsResult](r)
}

// TODO: links
// TODO: related
