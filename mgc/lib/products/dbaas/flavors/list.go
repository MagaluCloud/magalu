/*
Executor: list

# Summary

List available flavors.

# Description

Returns a list of available flavors. A flavor is a hardware template that defines the size of RAM and vcpu.

Version: 1.17.2

import "magalu.cloud/lib/products/dbaas/flavors"
*/
package flavors

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit  int    `json:"_limit,omitempty"`
	Offset int    `json:"_offset,omitempty"`
	Status string `json:"status,omitempty"`
}

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
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
	FamilyDescription string `json:"family_description"`
	FamilySlug        string `json:"family_slug"`
	Id                string `json:"id"`
	Label             string `json:"label"`
	Name              string `json:"name"`
	Ram               string `json:"ram"`
	Size              string `json:"size"`
	SkuReplica        string `json:"sku_replica"`
	SkuSource         string `json:"sku_source"`
	Vcpu              string `json:"vcpu"`
}

type ListResultResults []ListResultResultsItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/dbaas/flavors/list"), client, ctx)
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
