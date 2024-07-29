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
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type KubeconfigParameters struct {
	ClusterId string `json:"cluster_id"`
}

type KubeconfigConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

func (s *service) Kubeconfig(parameters KubeconfigParameters, configs KubeconfigConfigs) (string, error) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Kubeconfig", mgcCore.RefPath("/kubernetes/cluster/kubeconfig"), s.client, s.ctx)
	if err != nil {
		return "", err
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[KubeconfigParameters](parameters); err != nil {
		return "", err
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[KubeconfigConfigs](configs); err != nil {
		return "", err
	}

	result, err := exec.Execute(ctx, p, c)

	if err != nil {
		return "", err
	}

	output, err := mgcHelpers.ConvertResultReader[string](result)
	if err != nil {
		return "", err
	}

	return output, nil
}
