/*
Executor: get

# Summary

Get node pool by node_pool_id

# Description

Gets nodes in a node pool by nodepool_uuid.

Version: 0.1.0

import "magalu.cloud/lib/products/kubernetes/nodepool"
*/
package nodepool

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	ClusterId  string `json:"cluster_id"`
	NodePoolId string `json:"node_pool_id"`
}

type GetConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

// Object of the node pool response.
type GetResult struct {
	AutoScale        GetResultAutoScale        `json:"auto_scale"`
	CreatedAt        string                    `json:"created_at,omitempty"`
	Id               string                    `json:"id"`
	InstanceTemplate GetResultInstanceTemplate `json:"instance_template"`
	Labels           GetResultLabels           `json:"labels"`
	Name             string                    `json:"name"`
	Replicas         int                       `json:"replicas"`
	SecurityGroups   GetResultSecurityGroups   `json:"securityGroups,omitempty"`
	Status           GetResultStatus           `json:"status"`
	Tags             GetResultTags             `json:"tags,omitempty"`
	Taints           GetResultTaints           `json:"taints,omitempty"`
	UpdatedAt        string                    `json:"updated_at,omitempty"`
	Zone             *GetResultZone            `json:"zone"`
}

// Object specifying properties for updating workload resources in the Kubernetes cluster.

type GetResultAutoScale struct {
	MaxReplicas int `json:"max_replicas"`
	MinReplicas int `json:"min_replicas"`
}

// Template for the instance object used to create machine instances and managed instance groups.

type GetResultInstanceTemplate struct {
	DiskSize  int                             `json:"disk_size"`
	DiskType  string                          `json:"disk_type"`
	Flavor    GetResultInstanceTemplateFlavor `json:"flavor"`
	NodeImage string                          `json:"node_image"`
}

// Definition of CPU capacity, RAM, and storage for nodes.
type GetResultInstanceTemplateFlavor struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Ram  int    `json:"ram"`
	Size int    `json:"size"`
	Vcpu int    `json:"vcpu"`
}

// Key/value pairs attached to the object and used for specification.
type GetResultLabels struct {
}

type GetResultSecurityGroups []string

// Details about the status of the node pool or control plane.

type GetResultStatus struct {
	Messages GetResultStatusMessages `json:"messages"`
	State    string                  `json:"state"`
}

type GetResultStatusMessages []string

type GetResultTags []*string

type GetResultTaintsItem struct {
	Effect string `json:"effect"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type GetResultTaints []GetResultTaintsItem

type GetResultZone []string

func (s *service) Get(
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/kubernetes/nodepool/get"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[GetConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetResult](r)
}

// TODO: links
// TODO: related
