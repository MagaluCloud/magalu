/*
Executor: list

# Summary

# Lists all available flavors

# Description

Lists all available flavors.

Version: 0.1.0

import "magalu.cloud/lib/products/kubernetes/flavor"
*/
package flavor

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

// Response object for the Flavor request.
type ListResult struct {
	Results ListResultResults `json:"results"`
}

// Lists of available flavors provided by the application.
type ListResultResultsItem struct {
	Bastion      ListResultResultsItemBastion      `json:"bastion"`
	Controlplane ListResultResultsItemControlplane `json:"controlplane"`
	Nodepool     ListResultResultsItemNodepool     `json:"nodepool"`
}

// Definition of CPU capacity, RAM, and storage for nodes.
type ListResultResultsItemBastionItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Ram  int    `json:"ram"`
	Size int    `json:"size"`
	Sku  string `json:"sku"`
	Vcpu int    `json:"vcpu"`
}

type ListResultResultsItemBastion []ListResultResultsItemBastionItem

type ListResultResultsItemControlplane []ListResultResultsItemBastionItem

type ListResultResultsItemNodepool []ListResultResultsItemBastionItem

type ListResultResults []ListResultResultsItem

func (s *service) List(
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/kubernetes/flavor/list"), s.client, s.ctx)
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
