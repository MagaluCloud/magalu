/*
Executor: create

# Summary

# Create a cluster

# Description

Creates a Kubernetes cluster in Magalu Cloud.

Version: 0.1.0

import "magalu.cloud/lib/products/kubernetes/cluster"
*/
package cluster

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateParameters struct {
	AllowedCidrs       *CreateParametersAllowedCidrs `json:"allowed_cidrs,omitempty"`
	Description        *string                       `json:"description,omitempty"`
	EnabledBastion     *bool                         `json:"enabled_bastion,omitempty"`
	EnabledServerGroup *bool                         `json:"enabled_server_group,omitempty"`
	Name               string                        `json:"name"`
	NodePools          *CreateParametersNodePools    `json:"node_pools,omitempty"`
	Version            *string                       `json:"version,omitempty"`
}

type CreateParametersAllowedCidrs []string

// Object of the node pool request
type CreateParametersNodePoolsItem struct {
	AutoScale *CreateParametersNodePoolsItemAutoScale `json:"auto_scale,omitempty"`
	Flavor    string                                  `json:"flavor"`
	Name      string                                  `json:"name"`
	Replicas  int                                     `json:"replicas"`
	Tags      *CreateParametersNodePoolsItemTags      `json:"tags,omitempty"`
	Taints    *CreateParametersNodePoolsItemTaints    `json:"taints,omitempty"`
}

// Object specifying properties for updating workload resources in the Kubernetes cluster.

type CreateParametersNodePoolsItemAutoScale struct {
	MaxReplicas *int `json:"max_replicas"`
	MinReplicas *int `json:"min_replicas"`
}

type CreateParametersNodePoolsItemTags []string

type CreateParametersNodePoolsItemTaintsItem struct {
	Effect string `json:"effect"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type CreateParametersNodePoolsItemTaints []CreateParametersNodePoolsItemTaintsItem

type CreateParametersNodePools []CreateParametersNodePoolsItem

type CreateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Object of the cluster response request.
type CreateResult struct {
	AllowedCidrs *CreateResultAllowedCidrs `json:"allowed_cidrs,omitempty"`
	Id           string                    `json:"id"`
	Name         string                    `json:"name"`
	Status       CreateResultStatus        `json:"status"`
}

type CreateResultAllowedCidrs []string

// Details about the status of the Kubernetes cluster or node.

type CreateResultStatus struct {
	Message string `json:"message"`
	State   string `json:"state"`
}

func (s *service) Create(
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/kubernetes/cluster/create"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) CreateContext(
	ctx context.Context,
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/kubernetes/cluster/create"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateConfigs](configs); err != nil {
		return
	}

	sdkConfig := s.client.Sdk().Config().TempConfig()
	if c["serverUrl"] == nil && sdkConfig["serverUrl"] != nil {
		c["serverUrl"] = sdkConfig["serverUrl"]
	}

	if c["env"] == nil && sdkConfig["env"] != nil {
		c["env"] = sdkConfig["env"]
	}

	if c["region"] == nil && sdkConfig["region"] != nil {
		c["region"] = sdkConfig["region"]
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateResult](r)
}

// TODO: links
// TODO: related
