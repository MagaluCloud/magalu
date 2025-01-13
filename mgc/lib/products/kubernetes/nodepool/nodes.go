/*
Executor: nodes

# Summary

List nodes from a node pool by node_pool_id

# Description

Lists nodes in a node pool by nodepool_uuid.

Version: 0.1.0

import "github.com/MagaluCloud/magalu/mgc/lib/products/kubernetes/nodepool"
*/
package nodepool

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type NodesParameters struct {
	ClusterId  string `json:"cluster_id"`
	NodePoolId string `json:"node_pool_id"`
}

type NodesConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

// Object of the node response.
type NodesResult struct {
	Results NodesResultResults `json:"results"`
}

// Object of the node response request.
type NodesResultResultsItem struct {
	Addresses      NodesResultResultsItemAddresses      `json:"addresses"`
	Annotations    NodesResultResultsItemAnnotations    `json:"annotations"`
	ClusterName    string                               `json:"cluster_name"`
	CreatedAt      string                               `json:"created_at"`
	Flavor         string                               `json:"flavor"`
	Id             string                               `json:"id"`
	Infrastructure NodesResultResultsItemInfrastructure `json:"infrastructure"`
	Labels         NodesResultResultsItemLabels         `json:"labels"`
	Name           string                               `json:"name"`
	Namespace      string                               `json:"namespace"`
	NodeImage      string                               `json:"node_image"`
	NodepoolName   string                               `json:"nodepool_name"`
	Status         NodesResultResultsItemStatus         `json:"status"`
	Taints         NodesResultResultsItemTaints         `json:"taints"`
	Zone           *string                              `json:"zone"`
}

// Information for the node's address.
type NodesResultResultsItemAddressesItem struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}

type NodesResultResultsItemAddresses []NodesResultResultsItemAddressesItem

// Key/value pairs attached to the object and used for specification.
type NodesResultResultsItemAnnotations struct {
}

// Information about the node's infrastructure.
type NodesResultResultsItemInfrastructure struct {
	Allocatable             NodesResultResultsItemInfrastructureAllocatable `json:"allocatable"`
	Architecture            string                                          `json:"architecture"`
	Capacity                NodesResultResultsItemInfrastructureAllocatable `json:"capacity"`
	ContainerRuntimeVersion string                                          `json:"containerRuntimeVersion"`
	KernelVersion           string                                          `json:"kernelVersion"`
	KubeProxyVersion        string                                          `json:"kubeProxyVersion"`
	KubeletVersion          string                                          `json:"kubeletVersion"`
	OperatingSystem         string                                          `json:"operatingSystem"`
	OsImage                 string                                          `json:"osImage"`
}

// Information about node resources.
type NodesResultResultsItemInfrastructureAllocatable struct {
	Cpu              string `json:"cpu"`
	EphemeralStorage string `json:"ephemeral_storage"`
	Hugepages1gi     string `json:"hugepages_1Gi"`
	Hugepages2mi     string `json:"hugepages_2Mi"`
	Memory           string `json:"memory"`
	Pods             string `json:"pods"`
}

// Key/value pairs attached to the object and used for specification.
type NodesResultResultsItemLabels struct {
}

// Details about the status of the Kubernetes cluster or node.

type NodesResultResultsItemStatus struct {
	Message string `json:"message"`
	State   string `json:"state"`
}

type NodesResultResultsItemTaintsItem struct {
	Effect string `json:"effect"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type NodesResultResultsItemTaints []NodesResultResultsItemTaintsItem

type NodesResultResults []NodesResultResultsItem

func (s *service) Nodes(
	parameters NodesParameters,
	configs NodesConfigs,
) (
	result NodesResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Nodes", mgcCore.RefPath("/kubernetes/nodepool/nodes"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[NodesParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[NodesConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[NodesResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) NodesContext(
	ctx context.Context,
	parameters NodesParameters,
	configs NodesConfigs,
) (
	result NodesResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Nodes", mgcCore.RefPath("/kubernetes/nodepool/nodes"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[NodesParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[NodesConfigs](configs); err != nil {
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
	return mgcHelpers.ConvertResult[NodesResult](r)
}

// TODO: links
// TODO: related
