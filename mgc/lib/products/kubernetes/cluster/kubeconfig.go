/*
Executor: kubeconfig

# Summary

# Get kubeconfig cluster

# Description

Retrieves the kubeconfig of a Kubernetes cluster by cluster_uuid.

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

type KubeconfigParameters struct {
	ClusterId string `json:"cluster_id"`
}

type KubeconfigConfigs struct {
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

func Kubeconfig(
	client *mgcClient.Client,
	ctx context.Context,
	parameters KubeconfigParameters,
	configs KubeconfigConfigs,
) (
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Kubeconfig", mgcCore.RefPath("/kubernetes/cluster/kubeconfig"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[KubeconfigParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[KubeconfigConfigs](configs); err != nil {
		return
	}

	_, err = exec.Execute(ctx, p, c)
	return
}

// TODO: links
// TODO: related
