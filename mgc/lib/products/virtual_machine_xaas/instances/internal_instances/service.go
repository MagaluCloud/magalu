/*
import "magalu.cloud/lib/products/virtual_machine_xaas/instances/internal_instances"
*/
package internalInstances

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Create(parameters CreateParameters, configs CreateConfigs) (err error)
	Delete(configs DeleteConfigs) (result DeleteResult, err error)
	DeleteId(parameters DeleteIdParameters, configs DeleteIdConfigs) (err error)
	DeletePorts(parameters DeletePortsParameters, configs DeletePortsConfigs) (err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Update(parameters UpdateParameters, configs UpdateConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
