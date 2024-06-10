/*
import "magalu.cloud/lib/products/virtual_machine_xaas/xaas_instances"
*/
package xaasInstances

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Create(parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	Reboot(parameters RebootParameters, configs RebootConfigs) (err error)
	Retype(parameters RetypeParameters, configs RetypeConfigs) (err error)
	Start(parameters StartParameters, configs StartConfigs) (err error)
	Stop(parameters StopParameters, configs StopConfigs) (err error)
	Suspend(parameters SuspendParameters, configs SuspendConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
