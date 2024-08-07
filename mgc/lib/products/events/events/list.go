/*
Executor: list

# Summary

Lists all cloud events.

# Description

Lists all cloud events emitted by other products.

Version: 0.12.0

import "magalu.cloud/lib/products/events/events"
*/
package events

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit       *int                `json:"_limit,omitempty"`
	Offset      *int                `json:"_offset,omitempty"`
	Authid      *string             `json:"authid,omitempty"`
	Data        *ListParametersData `json:"data,omitempty"`
	Id          *string             `json:"id,omitempty"`
	ProductLike *string             `json:"product__like,omitempty"`
	SourceLike  *string             `json:"source__like,omitempty"`
	Tenantid    *string             `json:"tenantid,omitempty"`
	Time        *string             `json:"time,omitempty"`
	TypeLike    *string             `json:"type__like,omitempty"`
}

// The raw data event
type ListParametersData struct {
}

type ListConfigs struct {
	XTenantId *string `json:"X-Tenant-ID,omitempty"`
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Meta    ListResultMeta    `json:"meta"`
	Results ListResultResults `json:"results"`
}

type ListResultMeta struct {
	Count  int  `json:"count"`
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
	Total  int  `json:"total"`
}

// Represent all the fields available in event output, following the Cloud Events Spec.
type ListResultResultsItem struct {
	Authid      string                    `json:"authid"`
	Authtype    string                    `json:"authtype"`
	Data        ListResultResultsItemData `json:"data"`
	Id          string                    `json:"id"`
	Product     string                    `json:"product"`
	Region      *string                   `json:"region,omitempty"`
	Source      string                    `json:"source"`
	Specversion *string                   `json:"specversion,omitempty"`
	Subject     string                    `json:"subject"`
	Tenantid    string                    `json:"tenantid"`
	Time        string                    `json:"time"`
	Type        string                    `json:"type"`
}

// The raw event about the occurrence
type ListResultResultsItemData struct {
}

type ListResultResults []ListResultResultsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/events/events/list"), s.client, s.ctx)
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
