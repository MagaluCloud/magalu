/*
Executor: create

# Summary

# Create a node pool

# Description

Creates a node pool in a Kubernetes cluster.

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

type CreateParameters struct {
	AutoScale CreateParametersAutoScale `json:"auto_scale,omitempty"`
	ClusterId string                    `json:"cluster_id"`
	Flavor    string                    `json:"flavor"`
	Name      string                    `json:"name"`
	Replicas  int                       `json:"replicas"`
	Tags      CreateParametersTags      `json:"tags,omitempty"`
	Taints    CreateParametersTaints    `json:"taints,omitempty"`
}

// Object specifying properties for updating workload resources in the Kubernetes cluster.

type CreateParametersAutoScale struct {
	MaxReplicas int `json:"max_replicas"`
	MinReplicas int `json:"min_replicas"`
}

type CreateParametersTags []string

type CreateParametersTaintsItem struct {
	Effect string `json:"effect"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type CreateParametersTaints []CreateParametersTaintsItem

type CreateConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

// Object of the node pool response.
type CreateResult struct {
	AutoScale        CreateParametersAutoScale    `json:"auto_scale"`
	CreatedAt        string                       `json:"created_at,omitempty"`
	Id               string                       `json:"id"`
	InstanceTemplate CreateResultInstanceTemplate `json:"instance_template"`
	Labels           CreateResultLabels           `json:"labels"`
	Name             string                       `json:"name"`
	Replicas         int                          `json:"replicas"`
	SecurityGroups   CreateResultSecurityGroups   `json:"securityGroups"`
	Status           CreateResultStatus           `json:"status"`
	Tags             CreateResultTags             `json:"tags,omitempty"`
	Taints           CreateResultTaints           `json:"taints,omitempty"`
	UpdatedAt        string                       `json:"updated_at,omitempty"`
	Zone             CreateResultZone             `json:"zone"`
}

// Template for the instance object used to create machine instances and managed instance groups.

type CreateResultInstanceTemplate struct {
	DiskSize  int                                `json:"disk_size"`
	DiskType  string                             `json:"disk_type"`
	Flavor    CreateResultInstanceTemplateFlavor `json:"flavor"`
	NodeImage string                             `json:"node_image"`
}

// Definition of CPU capacity, RAM, and storage for nodes.
type CreateResultInstanceTemplateFlavor struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Ram  int    `json:"ram"`
	Size int    `json:"size"`
	Vcpu int    `json:"vcpu"`
}

// Key/value pairs attached to the object and used for specification.
type CreateResultLabels struct {
}

type CreateResultSecurityGroups []string

// Details about the status of the node pool or control plane.

type CreateResultStatus struct {
	Messages CreateResultStatusMessages `json:"messages"`
	State    string                     `json:"state"`
}

type CreateResultStatusMessages []string

type CreateResultTags []string

type CreateResultTaints []CreateParametersTaintsItem

type CreateResultZone []string

func Create(
	client *mgcClient.Client,
	ctx context.Context,
	parameters CreateParameters,
	configs CreateConfigs,
) (
	result CreateResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Create", mgcCore.RefPath("/kubernetes/nodepool/create"), client, ctx)
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

// TODO: links
// TODO: related
