/*
Executor: update

# Summary

Patch node pool replicas by node_pool_id

# Description

Updates nodes from a node pool by nodepool_uuid.

Version: 0.1.0

import "magalu.cloud/lib/products/kubernetes/nodepool"
*/
package nodepool

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type UpdateParameters struct {
	AutoScale  *UpdateParametersAutoScale `json:"auto_scale,omitempty"`
	ClusterId  string                     `json:"cluster_id"`
	NodePoolId string                     `json:"node_pool_id"`
	Replicas   *int                       `json:"replicas,omitempty"`
}

// Object specifying properties for updating workload resources in the Kubernetes cluster.

type UpdateParametersAutoScale struct {
	MaxReplicas int `json:"max_replicas"`
	MinReplicas int `json:"min_replicas"`
}

type UpdateConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Object of the node pool response.
type UpdateResult struct {
	AutoScale        UpdateParametersAutoScale    `json:"auto_scale"`
	CreatedAt        *string                      `json:"created_at,omitempty"`
	Id               string                       `json:"id"`
	InstanceTemplate UpdateResultInstanceTemplate `json:"instance_template"`
	Labels           UpdateResultLabels           `json:"labels"`
	Name             string                       `json:"name"`
	Replicas         int                          `json:"replicas"`
	SecurityGroups   *UpdateResultSecurityGroups  `json:"securityGroups,omitempty"`
	Status           UpdateResultStatus           `json:"status"`
	Tags             *UpdateResultTags            `json:"tags,omitempty"`
	Taints           *UpdateResultTaints          `json:"taints,omitempty"`
	UpdatedAt        *string                      `json:"updated_at,omitempty"`
	Zone             *UpdateResultZone            `json:"zone"`
}

// Template for the instance object used to create machine instances and managed instance groups.

type UpdateResultInstanceTemplate struct {
	DiskSize  int                                `json:"disk_size"`
	DiskType  string                             `json:"disk_type"`
	Flavor    UpdateResultInstanceTemplateFlavor `json:"flavor"`
	NodeImage string                             `json:"node_image"`
}

// Definition of CPU capacity, RAM, and storage for nodes.
type UpdateResultInstanceTemplateFlavor struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Ram  int    `json:"ram"`
	Size int    `json:"size"`
	Vcpu int    `json:"vcpu"`
}

// Key/value pairs attached to the object and used for specification.
type UpdateResultLabels struct {
}

type UpdateResultSecurityGroups []string

// Details about the status of the node pool or control plane.

type UpdateResultStatus struct {
	Messages UpdateResultStatusMessages `json:"messages"`
	State    string                     `json:"state"`
}

type UpdateResultStatusMessages []string

type UpdateResultTags []*string

type UpdateResultTaintsItem struct {
	Effect string `json:"effect"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type UpdateResultTaints []UpdateResultTaintsItem

type UpdateResultZone []string

func (s *service) Update(
	parameters UpdateParameters,
	configs UpdateConfigs,
) (
	result UpdateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Update", mgcCore.RefPath("/kubernetes/nodepool/update"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[UpdateParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[UpdateConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[UpdateResult](r)
}

// TODO: links
// TODO: related
