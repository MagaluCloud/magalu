/*
Executor: list

# Summary

# List all clusters

# Description

Lists all clusters for a user.

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

type ListConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

// Object of the clusters response request.
type ListResult struct {
	Results ListResultResults `json:"results"`
}

// Object of the cluster response request.
type ListResultResultsItem struct {
	Description   string                             `json:"description,omitempty"`
	Id            string                             `json:"id"`
	KubeApiServer ListResultResultsItemKubeApiServer `json:"kube_api_server,omitempty"`
	Name          string                             `json:"name"`
	ProjectId     string                             `json:"project_id"`
	Region        string                             `json:"region,omitempty"`
	Status        ListResultResultsItemStatus        `json:"status,omitempty"`
	Tags          ListResultResultsItemTags          `json:"tags,omitempty"`
	Version       string                             `json:"version,omitempty"`
}

// Information about the Kubernetes API Server of the cluster.
type ListResultResultsItemKubeApiServer struct {
	DisableApiServerFip bool   `json:"disable_api_server_fip,omitempty"`
	FixedIp             string `json:"fixed_ip,omitempty"`
	FloatingIp          string `json:"floating_ip,omitempty"`
	Port                int    `json:"port,omitempty"`
}

// Details about the status of the Kubernetes cluster or node.

type ListResultResultsItemStatus struct {
	Message string `json:"message"`
	State   string `json:"state"`
}

type ListResultResultsItemTags []*string

type ListResultResults []ListResultResultsItem

func List(
	client *mgcClient.Client,
	ctx context.Context,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/kubernetes/cluster/list"), client, ctx)
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
