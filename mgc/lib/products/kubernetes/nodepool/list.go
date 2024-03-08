/*
Executor: list

# Summary

List node pools by cluster_id

# Description

Gets a node pool from a Kubernetes cluster by cluster_uuid.

Version: 0.1.0

import "magalu.cloud/lib/products/kubernetes/nodepool"
*/
package nodepool

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	ClusterId string `json:"cluster_id"`
}

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

// Object of the node pool response in the cluster.
type ListResult struct {
	Results ListResultResults `json:"results"`
}

// Object of the node pool response.
type ListResultResultsItem struct {
	AutoScale        ListResultResultsItemAutoScale        `json:"auto_scale"`
	CreatedAt        string                                `json:"created_at,omitempty"`
	Id               string                                `json:"id"`
	InstanceTemplate ListResultResultsItemInstanceTemplate `json:"instance_template"`
	Labels           ListResultResultsItemLabels           `json:"labels"`
	Name             string                                `json:"name"`
	Replicas         int                                   `json:"replicas"`
	SecurityGroups   ListResultResultsItemSecurityGroups   `json:"securityGroups"`
	Status           ListResultResultsItemStatus           `json:"status"`
	Tags             ListResultResultsItemTags             `json:"tags,omitempty"`
	Taints           ListResultResultsItemTaints           `json:"taints,omitempty"`
	UpdatedAt        string                                `json:"updated_at,omitempty"`
	Zone             ListResultResultsItemZone             `json:"zone"`
}

// Object specifying properties for updating workload resources in the Kubernetes cluster.

type ListResultResultsItemAutoScale struct {
	MaxReplicas int `json:"max_replicas"`
	MinReplicas int `json:"min_replicas"`
}

// Template for the instance object used to create machine instances and managed instance groups.

type ListResultResultsItemInstanceTemplate struct {
	DiskSize  int                                         `json:"disk_size"`
	DiskType  string                                      `json:"disk_type"`
	Flavor    ListResultResultsItemInstanceTemplateFlavor `json:"flavor"`
	NodeImage string                                      `json:"node_image"`
}

// Definition of CPU capacity, RAM, and storage for nodes.
type ListResultResultsItemInstanceTemplateFlavor struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Ram  int    `json:"ram"`
	Size int    `json:"size"`
	Vcpu int    `json:"vcpu"`
}

// Key/value pairs attached to the object and used for specification.
type ListResultResultsItemLabels struct {
}

type ListResultResultsItemSecurityGroups []string

// Details about the status of the node pool or control plane.

type ListResultResultsItemStatus struct {
	Messages ListResultResultsItemStatusMessages `json:"messages"`
	State    string                              `json:"state"`
}

type ListResultResultsItemStatusMessages []string

type ListResultResultsItemTags []string

type ListResultResultsItemTaintsItem struct {
	Effect string `json:"effect"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type ListResultResultsItemTaints []ListResultResultsItemTaintsItem

type ListResultResultsItemZone []string

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
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/kubernetes/nodepool/list"), client, ctx)
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
