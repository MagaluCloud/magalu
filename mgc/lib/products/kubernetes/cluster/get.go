/*
Executor: get

# Summary

# List a cluster by uuid

# Description

Lists detailed cluster information by cluster_uuid.

Version: 0.1.0

import "magalu.cloud/lib/products/kubernetes/cluster"
*/
package cluster

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetParameters struct {
	ClusterId string `json:"cluster_id"`
}

type GetConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

// Object of the cluster response request by uuid.
type GetResult struct {
	Addons        GetResultAddons        `json:"addons,omitempty"`
	Controlplane  GetResultControlplane  `json:"controlplane,omitempty"`
	CreatedAt     string                 `json:"created_at,omitempty"`
	Description   string                 `json:"description,omitempty"`
	Id            string                 `json:"id"`
	KubeApiServer GetResultKubeApiServer `json:"kube_api_server,omitempty"`
	Name          string                 `json:"name"`
	Network       GetResultNetwork       `json:"network,omitempty"`
	NodePools     GetResultNodePools     `json:"node_pools,omitempty"`
	ProjectId     string                 `json:"project_id,omitempty"`
	Region        string                 `json:"region"`
	Status        GetResultStatus        `json:"status,omitempty"`
	Tags          GetResultTags          `json:"tags,omitempty"`
	UpdatedAt     string                 `json:"updated_at,omitempty"`
	Version       string                 `json:"version"`
}

// Object representing addons that extend the functionality of the Kubernetes cluster.

type GetResultAddons struct {
	Loadbalance string `json:"loadbalance"`
	Secrets     string `json:"secrets"`
	Volume      string `json:"volume"`
}

// Object of the node pool response.
type GetResultControlplane struct {
	AutoScale        GetResultControlplaneAutoScale        `json:"auto_scale"`
	CreatedAt        string                                `json:"created_at,omitempty"`
	Id               string                                `json:"id"`
	InstanceTemplate GetResultControlplaneInstanceTemplate `json:"instance_template"`
	Labels           GetResultControlplaneLabels           `json:"labels"`
	Name             string                                `json:"name"`
	Replicas         int                                   `json:"replicas"`
	SecurityGroups   GetResultControlplaneSecurityGroups   `json:"securityGroups,omitempty"`
	Status           GetResultControlplaneStatus           `json:"status"`
	Tags             GetResultControlplaneTags             `json:"tags,omitempty"`
	Taints           GetResultControlplaneTaints           `json:"taints,omitempty"`
	UpdatedAt        string                                `json:"updated_at,omitempty"`
	Zone             *GetResultControlplaneZone            `json:"zone"`
}

// Object specifying properties for updating workload resources in the Kubernetes cluster.

type GetResultControlplaneAutoScale struct {
	MaxReplicas int `json:"max_replicas"`
	MinReplicas int `json:"min_replicas"`
}

// Template for the instance object used to create machine instances and managed instance groups.

type GetResultControlplaneInstanceTemplate struct {
	DiskSize  int                                         `json:"disk_size"`
	DiskType  string                                      `json:"disk_type"`
	Flavor    GetResultControlplaneInstanceTemplateFlavor `json:"flavor"`
	NodeImage string                                      `json:"node_image"`
}

// Definition of CPU capacity, RAM, and storage for nodes.
type GetResultControlplaneInstanceTemplateFlavor struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Ram  int    `json:"ram"`
	Size int    `json:"size"`
	Vcpu int    `json:"vcpu"`
}

// Key/value pairs attached to the object and used for specification.
type GetResultControlplaneLabels struct {
}

type GetResultControlplaneSecurityGroups []string

// Details about the status of the node pool or control plane.

type GetResultControlplaneStatus struct {
	Messages GetResultControlplaneStatusMessages `json:"messages"`
	State    string                              `json:"state"`
}

type GetResultControlplaneStatusMessages []string

type GetResultControlplaneTags []*string

type GetResultControlplaneTaintsItem struct {
	Effect string `json:"effect"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type GetResultControlplaneTaints []GetResultControlplaneTaintsItem

type GetResultControlplaneZone []string

// Information about the Kubernetes API Server of the cluster.
type GetResultKubeApiServer struct {
	DisableApiServerFip bool   `json:"disable_api_server_fip,omitempty"`
	FixedIp             string `json:"fixed_ip,omitempty"`
	FloatingIp          string `json:"floating_ip,omitempty"`
	Port                int    `json:"port,omitempty"`
}

// Response object for the Kubernetes cluster network resource request.

type GetResultNetwork struct {
	Cidr     string `json:"cidr"`
	Name     string `json:"name"`
	SubnetId string `json:"subnet_id"`
	Uuid     string `json:"uuid"`
}

// Object of the node pool response.
type GetResultNodePoolsItem struct {
	AutoScale        GetResultControlplaneAutoScale         `json:"auto_scale"`
	CreatedAt        string                                 `json:"created_at,omitempty"`
	Id               string                                 `json:"id"`
	InstanceTemplate GetResultNodePoolsItemInstanceTemplate `json:"instance_template"`
	Labels           GetResultControlplaneLabels            `json:"labels"`
	Name             string                                 `json:"name"`
	Replicas         int                                    `json:"replicas"`
	SecurityGroups   GetResultControlplaneSecurityGroups    `json:"securityGroups,omitempty"`
	Status           GetResultControlplaneStatus            `json:"status"`
	Tags             GetResultControlplaneTags              `json:"tags,omitempty"`
	Taints           GetResultNodePoolsItemTaints           `json:"taints,omitempty"`
	UpdatedAt        string                                 `json:"updated_at,omitempty"`
	Zone             *GetResultControlplaneZone             `json:"zone"`
}

// Template for the instance object used to create machine instances and managed instance groups.

type GetResultNodePoolsItemInstanceTemplate struct {
	DiskSize  int                                         `json:"disk_size"`
	DiskType  string                                      `json:"disk_type"`
	Flavor    GetResultControlplaneInstanceTemplateFlavor `json:"flavor"`
	NodeImage string                                      `json:"node_image"`
}

type GetResultNodePoolsItemTaints []GetResultControlplaneTaintsItem

type GetResultNodePools []GetResultNodePoolsItem

// Details about the status of the Kubernetes cluster or node.

type GetResultStatus struct {
	Message string `json:"message"`
	State   string `json:"state"`
}

type GetResultTags []*string

func Get(
	client *mgcClient.Client,
	ctx context.Context,
	parameters GetParameters,
	configs GetConfigs,
) (
	result GetResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Get", mgcCore.RefPath("/kubernetes/cluster/get"), client, ctx)
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
